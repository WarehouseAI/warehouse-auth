FROM alpine:3.19.1

WORKDIR /usr/src/app
RUN apk -y update && apk add -y wget
RUN wget https://github.com/pressly/goose/releases/download/v3.4.1/goose_linux_x86_64
RUN chmod +x goose_linux_x86_64
COPY migrations .

ENTRYPOINT /usr/src/app/goose_linux_x86_64 -allow-missing postgres "host=$POSTGRES_HOST port=$POSTGRES_PORT user=$POSTGRES_USER password=$POSTGRES_PASSWORD dbname=$POSTGRES_DB sslmode=$POSTGRES_SSLMODE" up