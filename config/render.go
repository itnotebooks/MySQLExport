// Copyright 2021 EricWinn
// Author:   Eric Winn
// Email:    eng.eric.winn@gmail.com
// Time:     2021/10/12 18:04
// File:     render.py
// Software: GoLand

package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

// MySQLField MySQL数据连接信息
type MySQLField struct {
	MySQLEnable   bool   `yaml:"enable"`
	MySQLHost     string `yaml:"host"`
	MySQLPort     int    `yaml:"port"`
	MySQLUser     string `yaml:"user"`
	MySQLPassword string `yaml:"password"`
	MySQLDb       string `yaml:"database"`
}

// QueryField 要执行的SQL及结果要存入的Sheet页
type QueryField struct {
	SQL      string `yaml:"sql" json:"sql"`
	FileName string `yaml:"fileName" json:"fileName"`
}

// ConfigField 一级配置文件
type ConfigField struct {
	MySQL   MySQLField   `yaml:"mysql" json:"mysql"`
	Queries []QueryField `yaml:"queries" json:"queries"`
}

// GlobalConfig 配置变更存放于全局变量
var GlobalConfig ConfigField

func RenderConfig(c string) error {

	var err error
	var config ConfigField

	// 判断是否为文件并读取文件内容
	f, err := ioutil.ReadFile(c)
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(f, &config)
	if err != nil {
		log.Fatal(err)
	}

	// MySQL配置信息
	if config.MySQL.MySQLEnable {
		GlobalConfig.MySQL = config.MySQL
	}

	// Query语句
	GlobalConfig.Queries = config.Queries

	return nil
}
