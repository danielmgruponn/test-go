FROM golang:1.23.1-alpine3.20

WORKDIR /app
COPY /test-go .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/api

WORKDIR /app

EXPOSE 8000

CMD [ "/server" ]
