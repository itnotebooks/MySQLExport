// Copyright 2021 EricWinn
// Author:   Eric Winn
// Email:    eng.eric.winn@gmail.com
// Time:     2021/10/13 13:25
// File:     tools.py
// Software: GoLand

package tools

import (
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func GetBaseDir() string {
	baseDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return baseDir
}

// RandomString 生成随机密码
func RandomString(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func ConvertDateSymbolToString(s string) string {
	nowTime := time.Now()
	s = strings.ReplaceAll(s, "YYYY", nowTime.Format("2006"))
	// MM代表月份mm
	s = strings.ReplaceAll(s, "MM", nowTime.Format("01"))
	s = strings.ReplaceAll(s, "DD", nowTime.Format("02"))
	s = strings.ReplaceAll(s, "HH", nowTime.Format("15"))
	// FF代表分钟MM
	s = strings.ReplaceAll(s, "FF", nowTime.Format("04"))
	s = strings.ReplaceAll(s, "SS", nowTime.Format("05"))
	return s
}
