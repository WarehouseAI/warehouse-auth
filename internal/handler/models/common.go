package models

import (
	"auth-service/internal/domain"
	"auth-service/internal/warehousepb"

	"google.golang.org/protobuf/proto"
)

type (
	ErrorResponse struct {
		Code    int64
		Reason  string
		Details string
	}

	// todo: подобрать имя получше, мб разнести по другому структуры для логичности
	Result struct {
		Resp     *warehousepb.Response
		Accounts []*domain.Account
	}
)

func (r *Result) GetResponse() proto.Message {
	return r.Resp
}
