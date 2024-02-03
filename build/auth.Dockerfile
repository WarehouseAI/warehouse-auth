FROM golang:1.21.6-alpine AS build-stage

WORKDIR /auth-service
ENV GOPATH=/
ARG SERVER_PORT
ARG GRPC_PORT

# Download packages only if module files changed
COPY go.mod go.sum ./
RUN go mod download

COPY /proto ./proto
COPY /internal ./internal
COPY /configs ./configs
COPY /cmd ./cmd
RUN CGO_ENABLED=0 GOOS=linux go build -o /auth cmd/main.go

# Deploy the application binary into a lean image
FROM alpine:3.16 AS prod-stage

WORKDIR /

COPY --from=build-stage /auth /auth

EXPOSE ${SERVER_PORT}
EXPOSE ${GRPC_PORT}

# # Download alpine package and install psql-client for the script
# COPY wait-4-postgres.sh ./
# RUN apk update
# RUN apk add postgresql-client
# RUN chmod +x wait-4-postgres.sh

ENTRYPOINT ["/auth"]