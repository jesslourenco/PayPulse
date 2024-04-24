FROM golang:1.22

WORKDIR /app

COPY cmd/gopay/main.go .
COPY go.mod .
RUN go mod download
RUN go build -o api

EXPOSE 8080

CMD [ "./api" ]