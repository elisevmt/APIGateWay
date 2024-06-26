FROM golang:1.19.5-alpine3.16

WORKDIR /app

RUN ls -la

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

CMD ["go", "run", "/app/cmd/api/http/main.go"]
