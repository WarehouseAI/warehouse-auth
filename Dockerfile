FROM golang:1.22.2-bullseye as build-deps

WORKDIR /usr/src/backend

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download
RUN apt-get update
RUN apt install -y protobuf-compiler make
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

COPY . .
RUN make pb
RUN go build /usr/src/backend/cmd/warehouse/main.go

FROM alpine:3.19.1
WORKDIR /usr/src/app
ARG env

COPY --from=build-deps /usr/src/backend/run.sh run.sh
COPY --from=build-deps /usr/src/backend/main main
COPY --from=build-deps /usr/src/backend/configs/$env configs/

ARG service 
ENV LOG_PATH=/logs/$service.log

ENTRYPOINT ["./run.sh"]