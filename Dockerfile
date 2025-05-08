FROM golang:1.24

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

WORKDIR /app

COPY . .

ENV CGO_ENABLED=0
ENV GOOS=linux

RUN go build -o main .

WORKDIR /app

COPY assets .

CMD ["/app/main", "webserver"]
