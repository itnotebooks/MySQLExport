// Copyright 2021 EricWinn
// Author:   Eric Winn
// Email:    eng.eric.winn@gmail.com
// Time:     2021/10/12 18:03
// File:     const.py
// Software: GoLand

package config

import (
	"sync"
)

var UploadFiles []string
var WG sync.WaitGroup
var SftpWG sync.WaitGroup

var FileType = []string{
	"csv",
	"excel",
}
