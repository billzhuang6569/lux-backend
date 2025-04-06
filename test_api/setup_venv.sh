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

# 检查Python版本
check_python() {
    if command -v python3 &> /dev/null; then
        python_cmd="python3"
    elif command -v python &> /dev/null; then
        # 检查Python版本是否为Python 3
        version=$(python --version 2>&1 | awk '{print $2}' | cut -d '.' -f 1)
        if [ "$version" == "3" ]; then
            python_cmd="python"
        else
            echo_error "找不到Python 3。请安装Python 3后再试。"
            exit 1
        fi
    else
        echo_error "找不到Python。请安装Python 3后再试。"
        exit 1
    fi
    
    echo_success "找到Python: $($python_cmd --version)"
    return 0
}

# 创建虚拟环境
create_venv() {
    # 检查是否已存在venv目录
    if [ -d "venv" ]; then
        echo_warning "虚拟环境目录'venv'已存在。是否要删除并重新创建？(y/n)"
        read -r response
        if [ "$response" = "y" ] || [ "$response" = "Y" ]; then
            echo_info "删除现有虚拟环境..."
            rm -rf venv
        else
            echo_info "将使用现有虚拟环境。"
            return 0
        fi
    fi
    
    echo_info "创建Python虚拟环境..."
    $python_cmd -m venv venv
    
    if [ $? -ne 0 ]; then
        echo_error "创建虚拟环境失败。请确保已安装venv包。"
        echo_info "可以尝试运行: $python_cmd -m pip install --user virtualenv"
        exit 1
    fi
    
    echo_success "虚拟环境创建成功！"
}

# 激活虚拟环境并安装依赖
setup_venv() {
    echo_info "激活虚拟环境并安装依赖..."
    
    # 确定激活脚本路径
    if [ "$(uname)" == "Darwin" ] || [ "$(uname)" == "Linux" ]; then
        # macOS or Linux
        activate_script="venv/bin/activate"
    elif [ "$(expr substr $(uname -s) 1 5)" == "MINGW" ] || [ "$(expr substr $(uname -s) 1 6)" == "CYGWIN" ]; then
        # Windows
        activate_script="venv/Scripts/activate"
    else
        echo_error "无法确定操作系统类型。"
        exit 1
    fi
    
    # 提示用户如何激活虚拟环境
    echo_info "请运行以下命令激活虚拟环境:"
    echo_info "source $activate_script"
    echo_info "然后运行:"
    echo_info "pip install requests"
    echo_info "python test_api.py --url \"https://www.bilibili.com/video/BV1MN41127yH\""
    
    echo_success "设置完成！按照上述说明激活虚拟环境并运行测试。"
}

# 创建下载目录
create_download_dir() {
    if [ ! -d "../downloads" ]; then
        echo_info "创建下载目录..."
        mkdir -p ../downloads
        echo_success "下载目录创建成功！"
    else
        echo_info "下载目录已存在。"
    fi
}

