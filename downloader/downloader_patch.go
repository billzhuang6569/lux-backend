package downloader

// 添加进度回调相关的代码

// ProgressCallback 是进度回调函数类型
type ProgressCallback func(percent int)

// SetProgressCallback 设置进度回调函数
func (downloader *Downloader) SetProgressCallback(callback ProgressCallback) {
	downloader.progressCallback = callback
}

// 修改 Downloader 结构体，添加 progressCallback 字段
type Downloader struct {
	bar              *pb.ProgressBar
	option           Options
	progressCallback ProgressCallback
}

// 重写 writeFile 方法，添加进度回调支持
func (downloader *Downloader) writeFile(url string, file *os.File, headers map[string]string) (int64, error) {
	res, err := request.Request(http.MethodGet, url, nil, headers)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close() // nolint

	// 如果没有进度条，创建一个
	if downloader.bar == nil && res.ContentLength > 0 {
		downloader.bar = progressBar(res.ContentLength)
		if !downloader.option.Silent {
			downloader.bar.Start()
		}
	}

	barWriter := downloader.bar.NewProxyWriter(file)
	// Note that io.Copy reads 32kb(maximum) from input and writes them to output, then repeats.
	// So don't worry about memory.
	written, copyErr := io.Copy(barWriter, res.Body)
	if copyErr != nil && copyErr != io.EOF {
		return written, errors.Errorf("file copy error: %s", copyErr)
	}

	// 调用进度回调
	if downloader.progressCallback != nil && downloader.bar != nil {
		percent := int(float64(downloader.bar.Current()) * 100 / float64(downloader.bar.Total()))
		downloader.progressCallback(percent)
	}

	return written, nil
}