FROM alpine:3.20

WORKDIR /usr/src/app

RUN apk update
RUN apk add --no-cache wget
RUN wget https://github.com/pressly/goose/releases/download/v3.4.1/goose_linux_x86_64
RUN chmod +x goose_linux_x86_64
COPY migrations .

ENTRYPOINT /usr/src/app/goose_linux_x86_64 -allow-missing postgres "host=$POSTGRES_HOST port=5432 user=$POSTGRES_USER password=$POSTGRES_PASSWORD dbname=$POSTGRES_DB sslmode=disable" up