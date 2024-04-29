FROM golang:1.22.2-alpine as build-deps

WORKDIR /urs/src/backend

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY . .
RUN go build /usr/src/cmd/warehouse/main.go

FROM alpine:3.19.1

WORKDIR /usr/src/app
ARG env

COPY --from=build-deps /usr/src/backend/run.sh run.sh
COPY --from=build-deps /usr/src/backend/main main
COPY --from=build-deps /usr/src/backend/configs/$env/ config/

ARG service 
ENV LOG_PATH=/logs/$service.log

ENTRYPOINT ["./run.sh"]