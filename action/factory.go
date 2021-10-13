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
	"MySQLExport/tools/zip"
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
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
		configFile = fmt.Sprintf("%v/%v", tools.GetBaseDir(), strings.Split(configFile, configFile[0:2])[1])
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
		go ExecQueryAndWriteToCSV(query.SQL, target, query.FileName)
	}

	config.WG.Wait()
	// 压缩上传
	if globalConfig.Archive.Enable {
		password := globalConfig.Archive.PassWord
		zipFile := ArchiveCsv2ZipFile(target+"/", globalConfig.Archive.ZipFileName, password)
		UploadZipFile2SFTP(zipFile)
		config.WG.Wait()
		log.Println("Zip file password:", password)
	}

	return nil
}

// ExecQueryAndWriteToCSV 执行构架SQL并将结果写入到CSV文件
func ExecQueryAndWriteToCSV(sql, target, name string) error {
	fileName := fmt.Sprintf("%v/%v", target, name)
	Db := model.DB()

	// 执行sql语句
	log.Printf("[%v] SQL execing...\n", name)
	rows, err := Db.Raw(sql).Rows()
	if err != nil {
		config.WG.Done()
		log.Fatal(err)
	}
	defer rows.Close()

	// 调用WriteFile方法，将结果写入CSV文件
	log.Printf("[%v] CSV writing...\n", name)
	err = csv.WriteFile(fileName, rows)
	if err != nil {
		config.WG.Done()
		log.Fatal(err)
	}

	config.WG.Done()

	return nil
}

func ArchiveCsv2ZipFile(src, zipfileName, password string) string {
	// 替换zip文件名
	nowTime := time.Now()
	zipfileName = strings.ReplaceAll(zipfileName, "YYYY", nowTime.Format("2006"))
	// MM代表月份mm
	zipfileName = strings.ReplaceAll(zipfileName, "MM", nowTime.Format("01"))
	zipfileName = strings.ReplaceAll(zipfileName, "DD", nowTime.Format("02"))
	zipfileName = strings.ReplaceAll(zipfileName, "HH", nowTime.Format("15"))
	// FF代表分钟MM
	zipfileName = strings.ReplaceAll(zipfileName, "FF", nowTime.Format("04"))
	zipfileName = strings.ReplaceAll(zipfileName, "SS", nowTime.Format("05"))
	zipFilePath := src + "/" + zipfileName
	err := zip.ZipLib(zipFilePath, src, password)
	if err != nil {
		log.Fatal("Archive error,", err)
	}
	return zipFilePath
}

// UploadZipFile2SFTP zip文件上传到SFTP
func UploadZipFile2SFTP(src string) {

	var globalConfig = config.GlobalConfig

	for _, server := range globalConfig.Uploads {
		if server.Enable {
			switch server.Engine {
			case "sftp":
				config.WG.Add(1)
				go UploadToSftp(src, fmt.Sprintf("%v:%v", server.Host, server.Port), server.User,
					server.Password, server.TargetDir)
			}
		}
	}

}

func UploadToSftp(src, host, user, password, dest string) {
	fileName := filepath.Base(src)

	// 创建SFTP连接
	sftpClient := model.NewSFTP(host, user, password, dest)

	// 上传文件
	sftpClient.FileUpload(src, fileName)

	// 关闭链接
	defer sftpClient.Close()
	config.WG.Done()
}
