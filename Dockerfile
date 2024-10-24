FROM golang:1.21
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./data/ /app/data
COPY ./html/ /app/html
COPY ./server/ /app/server
COPY ./static /app/static
COPY config.json ./
RUN CGO_ENABLED=0 GOOS=linux go build -o ./modstats ./server

EXPOSE 8080
CMD ["./modstats"]
