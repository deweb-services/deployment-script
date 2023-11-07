FROM golang:1.21-alpine3.18 as builder

ADD . /deploy_script
WORKDIR /deploy_script
ARG RESULT

RUN apk update
RUN apk add --no-cache git gcc musl-dev
RUN go build -v -ldflags="-X 'main.version=$VERSION'" -o /tmp/deploy_script ./cmd/deploy_script

FROM alpine:3.18

RUN apk add --no-cache ca-certificates
RUN update-ca-certificates

ARG SERVER_CERT_PATH
ARG SERVER_KEY_PATH

COPY --from=builder /tmp/deploy_script /usr/bin/deploy_script

RUN chmod +x /usr/bin/deploy_script

ENTRYPOINT ["deploy_script"]