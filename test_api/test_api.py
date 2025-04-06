#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import requests
import json
import time
import sys
import os
from http.server import HTTPServer, BaseHTTPRequestHandler
import threading
import webbrowser
import argparse

# 测试服务器URL
API_BASE_URL = "http://localhost:8080"

# 彩色输出
class Colors:
    GREEN = '\033[92m'
    YELLOW = '\033[93m'
    RED = '\033[91m'
    BLUE = '\033[94m'
    PURPLE = '\033[95m'
    ENDC = '\033[0m'

def print_colored(text, color):
    print(f"{color}{text}{Colors.ENDC}")

def print_header(text):
    print_colored("\n" + "=" * 50, Colors.BLUE)
    print_colored(f"  {text}", Colors.PURPLE)
    print_colored("=" * 50, Colors.BLUE)

def print_success(text):
    print_colored(f"✅ {text}", Colors.GREEN)

def print_error(text):
    print_colored(f"❌ {text}", Colors.RED)

def print_warning(text):
    print_colored(f"⚠️  {text}", Colors.YELLOW)

def print_info(text):
    print_colored(f"ℹ️  {text}", Colors.BLUE)

def print_json(data):
    print(json.dumps(data, ensure_ascii=False, indent=2))

# 测试下载API
def test_download_api(video_url):
    print_header("测试下载API")
    
    print_info(f"测试URL: {video_url}")
    
    payload = {
        "url": video_url,
        "format": "mp4",
        "quality": "best"
    }
    
    try:
        response = requests.post(f"{API_BASE_URL}/api/download", json=payload)
        
        if response.status_code == 200:
            data = response.json()
            print_success("成功创建下载任务")
            print_json(data)
            return data.get("taskId")
        else:
            print_error(f"请求失败: {response.status_code}")
            print_json(response.json() if response.text else {"error": "无响应内容"})
            return None
    except Exception as e:
        print_error(f"请求异常: {str(e)}")
        return None

# 测试状态API
def test_status_api(task_id):
    print_header(f"测试状态API (任务ID: {task_id})")
    
    if not task_id:
        print_error("缺少任务ID，无法测试状态API")
        return False
    
    try:
        max_retries = 60  # 最多等待60次，每次5秒
        for i in range(max_retries):
            response = requests.get(f"{API_BASE_URL}/api/status/{task_id}")
            
            if response.status_code == 200:
                data = response.json()
                progress = data.get("progress", 0)
                status = data.get("status", "")
                
                print_info(f"进度: {progress}%, 状态: {status}")
                print_json(data)
                
                if status == "completed":
                    print_success("下载完成!")
                    return data
                elif status == "error":
                    print_error("下载失败")
                    return False
                
                # 如果进度没有达到100%，继续等待
                if i < max_retries - 1:
                    print_warning("等待5秒后再次检查...")
                    time.sleep(5)
            else:
                print_error(f"请求失败: {response.status_code}")
                print_json(response.json() if response.text else {"error": "无响应内容"})
                return False
        
        print_error("超时，下载未完成")
        return False
    except Exception as e:
        print_error(f"请求异常: {str(e)}")
        return False

# 测试支持的网站API
def test_supported_sites_api():
    print_header("测试支持的网站API")
    
    try:
        response = requests.get(f"{API_BASE_URL}/api/supported-sites")
        
        if response.status_code == 200:
            data = response.json()
            print_success("成功获取支持的网站列表")
            print_json(data)
            return True
        else:
            print_error(f"请求失败: {response.status_code}")
            print_json(response.json() if response.text else {"error": "无响应内容"})
            return False
    except Exception as e:
        print_error(f"请求异常: {str(e)}")
        return False

# 创建简单的下载页面
def create_download_page(download_info):
    download_url = download_info.get("downloadUrl", "")
    filename = download_info.get("filename", "unknown")
    
    html = f"""
    <!DOCTYPE html>
    <html lang="zh-CN">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>下载测试页面</title>
        <style>
            body {{
                font-family: Arial, sans-serif;
                max-width: 800px;
                margin: 0 auto;
                padding: 20px;
                text-align: center;
                background-color: #f5f5f5;
            }}
            .card {{
                background-color: white;
                border-radius: 10px;
                padding: 20px;
                margin: 20px 0;
                box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
            }}
            .btn {{
                display: inline-block;
                background-color: #4CAF50;
                color: white;
                padding: 10px 20px;
                text-decoration: none;
                border-radius: 5px;
                font-weight: bold;
                margin-top: 20px;
            }}
            h1 {{
                color: #333;
            }}
            .info {{
                color: #666;
                margin-bottom: 20px;
            }}
        </style>
    </head>
    <body>
        <div class="card">
            <h1>下载测试页面</h1>
            <p class="info">文件已成功下载并准备就绪！</p>
            <p>文件名: <strong>{filename}</strong></p>
            <a class="btn" href="{download_url}" target="_blank">下载文件</a>
        </div>
    </body>
    </html>
    """
    
    with open("download_test.html", "w", encoding="utf-8") as f:
        f.write(html)
    
    return os.path.abspath("download_test.html")

# 启动一个简单的HTTP服务器来显示下载页面
class SimpleHTTPRequestHandler(BaseHTTPRequestHandler):
    def do_GET(self):
        if self.path == '/':
            self.send_response(200)
            self.send_header('Content-type', 'text/html')
            self.end_headers()
            
            with open("download_test.html", 'rb') as file:
                self.wfile.write(file.read())
        else:
            self.send_response(404)
            self.end_headers()
            self.wfile.write(b'404 Not Found')
    
    def log_message(self, format, *args):
        # 禁止输出HTTP服务器日志
        return

def start_test_server():
    server = HTTPServer(('localhost', 9000), SimpleHTTPRequestHandler)
    server_thread = threading.Thread(target=server.serve_forever)
    server_thread.daemon = True
    server_thread.start()
    print_info("测试服务器已在 http://localhost:9000 启动")
    return server

def run_tests(video_url):
    print_header("开始API测试")
    
    # 测试支持的网站API
    test_supported_sites_api()
    
    # 测试下载API
    task_id = test_download_api(video_url)
    if not task_id:
        print_error("下载API测试失败，无法继续测试")
        return False
    
    # 测试状态API
    download_info = test_status_api(task_id)
    if not download_info:
        print_error("状态API测试失败或下载未完成")
        return False
    
    # 创建下载页面
    html_file = create_download_page(download_info)
    print_success(f"创建下载测试页面: {html_file}")
    
    # 启动测试服务器并打开浏览器
    server = start_test_server()
    webbrowser.open('http://localhost:9000')
    
    print_header("测试完成")
    print_info("网页已打开，您可以点击下载按钮测试文件下载")
    print_info("按Ctrl+C退出测试")
    
    try:
        while True:
            time.sleep(1)
    except KeyboardInterrupt:
        server.shutdown()
        print_info("测试服务器已关闭")
    
    return True

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description='测试Lux API端点')
    parser.add_argument('--url', type=str, default="https://www.bilibili.com/video/BV1MN41127yH",
                        help='视频链接')
    args = parser.parse_args()
    
    run_tests(args.url)