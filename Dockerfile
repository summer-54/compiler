FROM docker.io/library/golang:1.24 AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o compiler .

FROM gcc:15.2.0 as final

WORKDIR /app

COPY --from=builder /app/compiler ./compiler
COPY ./testlib.h ./testlib.h

CMD ["./compiler"]
