gen_dataservice_mocks:
	mockgen -package=mocks -source=internal/app/dataservice/dataservice.go \
	-destination=internal/pkg/mocks/mock_dataservice.go

gen_adapter_mocks:
	mockgen -package=mocks -source=internal/app/adapter/adapter.go \
	-destination=internal/pkg/mocks/mock_adapter.go

convert_proto:
	protoc -I=./proto --go_out=./internal/pkg --go-grpc_out=./internal/pkg user.proto auth.proto
