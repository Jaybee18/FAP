FROM golang:1.24.3-alpine3.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o /app/main .

FROM alpine:latest


WORKDIR /root/


COPY --from=builder /app/main .

EXPOSE 8080

# Command to run the executable
# This will run the binary named "main" located in the WORKDIR.
CMD ["./main"]
