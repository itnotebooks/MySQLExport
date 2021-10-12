// Copyright 2021 EricWinn
// Author:   Eric Winn
// Email:    eng.eric.winn@gmail.com
// Time:     2021/10/12 18:03
// File:     build.py
// Software: GoLand

package build

import "strings"

const (
	// VersionNumber 版本号
	VersionNumber = "1.0.0"
)

// Version 生成版本信息
func Version() string {
	var buf strings.Builder
	buf.WriteString(VersionNumber)
	return buf.String()
}
