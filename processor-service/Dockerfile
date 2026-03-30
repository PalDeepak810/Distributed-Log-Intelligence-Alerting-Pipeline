FROM golang:1.26-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o app

EXPOSE 8090

CMD ["./app"]