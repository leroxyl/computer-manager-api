FROM golang:1.20.6  AS builder

WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY main.go ./
COPY internal ./internal

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/computer-manager-api

FROM scratch

COPY --from=builder /bin/computer-manager-api /bin/computer-manager-api
EXPOSE 8080

# Run
ENTRYPOINT ["/bin/computer-manager-api"]