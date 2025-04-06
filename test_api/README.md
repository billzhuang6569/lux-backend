# Lux API 测试指南

本目录包含用于测试 Lux API 的脚本和说明。您可以使用以下方法之一进行测试。

## 方法 1: 使用 Docker (推荐)

如果您已安装 Docker 和 Docker Compose，这是最简单的测试方法。

### 步骤:

1. 创建下载目录:

```bash
mkdir -p ../downloads
```

2. 启动 Docker 容器:

```bash
docker-compose up --build
```

这将自动构建 API 服务器并运行测试脚本。测试结果将直接显示在控制台中。

## 方法 2: 使用 Python 脚本 (需要单独运行 API 服务器)

如果您能够通过其他方式启动 API 服务器，可以单独使用 Python 脚本进行测试。

### 前提条件:

- Python 3.6+
- pip (Python 包管理器)

### 步骤:

1. 安装依赖:

```bash
pip install requests
```

2. 确保 API 服务器已在 http://localhost:8080 运行

3. 运行测试脚本:

```bash
python test_api.py --url "https://www.bilibili.com/video/BV1MN41127yH"
```

## 测试内容

测试脚本会执行以下操作:

1. 测试 `/api/supported-sites` API - 获取支持的网站列表
2. 测试 `/api/download` API - 使用提供的 URL 创建下载任务
3. 测试 `/api/status/{taskId}` API - 监控下载进度
4. 创建简单的测试网页，用于测试下载链接

## 成功标准

测试成功的标志:

- 所有 API 请求返回 200 状态码
- 下载任务成功创建
- 下载进度能够被监控
- 最终能够通过下载链接获取视频文件

## 注意事项

- 测试可能需要几分钟时间，取决于视频大小和网络连接速度
- 脚本默认等待最多 5 分钟，如果下载时间更长，可以修改 `max_retries` 参数
- 下载的文件将保存在 `../downloads` 目录中