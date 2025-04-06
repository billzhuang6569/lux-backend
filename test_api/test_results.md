# Lux API 测试结果示例

以下是API测试的预期输出结果示例。

## 1. 测试支持的网站API

请求：
```
GET /api/supported-sites
```

成功响应：
```json
{
  "sites": [
    "YouTube",
    "Twitter",
    "Bilibili",
    "抖音",
    "优酷",
    "腾讯视频",
    "爱奇艺",
    "Instagram",
    "Facebook",
    "Vimeo",
    "微博",
    "小红书",
    "TikTok",
    "Reddit",
    "PornHub",
    "AcFun",
    "Tumblr",
    "芒果TV",
    "喜马拉雅",
    "映客直播",
    "斗鱼",
    "虎牙",
    "快手"
  ]
}
```

## 2. 测试下载API

请求：
```
POST /api/download
Content-Type: application/json

{
  "url": "https://www.bilibili.com/video/BV1MN41127yH",
  "format": "mp4",
  "quality": "best"
}
```

成功响应：
```json
{
  "taskId": "task_abcdef123456",
  "status": "pending",
  "message": "下载任务已创建",
  "progress": 0
}
```

## 3. 测试状态API

请求：
```
GET /api/status/task_abcdef123456
```

处理中响应：
```json
{
  "taskId": "task_abcdef123456",
  "status": "processing",
  "progress": 45,
  "message": "正在下载中..."
}
```

完成响应：
```json
{
  "taskId": "task_abcdef123456",
  "status": "completed",
  "progress": 100,
  "message": "下载已完成",
  "downloadUrl": "http://localhost:8080/downloads/video_abcdef123456.mp4",
  "filename": "video_abcdef123456.mp4"
}
```

## 4. 测试下载链接

一旦下载完成，可以通过浏览器访问 `downloadUrl` 来下载文件。测试网页会自动打开，显示下载按钮。点击按钮应该开始下载视频文件。

## 总结

如果所有API都返回了预期的响应，并且最终能够下载文件，则测试是成功的。测试脚本将在控制台中输出类似的结果，并显示每个测试步骤的状态。