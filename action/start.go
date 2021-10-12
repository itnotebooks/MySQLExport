// Copyright 2021 EricWinn
// Author:   Eric Winn
// Email:    eng.eric.winn@gmail.com
// Time:     2021/10/12 18:02
// File:     start.py
// Software: GoLand

package action

import (
	"MySQLExport/build"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

var AppFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    "config",
		Aliases: []string{"c"},
		Value:   "./config.yaml",
		Usage:   "配置文件",
	},
	//&cli.StringFlag{
	//	Name:    "type",
	//	Aliases: []string{"t"},
	//	Value:   "csv",
	//	Usage:   "输出文件类型，取值范围：csv | excel",
	//},
}

func Start() {
	app := &cli.App{
		Name:        "MySQLExport",
		Version:     build.Version(),
		Description: "MySQL数据导出工具，如需支持请联系Eric（eng.eric.winn@gmail.com）",
		Flags:       AppFlags,
		Action: func(c *cli.Context) error {
			return Factory(c)
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
