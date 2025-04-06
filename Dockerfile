FROM golang:1.22-alpine

# 安装必要的工具
RUN apk add --no-cache git ffmpeg ca-certificates

# 设置工作目录
WORKDIR /app

# 复制源代码
COPY . .

# 强制跳过go.sum检查
ENV GONOSUMDB=*

# 获取所有必要的依赖
RUN go get github.com/google/uuid
RUN go get github.com/cheggaaa/pb/v3
RUN go get github.com/fatih/color
RUN go get github.com/MercuryEngineering/CookieMonster
RUN go get github.com/kr/pretty
RUN go get github.com/pkg/errors

# 执行构建
RUN go build -o lux-api .

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