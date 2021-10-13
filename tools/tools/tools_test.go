// Copyright 2021 EricWinn
// Author:   Eric Winn
// Email:    eng.eric.winn@gmail.com
// Time:     2021/10/13 13:28
// File:     tools_test.py
// Software: GoLand

package tools

import (
	"log"
	"testing"
)

func TestRandomString(t *testing.T) {
	pass := RandomString(12)
	log.Println(pass)
}
