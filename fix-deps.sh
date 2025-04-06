#!/bin/bash

# 创建临时Dockerfile来生成正确的go.sum文件
cat > Dockerfile.fix << 'EOF'
FROM golang:1.22-alpine

RUN apk add --no-cache git

WORKDIR /app
COPY go.mod go.sum ./

# 更新go.sum
RUN go mod tidy

# 输出更新后的go.sum
CMD cat go.sum
EOF

# 构建修复镜像
echo "正在构建修复镜像..."
docker build -t lux-fix -f Dockerfile.fix .

# 从容器中提取更新后的go.sum
echo "正在提取更新后的go.sum文件..."
docker run --rm lux-fix > go.sum.new
mv go.sum.new go.sum

# 清理
docker rmi lux-fix
rm Dockerfile.fix

echo "go.sum已更新完成，现在可以重新构建应用" 