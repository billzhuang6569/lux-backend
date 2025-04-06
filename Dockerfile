FROM golang:1.22-alpine AS builder

# 安装必要的工具
RUN apk add --no-cache git

# 设置工作目录
WORKDIR /app

# 复制go模块文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 确保go.sum文件正确更新
RUN go mod tidy

# 复制源代码
COPY . .

# 构建应用
RUN go build -o lux-api .

# 使用轻量级镜像
FROM alpine:latest

# 安装运行时依赖
RUN apk add --no-cache ffmpeg ca-certificates

# 设置工作目录
WORKDIR /app

# 从构建阶段复制编译好的应用
COPY --from=builder /app/lux-api .

# 创建下载目录
RUN mkdir -p /data/downloads

# 暴露端口
EXPOSE 8080

# 设置环境变量
ENV HOST=0.0.0.0
ENV PORT=8080
ENV DOWNLOAD_DIR=/data/downloads
ENV MAX_CONCURRENT_DOWNLOADS=5
ENV DOWNLOAD_TIMEOUT=3600
ENV FILE_EXPIRY_TIME=86400
ENV CORS_ORIGINS=*

# 运行应用
CMD ["./lux-api"] 