FROM golang:1.24.4

WORKDIR /app

RUN apt-get update && \
    apt-get install -y libssl-dev pkg-config && \
    apt-get clean

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o kvstore-api ./cmd/server

CMD [ "./kvstore-api" ]
