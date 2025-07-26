FROM golang:latest

RUN go install github.com/swaggo/swag/cmd/swag@latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN swag init -g main.go --parseDependency --parseInternal

RUN go build -o main .
RUN chmod +x /app/main

EXPOSE 8080

CMD ["./main"]