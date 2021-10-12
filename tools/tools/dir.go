// Copyright 2021 EricWinn
// Author:   Eric Winn
// Email:    eng.eric.winn@gmail.com
// Time:     2021/10/12 17:10
// File:     dir.py
// Software: GoLand

package tools

import (
	"MySQLExport/config"
	"fmt"
	"log"
	"os"
	"time"
)

// MakeDir 递归创建目录: ${baseDir}/target/YYYY-mm-dd_HHMMSS
func MakeDir() string {

	targetDir := fmt.Sprintf("%v/target/%v", config.BaseDir, time.Now().Format("2006-01-02_150405"))

	// 判断目录是否存在,不存在则创建
	_, err := os.Stat(targetDir)
	if err != nil {
		// 创建目录
		err = os.MkdirAll(targetDir, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}

	return targetDir
}
