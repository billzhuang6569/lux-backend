#!/bin/bash

# 颜色输出
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 打印带颜色的消息
print_info() {
    echo -e "${BLUE}[INFO] $1${NC}"
}

print_success() {
    echo -e "${GREEN}[SUCCESS] $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}[WARNING] $1${NC}"
}

print_error() {
    echo -e "${RED}[ERROR] $1${NC}"
}

# 检查Docker是否安装
check_docker() {
    if ! command -v docker &> /dev/null; then
        print_error "Docker未安装。请安装Docker后再试。"
        return 1
    fi
    
    if ! command -v docker-compose &> /dev/null; then
        print_error "Docker Compose未安装。请安装Docker Compose后再试。"
        return 1
    fi
    
    return 0
}

# 检查Python是否安装
check_python() {
    if ! command -v python3 &> /dev/null; then
        print_error "Python 3未安装。请安装Python 3后再试。"
        return 1
    fi
    
    return 0
}

# 创建下载目录
create_download_dir() {
    if [ ! -d "../downloads" ]; then
        print_info "创建下载目录..."
        mkdir -p ../downloads
    fi
}

# 使用Docker运行测试
run_docker_test() {
    print_info "使用Docker运行测试..."
    
    # 创建下载目录
    create_download_dir
    
    # 启动Docker容器
    print_info "启动Docker容器..."
    docker-compose up --build
    
    # 检查结果
    if [ $? -eq 0 ]; then
        print_success "Docker测试完成。"
    else
        print_error "Docker测试失败。"
    fi
}

# 使用Python运行测试
run_python_test() {
    print_info "使用Python运行测试..."
    
    # 安装依赖
    print_info "安装依赖..."
    pip3 install requests || {
        print_error "安装依赖失败。"
        return 1
    }
    
    # 运行测试脚本
    print_info "运行测试脚本..."
    python3 test_api.py --url "https://www.bilibili.com/video/BV1MN41127yH"
    
    # 检查结果
    if [ $? -eq 0 ]; then
        print_success "Python测试完成。"
    else
        print_error "Python测试失败。"
    fi
}

# 主函数
main() {
    echo "==================================================="
    echo "      Lux API 测试脚本                            "
    echo "==================================================="
    echo ""
    echo "选择测试方法:"
    echo "1. 使用Docker (推荐，自动启动API服务器)"
    echo "2. 使用Python (需要手动启动API服务器)"
    echo "3. 退出"
    echo ""
    read -p "请输入选项 (1-3): " choice
    
    case $choice in
        1)
            if check_docker; then
                run_docker_test
            fi
            ;;
        2)
            if check_python; then
                run_python_test
            fi
            ;;
        3)
            print_info "退出测试。"
            exit 0
            ;;
        *)
            print_error "无效的选项。"
            main
            ;;
    esac
}

# 运行主函数
main