FROM golang:1.9.0
ARG  APP_ROOT
WORKDIR ${APP_ROOT}
COPY ./ ${APP_ROOT}
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest
ARG APP_ROOT
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app/
RUN mkdir -p ./conf/dev && touch ./conf/dev/app.conf
COPY --from=0 ${APP_ROOT}/main .
COPY --from=0 ${APP_ROOT}/conf/app.conf ./conf/app.conf
COPY --from=0 ${APP_ROOT}/views ./views
COPY --from=0 ${APP_ROOT}/static ./static
ENV DREAMENV TEST
ENV INDOCKER TRUE
ENTRYPOINT ["/app/main"]
