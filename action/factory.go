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
)

var (
	UploadFiles = config.UploadFiles
)

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
		zipFile := ArchiveCsv2ZipFile(target+"/", globalConfig.Archive.ZipFileName, globalConfig.Archive.Encrypt,
			globalConfig.Archive.PassWord,
			globalConfig.Archive.EncryptionMethod)
		UploadFiles = append(UploadFiles, zipFile)
	}

	config.WG.Add(1)
	go UploadFile2SFTPs()

	config.SftpWG.Wait()
	config.WG.Wait()

	if globalConfig.Archive.Enable && globalConfig.Archive.Encrypt {
		log.Println("Zip file password:", globalConfig.Archive.PassWord)
	}
	return nil
}

// ExecQueryAndWriteToCSV 执行构架SQL并将结果写入到CSV文件
func ExecQueryAndWriteToCSV(sql, target, name string) {
	fileName := fmt.Sprintf("%v/%v", target, name)
	Db := model.DB()

	// 执行sql语句
	log.Printf("[%v] SQL execing...\n", name)
	rows, err := Db.Raw(sql).Rows()
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// 调用WriteFile方法，将结果写入CSV文件
	log.Printf("[%v] CSV writing...\n", name)
	err = csv.WriteFile(fileName, rows)
	if err != nil {
		log.Fatal(err)
	}

	config.WG.Done()
	if !config.GlobalConfig.Archive.Enable {
		UploadFiles = append(UploadFiles, fileName)
	}
	return
}

func ArchiveCsv2ZipFile(src, zipfileName string, encrypt bool, password, enc string) string {
	// 将zip包名中的日期符号转换为实际数值
	zipfileName = tools.ConvertDateSymbolToString(zipfileName)

	zipFilePath := src + "/" + zipfileName
	err := zip.ZipLib(zipFilePath, src, encrypt, password, enc)
	if err != nil {
		log.Fatal("Archive error,", err)
	}
	return zipFilePath
}

// UploadFile2SFTPs 文件上传到SFTPs
func UploadFile2SFTPs() {

	var globalConfig = config.GlobalConfig

	for _, server := range globalConfig.Uploads {
		if server.Enable {
			switch server.Engine {
			case "sftp":
				config.SftpWG.Add(1)
				go UploadToSftp(fmt.Sprintf("%v:%v", server.Host, server.Port), server.User,
					server.Password, server.TargetDir)
			}
		}
	}
	config.SftpWG.Wait()
	config.WG.Done()
}

func UploadToSftp(host, user, password, dest string) {

	// 创建SFTP连接
	sftpClient := model.NewSFTP(host, user, password, dest)

	for _, f := range UploadFiles {
		fileName := filepath.Base(f)
		// 上传文件
		sftpClient.FileUpload(f, fileName)
	}
	// 关闭链接
	defer sftpClient.Close()
	config.SftpWG.Done()
}
