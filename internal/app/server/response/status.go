package response

import "net/http"

func Status(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
}
