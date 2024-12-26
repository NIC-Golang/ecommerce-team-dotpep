FROM golang:1.23

WORKDIR /app

COPY src/core-shop-service/go.* ./

RUN go mod download

COPY . .
WORKDIR /app/src/core-shop-service/cmd
RUN go build -o /app/main .

CMD ["/app/main"]


