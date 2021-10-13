// Copyright 2021 EricWinn
// Author:   Eric Winn
// Email:    eng.eric.winn@gmail.com
// Time:     2021/10/11 17:08
// File:     db.py
// Software: GoLand

package model

import "C"
import (
	"MySQLExport/config"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"time"
)

var db *gorm.DB

func DbInit() {
	globalConfig := config.GlobalConfig
	// 创建DB连接
	newDb, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		globalConfig.MySQL.MySQLUser,
		globalConfig.MySQL.MySQLPassword,
		globalConfig.MySQL.MySQLHost,
		globalConfig.MySQL.MySQLPort,
		globalConfig.MySQL.MySQLDb),
	)
	if err != nil {
		log.Fatal("MySQL连接建立失败!!!")
	}

	db = newDb
	Db := db.DB()
	Db.SetConnMaxLifetime(time.Minute * 10)
	Db.SetMaxOpenConns(50)
	Db.SetMaxIdleConns(15)
}

func DB() *gorm.DB {
	db.LogMode(false)
	return db.New()
}