# 设置测试API模拟服务器
setup_mock_server() {
    echo_info "创建API模拟服务器脚本..."
    
    # 写入模拟服务器脚本
    cat > mock_api_server.py << 'EOF'
#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import http.server
import json
import socketserver
import os
import time
import threading
import uuid
import shutil
from urllib.parse import urlparse, parse_qs

# 模拟任务状态
tasks = {}
DOWNLOAD_DIR = os.path.abspath("../downloads")

# 确保下载目录存在
if not os.path.exists(DOWNLOAD_DIR):
    os.makedirs(DOWNLOAD_DIR)

class MockAPIHandler(http.server.BaseHTTPRequestHandler):
    def do_OPTIONS(self):
        self.send_response(200)
        self.send_header('Access-Control-Allow-Origin', '*')
        self.send_header('Access-Control-Allow-Methods', 'GET, POST, OPTIONS')
        self.send_header('Access-Control-Allow-Headers', 'Content-Type')
        self.end_headers()
    
    def do_GET(self):
        # 设置CORS头
        self.send_header('Access-Control-Allow-Origin', '*')
        self.send_header('Access-Control-Allow-Methods', 'GET, POST, OPTIONS')
        
        # 根据路径处理不同的API
        if self.path.startswith('/api/supported-sites'):
            self.handle_supported_sites()
        elif self.path.startswith('/api/status/'):
            self.handle_status()
        elif self.path.startswith('/downloads/'):
            self.handle_download_file()
        else:
            self.send_response(404)
            self.send_header('Content-type', 'application/json')
            self.end_headers()
            self.wfile.write(json.dumps({'status': 'error', 'message': '未找到资源'}).encode())
    
    def do_POST(self):
        # 设置CORS头
        self.send_header('Access-Control-Allow-Origin', '*')
        self.send_header('Access-Control-Allow-Methods', 'GET, POST, OPTIONS')
        
        if self.path.startswith('/api/download'):
            self.handle_download()
        else:
            self.send_response(404)
            self.send_header('Content-type', 'application/json')
            self.end_headers()
            self.wfile.write(json.dumps({'status': 'error', 'message': '未找到资源'}).encode())
    
    def handle_supported_sites(self):
        # 返回支持的网站列表
        self.send_response(200)
        self.send_header('Content-type', 'application/json')
        self.end_headers()
        
        sites = [
            "YouTube", "Twitter", "Bilibili", "抖音", "优酷", 
            "腾讯视频", "爱奇艺", "Instagram", "Facebook", "Vimeo",
            "微博", "小红书", "TikTok", "Reddit", "PornHub",
            "AcFun", "Tumblr", "芒果TV", "喜马拉雅", "映客直播",
            "斗鱼", "虎牙", "快手"
        ]
        
        response = {'sites': sites}
        self.wfile.write(json.dumps(response, ensure_ascii=False).encode())
    
    def handle_download(self):
        # 获取请求内容长度
        content_length = int(self.headers['Content-Length'])
        post_data = self.rfile.read(content_length).decode('utf-8')
        
        try:
            request = json.loads(post_data)
            url = request.get('url', '')
            format = request.get('format', 'mp4')
            quality = request.get('quality', 'best')
            
            # 验证参数
            if not url:
                self.send_response(400)
                self.send_header('Content-type', 'application/json')
                self.end_headers()
                self.wfile.write(json.dumps({'status': 'error', 'message': 'URL不能为空'}).encode())
                return
            
            # 创建任务ID
            task_id = str(uuid.uuid4())[:12]
            
            # 创建任务
            tasks[task_id] = {
                'url': url,
                'format': format,
                'quality': quality,
                'status': 'pending',
                'progress': 0,
                'message': '下载任务已创建',
                'created_at': time.time()
            }
            
            # 启动模拟下载线程
            thread = threading.Thread(target=self.simulate_download, args=(task_id,))
            thread.daemon = True
            thread.start()
            
            # 返回任务ID
            self.send_response(200)
            self.send_header('Content-type', 'application/json')
            self.end_headers()
            
            response = {
                'taskId': task_id,
                'status': 'pending',
                'message': '下载任务已创建',
                'progress': 0
            }
            
            self.wfile.write(json.dumps(response).encode())
            
        except json.JSONDecodeError:
            self.send_response(400)
            self.send_header('Content-type', 'application/json')
            self.end_headers()
            self.wfile.write(json.dumps({'status': 'error', 'message': '无效的JSON格式'}).encode())
    
    def handle_status(self):
        # 从路径中提取任务ID
        parts = self.path.split('/')
        if len(parts) < 3:
            self.send_response(400)
            self.send_header('Content-type', 'application/json')
            self.end_headers()
            self.wfile.write(json.dumps({'status': 'error', 'message': '无效的请求路径'}).encode())
            return
        
        task_id = parts[-1]
        
        # 检查任务是否存在
        if task_id not in tasks:
            self.send_response(404)
            self.send_header('Content-type', 'application/json')
            self.end_headers()
            self.wfile.write(json.dumps({'status': 'error', 'message': '任务不存在'}).encode())
            return
        
        # 获取任务状态
        task = tasks[task_id]
        
        # 返回状态
        self.send_response(200)
        self.send_header('Content-type', 'application/json')
        self.end_headers()
        
        response = {
            'taskId': task_id,
            'status': task['status'],
            'progress': task['progress'],
            'message': task['message']
        }
        
        # 如果任务已完成，添加下载URL
        if task['status'] == 'completed':
            filename = f"video_{task_id}.{task['format']}"
            response['downloadUrl'] = f"http://localhost:8080/downloads/{filename}"
            response['filename'] = filename
        
        self.wfile.write(json.dumps(response).encode())
    
    def handle_download_file(self):
        # 从路径中提取文件名
        filename = os.path.basename(self.path)
        file_path = os.path.join(DOWNLOAD_DIR, filename)
        
        # 检查文件是否存在
        if not os.path.exists(file_path) and not self.create_dummy_file(file_path):
            self.send_response(404)
            self.send_header('Content-type', 'application/json')
            self.end_headers()
            self.wfile.write(json.dumps({'status': 'error', 'message': '文件不存在'}).encode())
            return
        
        # 返回文件
        self.send_response(200)
        
        # 设置Content-Type
        if filename.endswith('.mp4'):
            self.send_header('Content-Type', 'video/mp4')
        elif filename.endswith('.mp3'):
            self.send_header('Content-Type', 'audio/mpeg')
        elif filename.endswith('.webm'):
            self.send_header('Content-Type', 'video/webm')
        else:
            self.send_header('Content-Type', 'application/octet-stream')
        
        # 设置Content-Disposition
        self.send_header('Content-Disposition', f'attachment; filename="{filename}"')
        
        # 获取文件大小
        file_size = os.path.getsize(file_path)
        self.send_header('Content-Length', str(file_size))
        
        self.end_headers()
        
        # 发送文件内容
        with open(file_path, 'rb') as f:
            self.wfile.write(f.read())
    
    def create_dummy_file(self, file_path):
        """创建一个虚拟的视频文件用于测试"""
        try:
            # 检查是否有示例视频
            sample_video = os.path.join(os.path.dirname(os.path.abspath(__file__)), 'sample.mp4')
            
            if os.path.exists(sample_video):
                # 复制示例视频
                shutil.copy(sample_video, file_path)
            else:
                # 创建一个1MB的随机文件
                with open(file_path, 'wb') as f:
                    f.write(b'\0' * 1024 * 1024)
            
            return True
        except Exception as e:
            print(f"创建虚拟文件失败: {e}")
            return False
    
    def simulate_download(self, task_id):
        """模拟下载过程"""
        task = tasks[task_id]
        
        # 更新任务状态为处理中
        task['status'] = 'processing'
        task['message'] = '正在下载中...'
        
        # 模拟下载进度
        for i in range(1, 11):
            if task_id not in tasks:
                return  # 任务被删除
            
            progress = i * 10
            task['progress'] = progress
            
            # 每次更新后暂停1秒
            time.sleep(1)
        
        # 下载完成
        if task_id in tasks:
            task['status'] = 'completed'
            task['progress'] = 100
            task['message'] = '下载已完成'
            
            # 创建下载文件
            filename = f"video_{task_id}.{task['format']}"
            file_path = os.path.join(DOWNLOAD_DIR, filename)
            
            # 创建一个虚拟文件
            self.create_dummy_file(file_path)
    
    def log_message(self, format, *args):
        """自定义日志格式"""
        print(f"[{self.log_date_time_string()}] {format % args}")

def run_server(port=8080):
    """启动模拟API服务器"""
    server_address = ('', port)
    httpd = socketserver.ThreadingTCPServer(server_address, MockAPIHandler)
    print(f"启动模拟API服务器在端口 {port}...")
    print(f"支持的API端点:")
    print(f"  GET  /api/supported-sites")
    print(f"  POST /api/download")
    print(f"  GET  /api/status/<task_id>")
    print(f"  GET  /downloads/<filename>")
    print(f"按Ctrl+C停止服务器")
    try:
        httpd.serve_forever()
    except KeyboardInterrupt:
        pass
    httpd.server_close()
    print("服务器已停止")

if __name__ == "__main__":
    port = 8080
    run_server(port)
EOF
    
    chmod +x mock_api_server.py
    echo_success "模拟服务器脚本创建成功！"
}

# 运行服务器说明
run_instructions() {
    echo_info "==================================================="
    echo_info "测试说明"
    echo_info "==================================================="
    echo_info ""
    echo_info "1. 首先激活虚拟环境:"
    echo_info "   source venv/bin/activate"
    echo_info ""
    echo_info "2. 安装依赖:"
    echo_info "   pip install requests"
    echo_info ""
    echo_info "3. 启动模拟API服务器(在一个终端窗口):"
    echo_info "   python mock_api_server.py"
    echo_info ""
    echo_info "4. 在另一个终端窗口运行测试脚本:"
    echo_info "   python test_api.py"
    echo_info ""
    echo_info "5. 查看测试结果和生成的下载页面"
    echo_info ""
    echo_info "完成后，使用 'deactivate' 命令退出虚拟环境"
}

# 主程序
main() {
    echo "==================================================="
    echo "       Python虚拟环境设置和API测试向导              "
    echo "==================================================="
    
    # 检查Python
    check_python
    
    # 创建下载目录
    create_download_dir
    
    # 创建虚拟环境
    create_venv
    
    # 设置模拟服务器
    setup_mock_server
    
    # 设置虚拟环境
    setup_venv
    
    # 显示运行说明
    run_instructions
}

# 运行主程序
main