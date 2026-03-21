# Собираем в гошке
FROM golang:1.23 AS build

ENV BIN_FILE=/opt/calendar/calendar-sender
ENV CODE_DIR=/go/src/

WORKDIR ${CODE_DIR}

# Кэшируем слои с модулями
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

# Собираем статический бинарник Go (без зависимостей на Си API),
# иначе он не будет работать в alpine образе.
ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} cmd/sender/*

# На выходе тонкий образ
FROM alpine:3.9

LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="calendar"
LABEL MAINTAINERS="student@otus.ru"

EXPOSE 8080
EXPOSE 8090

ENV BIN_FILE="/opt/calendar/calendar-sender"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ENV CONFIG_FILE=/etc/calendar/sender_config.yaml
COPY ./configs/integrationtest/sender_config.yaml ${CONFIG_FILE}

CMD ["/opt/calendar/calendar-sender", "--config", "/etc/calendar/sender_config.yaml"]
