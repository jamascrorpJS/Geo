FROM golang:1.24.4-alpine

RUN mkdir /application

WORKDIR /application

COPY go.mod go.sum .

RUN go mod download

COPY . .

RUN go build -o main cmd/main.go
CMD ["./main"]