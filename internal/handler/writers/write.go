package writers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/warehouse/auth-service/internal/domain"

	"google.golang.org/protobuf/proto"
)

func SendBytes(w http.ResponseWriter, code int, data []byte) {
	w.WriteHeader(code)
	if _, err := w.Write(data); err != nil {
		fmt.Println("write manager couriers report result failed: ", err)
	}
}

func SendJSON(w http.ResponseWriter, code int, payload interface{}) {
	res, _ := json.Marshal(payload)
	// TODO: Добавить обработку ошибок

	w.Header().Set(domain.HeaderContentType, domain.JsonContentType)
	SendBytes(w, code, res)
}

func SendProto(w http.ResponseWriter, code int, m proto.Message) {
	msg, err := proto.Marshal(m)
	if err != nil {
		// TODO: добавить логгер
		return
	}

	w.Header().Set(domain.HeaderContentType, domain.ProtoContentType)
	SendBytes(w, code, msg)
}
