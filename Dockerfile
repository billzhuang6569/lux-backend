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

# 首先构建原始的lux命令行工具并安装到/bin目录
RUN go build -o /bin/lux main.go && chmod +x /bin/lux && /bin/lux --help || echo "警告: lux命令行工具构建或运行失败"

# 修改main.go中的名称避免冲突（临时备份）
RUN cp main.go main.go.bak

# 然后构建API服务
RUN go build -o lux-api .

# 恢复main.go
RUN mv main.go.bak main.go

# 创建下载目录
RUN mkdir -p /data/downloads

# 设置环境变量
ENV PATH="/bin:/usr/local/bin:/usr/bin:${PATH}"
ENV HOST=0.0.0.0
ENV PORT=8080
ENV DOWNLOAD_DIR=/data/downloads
ENV MAX_CONCURRENT_DOWNLOADS=5
ENV DOWNLOAD_TIMEOUT=3600
ENV FILE_EXPIRY_TIME=86400
ENV CORS_ORIGINS=*

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./lux-api"] 