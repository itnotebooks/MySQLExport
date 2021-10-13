// Copyright 2021 EricWinn
// Author:   Eric Winn
// Email:    eng.eric.winn@gmail.com
// Time:     2021/10/13 12:11
// File:     zip.py
// Software: GoLand

package zip

import (
	"github.com/itnotebooks/zip"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// ZipLib 压缩递归压缩
func ZipLib(dst, src, password string) (err error) {
	var dstFileBaseName = ""
	// 创建压缩文件对象
	zfile, err := os.Create(dst)
	defer zfile.Close()

	if err != nil {
		return err
	}

	// 通过文件对象生成写入对象
	zFileWriter := zip.NewWriter(zfile)
	defer func() {
		// 检测一下是否成功关闭
		if err := zFileWriter.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	info, err := os.Stat(src)
	if err != nil {
		return err
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(src)

		if !strings.HasSuffix(src, "/") {
			dstName := filepath.Base(dst)
			dstFileBaseName = strings.TrimSuffix(dstName, filepath.Ext(dstName))
		}
	}

	// 将文件写入 zFileWriter 对象 ，可能会有很多个目录及文件，递归处理
	return filepath.Walk(src, func(path string, fi os.FileInfo, errBack error) (err error) {
		if errBack != nil {
			return errBack
		}

		if strings.HasSuffix(path, ".zip") {
			return
		}
		//创建文件头
		header, err := zip.FileInfoHeader(fi)
		if err != nil {
			return err
		}

		if baseDir != "" {
			// 如果原目录是以"/"结尾，表示打包指定目录时不包含该目录
			if strings.HasSuffix(src, "/") {
				header.Name = strings.TrimPrefix(path, src)
			} else {
				header.Name = filepath.Join(dstFileBaseName, baseDir, strings.TrimPrefix(path, src))
			}
		}

		if fi.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate

		}

		// 设置密码
		header.SetPassword(password)

		// 写入文件头信息
		fh, err := zFileWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		// 判断是否是标准文件
		if !header.Mode().IsRegular() {
			return nil
		}

		// 打开要压缩的文件
		file, err := os.Open(path)
		defer file.Close()
		if err != nil {
			return
		}

		// 将文件对象拷贝到 writer 结构中
		ret, err := io.Copy(fh, file)
		if err != nil {
			return
		}

		log.Printf("added： %s, total: %d\n", path, ret)

		return nil
	})
}
