# Builder Image
# ---------------------------------------------------
FROM golang:alpine AS go-builder

RUN apk update && apk add --no-cache git

WORKDIR /usr/src/app

COPY . ./

RUN go mod download \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -a -o main cmd/main/main.go


# Final Image
# ---------------------------------------------------
FROM alpine:edge

ARG SERVICE_NAME="whatsapp-rest"
ENV PATH="$PATH:/usr/app/${SERVICE_NAME}" \
    CONFIG_ENV="production" \
    PRODUCTION_ROUTER_BASE_PATH="/api"

WORKDIR /usr/app/${SERVICE_NAME}

COPY --from=go-builder /usr/src/app/config/ ./config
COPY --from=go-builder /usr/src/app/main ./whatsapp-rest

RUN chmod 777 config/stores config/uploads

EXPOSE 3000
HEALTHCHECK --interval=5s --timeout=3s CMD ["curl", "http://127.0.0.1:3000${PRODUCTION_ROUTER_BASE_PATH}/health"] || exit 1

VOLUME ["/usr/app/config/stores","/usr/app/config/uploads"]
CMD ["whatsapp-rest"]
