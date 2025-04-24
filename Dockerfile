FROM golang:1.16-alpine as builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o app .

FROM alpine:latest  
WORKDIR /root/
COPY --from=builder /app/app .
COPY --from=builder /app/config.json .

CMD ["./app"]
    