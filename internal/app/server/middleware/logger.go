package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"time"

	"github.com/sirupsen/logrus"
)

func WriteLog(logger *logrus.Logger, fn http.HandlerFunc) http.HandlerFunc {
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

		if recorder.Result().StatusCode != 200 || recorder.Result().StatusCode != 201 {
			if recorder.Result().StatusCode == 500 {
				logger.WithFields(logrus.Fields{
					"request_runtime": time.Since(startTime),
					"response_status": recorder.Code,
					"response_body":   recorder.Result().Body,
				}).Warn(fmt.Sprintf("Can't proceed %s request.", r.URL.Path))
			} else {
				logger.WithFields(logrus.Fields{
					"request_runtime": time.Since(startTime),
					"response_status": recorder.Code,
					"response_body":   recorder.Result().Body,
				}).Info(fmt.Sprintf("Can't proceed %s request.", r.URL.Path))
			}
		} else {
			logger.WithFields(logrus.Fields{
				"request_runtime": time.Since(startTime),
				"response_status": recorder.Code,
			}).Info("Request successfully completed.")
		}
	}
}
