FROM golang:latest as builder

WORKDIR /app
ENV CGO_ENABLED=0
COPY . .

RUN go build -o . cmd/main.go

FROM alpine:latest

COPY --from=builder /app/main /app/main
COPY --from=builder /app/configs /app/configs
WORKDIR /app

EXPOSE 6565
EXPOSE 8080
RUN apk --no-cache add ca-certificates

ENTRYPOINT ["./main"]