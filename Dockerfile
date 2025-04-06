FROM golang:1.19-alpine AS builder

WORKDIR /app

# 复制go mod文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 编译
RUN CGO_ENABLED=0 GOOS=linux go build -o lux-api .

# 使用小型镜像
FROM alpine:latest

# 安装ffmpeg (用于视频处理)
RUN apk add --no-cache ffmpeg ca-certificates

WORKDIR /app

# 复制二进制文件
COPY --from=builder /app/lux-api .

# 创建下载目录
RUN mkdir -p /app/downloads

# 设置默认环境变量
ENV PORT=8080
ENV HOST=0.0.0.0
ENV CORS_ORIGINS=*
ENV DOWNLOAD_DIR=/app/downloads
ENV MAX_CONCURRENT_DOWNLOADS=5
ENV DOWNLOAD_TIMEOUT=3600
ENV FILE_EXPIRY_TIME=86400
ENV RATE_LIMIT=60

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./lux-api"]