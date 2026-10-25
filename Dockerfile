FROM golang:1.25-alpine AS builder
WORKDIR /app

ENV GOPROXY=https://goproxy.cn,direct

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /blog ./cmd/server

FROM alpine:3.20
WORKDIR /app

# 拷贝编译好的程序
COPY --from=builder /blog /app/blog
# 拷贝证书，https接口必备
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8080

# 启动命令
CMD ["./blog"]
