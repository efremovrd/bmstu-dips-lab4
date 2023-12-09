FROM golang:1.21.2

COPY ./config/config.go /app/config/config.go
COPY ./config/payment-config.yml /app/config/config.yml
COPY ./payment-service/ /app/payment-service
COPY ./pkg /app/pkg
COPY ./go.mod /app/go.mod
COPY ./go.sum /app/go.sum

RUN export GOPATH=/app

WORKDIR /app

# RUN go mod tidy

RUN go mod download

RUN go build payment-service/cmd/api/main.go

ENTRYPOINT [ "./main" ]
