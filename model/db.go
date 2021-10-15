// Copyright 2021 EricWinn
// Author:   Eric Winn
// Email:    eng.eric.winn@gmail.com
// Time:     2021/10/11 17:08
// File:     db.py
// Software: GoLand

package model

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
	Db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		globalConfig.MySQL.MySQLUser,
		globalConfig.MySQL.MySQLPassword,
		globalConfig.MySQL.MySQLHost,
		globalConfig.MySQL.MySQLPort,
		globalConfig.MySQL.MySQLDb),
	)
	if err != nil {
		log.Println(err)
		log.Fatal("MySQL连接建立失败!!!")
	}

	db = Db
	db.DB().SetConnMaxLifetime(time.Minute * 10)
	db.DB().SetMaxOpenConns(50)
	db.DB().SetMaxIdleConns(15)
}

func DB() *gorm.DB {
	db.LogMode(config.GlobalConfig.DEBUG)
	return db.New()
}
