# The build stage
FROM golang:1.21 as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api *.go

# The run stage
FROM scratch
WORKDIR /app
COPY --from=builder /app/api .
EXPOSE 8080
CMD ["./api"]