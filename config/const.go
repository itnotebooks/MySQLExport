// Copyright 2021 EricWinn
// Author:   Eric Winn
// Email:    eng.eric.winn@gmail.com
// Time:     2021/10/12 18:03
// File:     const.py
// Software: GoLand

package config

import (
	"os"
	"path/filepath"
	"sync"
)

var WG sync.WaitGroup

// 程序根目录
var BaseDir, _ = filepath.Abs(filepath.Dir(os.Args[0]))

var FileType = []string{
	"csv",
	"excel",
}

