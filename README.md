# Lux API 服务器

这是基于 [Lux](https://github.com/iawia002/lux) 视频下载器的 API 服务器，提供了与前端交互的 RESTful API。

## 功能

- 从各种视频网站下载视频
- 支持选择视频格式和质量
- 实时进度跟踪
- 支持多种输出格式（mp4, flv, mp3, webm, m4a）
- 支持多线程下载
- 下载完成后提供文件下载链接

## API 端点

### 1. 下载视频 API

- **端点**: `/api/download`
- **方法**: POST
- **功能**: 接收视频 URL 和下载选项，创建下载任务
- **请求体**:

```json
{
  "url": "https://www.example.com/video",
  "format": "mp4",
  "quality": "best"
}
```

- **成功响应**:

```json
{
  "taskId": "task_12345abcde",
  "status": "pending",
  "message": "下载任务已创建",
  "progress": 0
}
```

### 2. 获取下载状态 API

- **端点**: `/api/status/{taskId}`
- **方法**: GET
- **功能**: 获取特定下载任务的状态和进度

- **成功响应**:

```json
{
  "taskId": "task_12345abcde",
  "status": "processing",
  "progress": 45,
  "message": "正在下载中..."
}
```

### 3. 获取支持的网站 API

- **端点**: `/api/supported-sites`
- **方法**: GET
- **功能**: 返回 Lux 支持的所有视频网站列表

## 配置

通过环境变量进行配置:

```
# 服务器配置
PORT=8080                           # API 服务器端口
HOST=0.0.0.0                        # 监听地址
CORS_ORIGINS=https://your-frontend.vercel.app  # 允许的前端域名

# 下载配置
DOWNLOAD_DIR=/tmp/lux-downloads     # 下载文件存储目录
MAX_CONCURRENT_DOWNLOADS=5          # 最大并发下载数
DOWNLOAD_TIMEOUT=3600               # 下载超时时间（秒）
FILE_EXPIRY_TIME=86400              # 文件过期时间（秒）

# 安全配置
API_KEY=your_secret_key             # API 密钥（可选）
RATE_LIMIT=60                       # 每分钟最大请求数
```

## 运行

```bash
go run main.go
```

## 构建

```bash
go build -o lux-api
```

## Docker

```bash
docker build -t lux-api .
docker run -p 8080:8080 -e PORT=8080 lux-api
```