FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init main.go

RUN go build -o main .
RUN chmod +x /app/main

EXPOSE 8080

CMD ["./main"]

# CMD ["go", "run", "main.go"]