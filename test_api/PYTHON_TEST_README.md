# Lux API Python测试指南

这个指南将帮助您使用Python虚拟环境来测试Lux API。我们提供了一个模拟API服务器和测试脚本，无需实际的Go环境即可进行测试。

## 快速开始

运行自动化设置脚本:

```bash
./setup_venv.sh
```

这个脚本会:
1. 检查Python安装
2. 创建Python虚拟环境
3. 创建一个模拟的API服务器
4. 提供详细的使用说明

## 手动步骤

如果您想手动完成设置，请按照以下步骤操作:

### 1. 创建并激活Python虚拟环境

```bash
# 创建虚拟环境
python3 -m venv venv

# 激活虚拟环境 (Linux/macOS)
source venv/bin/activate

# 激活虚拟环境 (Windows)
# venv\Scripts\activate
```

### 2. 安装依赖

```bash
pip install requests
```

如果您想创建示例视频(可选):
```bash
pip install moviepy
python create_sample_video.py
```

### 3. 启动模拟API服务器

在一个终端窗口中:

```bash
python mock_api_server.py
```

### 4. 运行测试脚本

在另一个终端窗口中:

```bash
# 确保已激活虚拟环境
source venv/bin/activate  # Linux/macOS
# venv\Scripts\activate   # Windows

# 运行测试
python test_api.py --url "https://www.bilibili.com/video/BV1MN41127yH"
```

## 模拟API服务器说明

模拟API服务器实现了三个主要API端点:

1. `GET /api/supported-sites` - 返回支持的网站列表
2. `POST /api/download` - 创建一个模拟下载任务
3. `GET /api/status/{taskId}` - 获取任务状态
4. `GET /downloads/{filename}` - 下载文件

模拟服务器会:
- 模拟下载进度 (每秒增加10%)
- 生成一个虚拟视频文件用于下载
- 自动返回完整的API响应

## 预期结果

测试脚本将显示每个API调用的详细信息和响应。如果一切正常:

1. 您应该看到支持的网站列表
2. 一个下载任务将被创建，并显示任务ID
3. 您可以监控下载进度，从0%到100%
4. 下载完成后，将打开一个网页显示下载链接
5. 点击下载按钮，您应该能够下载视频文件

## 故障排除

- **端口冲突**: 如果端口8080已被占用，您可以修改`mock_api_server.py`中的端口号
- **Python版本**: 确保使用Python 3.6或更高版本
- **视频创建失败**: 如果`create_sample_video.py`失败，模拟服务器将创建一个简单的二进制文件

## 测试后清理

测试完成后，使用以下命令退出虚拟环境:

```bash
deactivate
```