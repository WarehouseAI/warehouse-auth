package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"time"

	"github.com/sirupsen/logrus"
)

func LoggerMW(logger *logrus.Logger, fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dump, err := httputil.DumpRequest(r, true)

		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}

		startTime := time.Now()

		logger.WithFields(logrus.Fields{
			"request_type":     r.Method,
			"request_endpoint": r.URL.Path,
			"request_dump":     fmt.Sprintf("%q", dump),
		}).Info(fmt.Sprintf("New %s request accepted.", r.URL.Path))

		recorder := httptest.NewRecorder()
		fn(recorder, r)

		if recorder.Result().StatusCode != 200 && recorder.Result().StatusCode != 201 {
			var jsonBody map[string]interface{}
			json.Unmarshal(recorder.Body.Bytes(), &jsonBody)

			if recorder.Result().StatusCode == 500 {
				logger.WithFields(logrus.Fields{
					"request_runtime": time.Since(startTime),
					"response_status": recorder.Code,
					"response_body":   jsonBody,
				}).Warn(fmt.Sprintf("Can't proceed %s request.", r.URL.Path))
			} else {
				logger.WithFields(logrus.Fields{
					"request_runtime": time.Since(startTime),
					"response_status": recorder.Code,
					"response_body":   jsonBody,
				}).Info(fmt.Sprintf("Can't proceed %s request.", r.URL.Path))
			}
		} else {
			logger.WithFields(logrus.Fields{
				"request_runtime": time.Since(startTime),
				"response_status": recorder.Code,
			}).Info("Request successfully completed.")
		}

		for key, values := range recorder.Header() {
			w.Header()[key] = values
		}
		w.WriteHeader(recorder.Code)
		_, _ = w.Write(recorder.Body.Bytes())
	}
}
