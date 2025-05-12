FROM golang:1.23 AS build

WORKDIR /app

COPY . .

RUN go build -o main main.go

FROM golang:1.23

WORKDIR /app

COPY --from=build /app/main /app
COPY --from=build /app/migrations /app/migrations

CMD ["./main"]
