#!/usr/bin/env python3
# -*- coding: utf-8 -*-

"""
创建一个简单的样本视频文件用于测试
这个脚本需要安装 moviepy 库：pip install moviepy
"""

import os
import sys
from moviepy.editor import ColorClip, TextClip, CompositeVideoClip

def create_sample_video(output_path="sample.mp4", duration=5):
    """创建一个简单的示例视频"""
    try:
        # 创建背景
        size = (640, 480)
        background = ColorClip(size, color=(0, 0, 0), duration=duration)
        
        # 创建文本
        text = TextClip("Lux API测试视频", fontsize=70, color="white", duration=duration)
        text = text.set_position("center")
        
        # 组合所有元素
        video = CompositeVideoClip([background, text])
        
        # 写入文件
        video.write_videofile(output_path, fps=24, codec="libx264", audio=False)
        
        print(f"示例视频已创建: {output_path}")
        return True
    except Exception as e:
        print(f"创建视频失败: {e}")
        return False

if __name__ == "__main__":
    output_path = "sample.mp4"
    if len(sys.argv) > 1:
        output_path = sys.argv[1]
    
    create_sample_video(output_path)