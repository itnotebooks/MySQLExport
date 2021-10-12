// Copyright 2021 EricWinn
// Author:   Eric Winn
// Email:    eng.eric.winn@gmail.com
// Time:     2021/10/12 18:01
// File:     factory.py
// Software: GoLand

package action

import (
	"MySQLExport/config"
	"MySQLExport/model"
	"MySQLExport/tools/csv"
	"MySQLExport/tools/tools"
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"strings"
)

type Result struct {
	code           int
	target_content string
}

func Factory(c *cli.Context) error {
	var err error
	configFile := c.String("config")

	// 如果配置文件是相对路径，自动替换为绝对路径
	if !strings.HasPrefix(configFile, "/") {
		configFile = fmt.Sprintf("%v/%v", config.BaseDir, strings.Split(configFile, configFile[0:2])[1])
	}

	// 判断配置文件是否存在
	fileStat, err := os.Stat(configFile)
	if err != nil {
		log.Fatal(err)
	}

	// 判断是否为目录
	if fileStat.IsDir() {
		log.Fatalf("%v is a directory", configFile)
	}

	// 读取配置文件
	err = config.RenderConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}

	// 获取全局配置信息
	globalConfig := config.GlobalConfig

	// 创建文件输出目录: ${baseDir}/target/YYYY-mm-dd_HHMMSS
	target := tools.MakeDir()

	// DB初始化
	model.DbInit()
	for _, query := range globalConfig.Queries {
		config.WG.Add(1)
		go ExecQueryAndWriteToCSV(query.SQL, fmt.Sprintf("%v/%v", target, query.FileName))
	}

	config.WG.Wait()
	return nil
}

// ExecQueryAndWriteToCSV 执行构架SQL并将结果写入到CSV文件
func ExecQueryAndWriteToCSV(sql, fileName string) {
	Db := model.DB()
	// 执行sql语句
	rows, err := Db.Raw(sql).Rows()
	if err != nil {
		config.WG.Done()
		log.Fatal(err)
	}
	defer rows.Close()
	// 调用WriteFile方法，将结果写入CSV文件
	err = csv.WriteFile(fileName, rows)
	if err != nil {
		config.WG.Done()
		log.Fatal(err)
	}

	config.WG.Done()

	return
}
