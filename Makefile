.PHONY: pb

gen_dataservice_mocks:
	mockgen -package=mocks -source=internal/app/dataservice/dataservice.go \
	-destination=internal/pkg/mocks/mock_dataservice.go

gen_adapter_mocks:
	mockgen -package=mocks -source=internal/app/adapter/adapter.go \
	-destination=internal/pkg/mocks/mock_adapter.go

pb:
	./build_proto.sh

migrate.create:
	goose -dir migrations create ${name} sql

setup.extensions:
	docker run -it -d --env-file=.env --name db-ext-setup alpine:3.19 ;\
	docker cp ./migrations/extensions/xid.sql db-ext-setup:/ ;\
	docker cp ./setupxid.sh db-ext-setup:/ ;\
	docker exec db-ext-setup /bin/sh -c "apk --update add postgresql-client && sh setupxid.sh" ;\
	docker stop db-ext-setup ;\
	docker rm db-ext-setup