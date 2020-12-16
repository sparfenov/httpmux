FROM golang:1.15 as builder

ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0

WORKDIR /src/

COPY . ./

RUN go build -a ./cmd/app/...


FROM alpine:3.12

RUN apk add --no-cache ca-certificates \
    && update-ca-certificates

WORKDIR /app/

COPY --from=builder ./src/app .

ENV PATH="/app:${PATH}"

EXPOSE 80
ENTRYPOINT ["./app"]
