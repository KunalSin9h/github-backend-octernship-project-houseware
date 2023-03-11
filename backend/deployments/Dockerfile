FROM golang:1.20.2-alpine3.17 as builder

WORKDIR /auth-service
COPY go.mod  .
COPY go.sum .
RUN go mod download


COPY cmd ./cmd
COPY data ./data

RUN cd cmd/api/
RUN go test .
RUN CGO_ENABLED=0 go build -o main .

FROM alpine:3.17

WORKDIR /auth-service
COPY --from=builder /auth-service/main .

CMD ["./main"]