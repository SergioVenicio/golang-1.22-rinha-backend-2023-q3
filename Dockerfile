FROM golang:1.22.0 as builder

WORKDIR /app


COPY go.* ./
COPY ./src ./src

RUN go mod tidy

RUN CGO_ENABLED=0 go build -v -o ./bin/rinha ./src/cmd/rinha.go

FROM alpine:3.19.1

EXPOSE 8080

COPY --from=builder /app/bin/rinha .

ENV GOGC 1000
ENV GOMAXPROCS 8

CMD ["/rinha"]
