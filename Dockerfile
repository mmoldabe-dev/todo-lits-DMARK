
FROM golang:1.21-alpine AS builder


RUN apk add --no-cache git gcc musl-dev


WORKDIR /app

COPY go.mod go.sum ./


RUN go mod download


COPY . .


RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main ./app/cmd


FROM alpine:latest


RUN apk --no-cache add ca-certificates tzdata


WORKDIR /root/

COPY --from=builder /app/main .


COPY --from=builder /app/migrations ./migrations


RUN adduser -D -s /bin/sh todouser
USER todouser


EXPOSE 8080


CMD ["./main"]