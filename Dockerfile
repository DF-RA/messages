FROM golang:1.23.2-alpine

WORKDIR /go/src/app

COPY . .

RUN go mod download

RUN go build -o app main.go

CMD ["./app"]

