#!/bin/bash

# 颜色输出
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 打印带颜色的消息
echo_info() {
    echo -e "${BLUE}[INFO] $1${NC}"
}

echo_success() {
    echo -e "${GREEN}[SUCCESS] $1${NC}"
}

echo_warning() {
    echo -e "${YELLOW}[WARNING] $1${NC}"
}

echo_error() {
    echo -e "${RED}[ERROR] $1${NC}"
}

# 检查Docker是否安装
check_docker() {
    if ! command -v docker &> /dev/null; then
        echo_error "Docker未安装，请先安装Docker"
        exit 1
    fi
    
    if ! command -v docker-compose &> /dev/null; then
        echo_error "Docker Compose未安装，请先安装Docker Compose"
        exit 1
    fi
    
    echo_success "Docker环境检查通过"
}

# 启动Docker容器
start_docker() {
    echo_info "启动Lux API Docker容器..."
    docker-compose -f docker-compose-real.yml up -d
    
    if [ $? -ne 0 ]; then
        echo_error "启动Docker容器失败"
        exit 1
    fi
    
    echo_success "Docker容器已启动"
    echo_info "等待API服务初始化..."
    sleep 5
}

# 测试视频下载
test_download() {
    local url="$1"
    local format="$2"
    local quality="$3"
    
    echo_info "测试下载视频: $url"
    echo_info "格式: $format, 质量: $quality"
    
    # 发送下载请求
    response=$(curl -s -X POST http://localhost:8080/api/download \
      -H "Content-Type: application/json" \
      -d "{\"url\": \"$url\", \"format\": \"$format\", \"quality\": \"$quality\"}")
    
    # 提取任务ID
    task_id=$(echo $response | grep -o '"taskId":"[^"]*"' | cut -d'"' -f4)
    
    if [ -z "$task_id" ]; then
        echo_error "无法获取任务ID，响应: $response"
        return 1
    fi
    
    echo_success "创建下载任务成功，任务ID: $task_id"
    echo_info "开始监控下载进度..."
    
    # 监控下载进度
    status="pending"
    progress=0
    
    while [ "$status" != "completed" ] && [ "$status" != "error" ]; do
        response=$(curl -s -X GET http://localhost:8080/api/status/$task_id)
        status=$(echo $response | grep -o '"status":"[^"]*"' | cut -d'"' -f4)
        progress=$(echo $response | grep -o '"progress":[^,}]*' | cut -d':' -f2)
        message=$(echo $response | grep -o '"message":"[^"]*"' | cut -d'"' -f4)
        
        echo_info "状态: $status, 进度: $progress%, 消息: $message"
        
        if [ "$status" = "error" ]; then
            echo_error "下载失败: $message"
            return 1
        fi
        
        if [ "$status" != "completed" ]; then
            sleep 2
        fi
    done
    
    # 提取下载URL
    download_url=$(echo $response | grep -o '"downloadUrl":"[^"]*"' | cut -d'"' -f4)
    filename=$(echo $response | grep -o '"filename":"[^"]*"' | cut -d'"' -f4)
    
    if [ -z "$download_url" ]; then
        echo_error "无法获取下载URL"
        return 1
    fi
    
    echo_success "下载完成！"
    echo_info "文件名: $filename"
    echo_info "下载链接: $download_url"
    
    # 创建HTML测试页面
    cat > test_download.html << EOF
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Lux 下载测试</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
            text-align: center;
            background-color: #f5f5f5;
        }
        .card {
            background-color: white;
            border-radius: 10px;
            padding: 20px;
            margin: 20px 0;
            box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
        }
        .btn {
            display: inline-block;
            background-color: #4CAF50;
            color: white;
            padding: 10px 20px;
            text-decoration: none;
            border-radius: 5px;
            font-weight: bold;
            margin-top: 20px;
        }
        h1 {
            color: #333;
        }
        .info {
            color: #666;
            margin-bottom: 20px;
        }
    </style>
</head>
<body>
    <div class="card">
        <h1>Lux 下载测试成功</h1>
        <p class="info">视频已成功从Bilibili下载！</p>
        <p>视频链接: <a href="$url" target="_blank">$url</a></p>
        <p>文件名: <strong>$filename</strong></p>
        <a class="btn" href="$download_url" target="_blank">下载视频</a>
    </div>
</body>
</html>
EOF
    
    echo_success "测试页面已创建: test_download.html"
    
    # 尝试打开浏览器
    if command -v open &> /dev/null; then
        open test_download.html
    elif command -v xdg-open &> /dev/null; then
        xdg-open test_download.html
    else
        echo_info "请手动打开文件: $(pwd)/test_download.html"
    fi
    
    return 0
}

# 停止Docker容器
stop_docker() {
    echo_info "停止Docker容器..."
    docker-compose -f docker-compose-real.yml down
    echo_success "Docker容器已停止"
}

# 清理测试文件
cleanup() {
    if [ "$1" = "all" ]; then
        echo_info "清理下载文件..."
        rm -rf ../downloads/*
    fi
    
    echo_info "清理测试文件..."
    rm -f test_download.html
    echo_success "清理完成"
}

# 主函数
main() {
    echo "==================================================="
    echo "      Lux API 真实下载测试                        "
    echo "==================================================="
    
    # 检查Docker
    check_docker
    
    # 清理旧的测试文件
    cleanup
    
    # 创建下载目录
    mkdir -p ../downloads
    
    # 启动Docker
    start_docker
    
    # 测试下载
    url="https://www.bilibili.com/video/BV1MN41127yH"
    test_download "$url" "mp4" "best"
    
    # 提示用户
    echo_info "请查看下载的视频。测试完成后，您可以选择停止Docker容器。"
    echo_info "按任意键停止Docker容器，或按Ctrl+C退出并保持容器运行..."
    read -n 1 -s
    
    # 停止Docker
    stop_docker
}

# 运行主函数
main