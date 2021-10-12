// Copyright 2021 EricWinn
// Author:   Eric Winn
// Email:    eng.eric.winn@gmail.com
// Time:     2021/10/12 18:04
// File:     main.py
// Software: GoLand

package main

import "MySQLExport/action"

func main() {
	// 解析配置文件中的信息放入到全局变量中
	//conf := config.GetConfig()
	//config.RenderConfig(conf)

	action.Start()
}
