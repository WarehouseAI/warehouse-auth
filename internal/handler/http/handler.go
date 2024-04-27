package http

import (
	"auth-service/internal/domain"
	"auth-service/internal/handler/converters"
	"auth-service/internal/handler/models"
	"auth-service/internal/handler/writers"
	"auth-service/internal/pkg/errors"
	"auth-service/internal/pkg/logger"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

const (
	IpHeader = "X-Forwarded-For"
)

type (
	Handler interface {
		FillHandlers(router *mux.Router)
		Shutdown()
	}

	WarehouseRequestHandler interface {
		HandleJsonRequest(router *mux.Router, main, path, method string, handler jsonHandler)
		HandleJsonRequestWithMiddleware(router *mux.Router, main, path, method string, handler jsonHandler, middleware func(http.Handler) http.Handler)
	}

	warehouseRequestHandler struct {
		log           logger.Logger
		cookieTimeout time.Duration
	}

	jsonResponse struct {
		Cookies []http.Cookie
		Data    interface{}   `json:"data;omitempty"`
		Code    int           `json:"code"`
		Error   *errors.Error `json:"error;omitempty"`
	}
	jsonHandler func(ctx context.Context, acc *domain.Account, r *http.Request) jsonResponse
)

func NewWarehouseJsonRequestHandler(
	log logger.Logger,
	cookieTimeout time.Duration,
) WarehouseRequestHandler {
	return &warehouseRequestHandler{
		log:           log,
		cookieTimeout: cookieTimeout,
	}
}

func whJsonErrorResponse(err *errors.Error) jsonResponse {
	return jsonResponse{
		Code:  int(err.Code),
		Error: err,
	}
}

func whJsonSuccessResponse(data interface{}, statusCode int, cookie []http.Cookie) jsonResponse {
	return jsonResponse{
		Data:    data,
		Code:    statusCode,
		Cookies: cookie,
	}
}

func createCookie(name string, value interface{}) (http.Cookie, *errors.Error) {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return http.Cookie{}, errors.WD(errors.InternalError, err)
	}

	return http.Cookie{
		Name:     name,
		Value:    string(jsonData),
		Path:     "/",
		MaxAge:   0,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	}, nil
}

func (wh *warehouseRequestHandler) HandleJsonRequestWithMiddleware(
	router *mux.Router,
	main, path, method string,
	handler jsonHandler,
	middleware func(http.Handler) http.Handler,
) {
	router.Handle(path, middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wh.requestWithJsonResult(main, path, method, w, r, handler)
	}))).Methods(method)
}

func (wh *warehouseRequestHandler) HandleJsonRequest(router *mux.Router, main, path, method string, handler jsonHandler) {
	router.Handle(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wh.requestWithJsonResult(main, path, method, w, r, handler)
	})).Methods(method)
}

func (wh *warehouseRequestHandler) requestWithJsonResult(
	main, path, method string,
	w http.ResponseWriter,
	r *http.Request,
	handler jsonHandler,
) {
	defer r.Body.Close()

	accPayload := r.Context().Value(domain.AccountCtxKey)
	var acc *domain.Account

	if accPayload != nil {
		acc = accPayload.(*domain.Account)
	} else {
		acc = nil
	}

	res := handler(r.Context(), acc, r)
	var resBytes []byte
	if res.Error != nil {
		errResp := converters.MakeJsonErrorResponseWithErrorsError(res.Error)
		wh.logRequest(main, path, method, r.Header.Get(IpHeader), false, &errResp, acc)

		writers.SendJSON(w, res.Code, errResp)
		return
	} else {
		resBytes, _ = json.Marshal(res.Data)
		wh.logRequest(main, path, method, r.Header.Get(IpHeader), true, nil, acc)

		if res.Cookies != nil {
			for _, cookie := range res.Cookies {
				cookie.Expires = time.Now().Add(wh.cookieTimeout)
				http.SetCookie(w, &cookie)
			}
		}

		writers.SendBytes(w, res.Code, resBytes)
		return
	}
}

func (wh *warehouseRequestHandler) logRequest(main, path, method, ip string, success bool, err *models.ErrorResponse, acc *domain.Account) {
	fields := []zap.Field{
		zap.String("method", method), zap.String("path", main+path), zap.String("ip", ip),
	}

	if acc != nil {
		fields = append(fields, zap.String("acc", fmt.Sprintf("%+v", acc)))
	}

	if success {
		wh.log.Info("request", fields...)
	} else {
		fields = append(fields, zap.String("error", fmt.Sprintf("details=%+v, reason=%+v, status=%+v", err.Details, err.Reason, err.Code)))
		wh.log.Zap().Warn("request", fields...)
	}
}
