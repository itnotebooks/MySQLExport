// Copyright 2021 EricWinn
// Author:   Eric Winn
// Email:    eng.eric.winn@gmail.com
// Time:     2021/10/11 17:17
// File:     global.py
// Software: GoLand

package model

import "time"

type mysql struct {
	Host     string
	User     string
	Password string
	Db       string
	Port     string
}

type general struct {
	SecretKey string
	Host      string
	Hours     time.Duration
	GrpcAddr  string
}

type DbInfo struct {
	Host     string
	User     string
	Password string
	Port     string
	Db       string
}

type Config struct {
	General general
	Mysql   mysql
}

var C Config
