FROM golang:1.19.3-alpine

WORKDIR /cdn

COPY go.mod ./

COPY . .

RUN go mod download

RUN go build -o /cdn

EXPOSE 8080

CMD ["/cdn"]
