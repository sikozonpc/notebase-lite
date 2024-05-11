FROM golang:1.22 AS builder
RUN apt-get update
ENV GO111MODULE=on \
  CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64
WORKDIR /app
COPY . .
RUN go install
RUN CGO_ENABLED=0 GOOS=linux go build -o /api *.go

FROM scratch
WORKDIR /
COPY --from=builder /api /api

CMD ["/api"]