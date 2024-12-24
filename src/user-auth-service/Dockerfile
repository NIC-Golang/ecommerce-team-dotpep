FROM golang:1.23

WORKDIR /app

COPY src/user-auth-service/go.* ./
RUN go mod download

COPY . .
WORKDIR /app/src/user-auth-service/cmd/auth
RUN go build -o /app/main .

CMD ["/app/main"]



