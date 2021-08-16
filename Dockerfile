FROM golang:alpine as builder
RUN apk update && apk add --no-cache git
WORKDIR /app 
COPY go.mod go.sum ./
RUN go mod download
COPY . .
WORKDIR /app/cmd
RUN go build -o main .
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/cmd/main .
COPY --from=builder /app/.env .

EXPOSE 8080

CMD ["./main"]
