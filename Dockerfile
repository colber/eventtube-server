FROM golang:1.14.1-buster AS builder

WORKDIR /app
COPY . /app

RUN go get "github.com/nats-io/go-nats"
RUN go get "github.com/gorilla/websocket"

# For prometheus
RUN go get "github.com/prometheus/client_golang/prometheus"
RUN go get "github.com/prometheus/client_golang/prometheus/promauto"
RUN go get "github.com/prometheus/client_golang/prometheus/promhttp"

# CMD go run main.go

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o dist .
FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /app

COPY --from=builder /app/dist .
COPY --from=builder /app/config.json .
COPY --from=builder /app/sdk.js .


EXPOSE 9005

CMD ["./dist"] 
