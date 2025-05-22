FROM golang:1.21 AS builder

# 작업 디렉터리 생성
WORKDIR /app

COPY . .

RUN go mod tidy
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o StoryShift ./cmd/main.go

FROM alpine:latest

COPY --from=builder /app/config /config
COPY --from=builder /app/StoryShift /StoryShift

ENTRYPOINT ["/StoryShift"]
