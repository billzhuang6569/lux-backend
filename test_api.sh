#!/bin/bash

# API测试脚本
SERVER_URL="http://localhost:8080"

# 测试函数 - 格式化输出测试结果
function test_result() {
  if [ $1 -eq 0 ]; then
    echo -e "\e[32m测试通过: $2\e[0m"
  else
    echo -e "\e[31m测试失败: $2\e[0m"
    echo "错误信息: $3"
  fi
}

echo "===== Lux API 测试 ====="
echo

# 测试 1: 获取支持的网站列表
echo "测试 1: 获取支持的网站列表"
SITES_RESPONSE=$(curl -s -X GET $SERVER_URL/api/supported-sites)
if [[ $SITES_RESPONSE == *"sites"* ]]; then
  test_result 0 "支持的网站列表API" "$SITES_RESPONSE"
else
  test_result 1 "支持的网站列表API" "$SITES_RESPONSE"
fi
echo

# 测试 2: 创建下载任务
echo "测试 2: 创建下载任务"
DOWNLOAD_RESPONSE=$(curl -s -X POST $SERVER_URL/api/download \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://www.bilibili.com/video/BV1MN41127yH/?spm_id_from=333.999.0.0",
    "format": "mp4",
    "quality": "best"
  }')

if [[ $DOWNLOAD_RESPONSE == *"taskId"* ]]; then
  # 提取任务ID
  TASK_ID=$(echo $DOWNLOAD_RESPONSE | sed 's/.*"taskId":"\([^"]*\)".*/\1/')
  test_result 0 "创建下载任务" "$DOWNLOAD_RESPONSE"
else
  test_result 1 "创建下载任务" "$DOWNLOAD_RESPONSE"
  exit 1
fi
echo

# 测试 3: 获取下载状态
echo "测试 3: 获取下载状态"
echo "任务ID: $TASK_ID"
echo "正在检查下载进度..."

# 轮询任务状态，直到完成或出错
MAX_CHECKS=30
for (( i=1; i<=$MAX_CHECKS; i++ )); do
  STATUS_RESPONSE=$(curl -s -X GET $SERVER_URL/api/status/$TASK_ID)
  
  # 提取状态
  TASK_STATUS=$(echo $STATUS_RESPONSE | sed 's/.*"status":"\([^"]*\)".*/\1/')
  TASK_PROGRESS=$(echo $STATUS_RESPONSE | sed 's/.*"progress":\([^,}]*\).*/\1/')
  
  echo "状态检查 $i: 状态=$TASK_STATUS, 进度=$TASK_PROGRESS%"
  
  if [[ $TASK_STATUS == "completed" ]]; then
    DOWNLOAD_URL=$(echo $STATUS_RESPONSE | sed 's/.*"downloadUrl":"\([^"]*\)".*/\1/')
    test_result 0 "下载完成" "下载URL: $DOWNLOAD_URL"
    break
  elif [[ $TASK_STATUS == "error" ]]; then
    ERROR_MSG=$(echo $STATUS_RESPONSE | sed 's/.*"message":"\([^"]*\)".*/\1/')
    test_result 1 "下载任务" "错误: $ERROR_MSG"
    break
  fi
  
  # 等待5秒再次检查
  if [ $i -lt $MAX_CHECKS ]; then
    sleep 5
  fi
done

# 如果达到最大检查次数仍未完成
if [ $i -gt $MAX_CHECKS ]; then
  test_result 1 "下载任务" "超时 - 任务未在预期时间内完成"
fi

echo
echo "===== 测试完成 ====="