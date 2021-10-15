// Copyright 2021 EricWinn
// Author:   Eric Winn
// Email:    eng.eric.winn@gmail.com
// Time:     2021/10/12 21:58
// File:     sftp.py
// Software: GoLand

package model

import (
	"MySQLExport/tools/tools"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"os"
	"time"
)

type SFTP struct {
	Host     string
	User     string
	PassWord string
	Target   string
	SSHChn   *ssh.Client
	SFTPChn  *sftp.Client
}

// ClientInit 初始化SFTP连接
func (s *SFTP) ClientInit() {
	var err error
	sshConfig := &ssh.ClientConfig{
		User: s.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(s.PassWord),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		ClientVersion:   "",
		Timeout:         10 * time.Second,
	}

	// 建立连接
	s.SSHChn, err = ssh.Dial("tcp", s.Host, sshConfig)
	if err != nil {
		log.Fatalln(err.Error())
	}

	// 创建SFTP客户端
	s.SFTPChn, err = sftp.NewClient(s.SSHChn)
	if err != nil {
		log.Fatal(err)
	}

}

func (s *SFTP) FileUpload(src, fileName string) {
	var err error
	//// 获取当前目录
	//cwd, err := s.SFTPChn.Getwd()
	//if err != nil {
	//	log.Fatal(err)
	//}

	remoteTargetDir := s.Target
	// 判断远程目录是否需要按日期生成
	remoteTargetDir = tools.ConvertDateSymbolToString(remoteTargetDir)

	// 不存在则创建，存在则不做动作
	err = s.SFTPChn.MkdirAll(remoteTargetDir)

	if err != nil {
		log.Fatalf("%v: %v", remoteTargetDir, err)
	}

	// 上传文件到远程目录
	remoteFilePath := remoteTargetDir + "/" + fileName
	log.Printf("[%v]\tStaring uploading...", src)
	// Create方法，如果远程文件存在，则替换，不存在则创建
	remoteFile, err := s.SFTPChn.Create(remoteFilePath)
	if err != nil {
		log.Fatalln(err)
	}
	defer remoteFile.Close()

	// 获取本地文件对象
	localFilePath := src
	localFile, err := os.Open(localFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer localFile.Close()

	// 将本地文件对象拷贝到远程
	n, err := io.Copy(remoteFile, localFile)
	if err != nil {
		log.Fatalln(err.Error())
	}

	// 获取本地文件大小
	localFileInfo, err := os.Stat(localFilePath)
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Printf("[%v(%v) ---> %s(%v)] : Files upload success!!!", localFilePath, tools.FormatFileSize(localFileInfo.Size()),
		remoteFilePath,
		tools.FormatFileSize(n))

}

func (s *SFTP) Close() {
	s.SSHChn.Close()
	s.SFTPChn.Close()
	log.Printf("[%v] SFTP connection closed", s.Host)
}

// NewSFTP  创建一个新的SFTP客户端
func NewSFTP(host, user, password, dest string) *SFTP {
	s := &SFTP{
		Host:     host,
		User:     user,
		PassWord: password,
		Target:   dest,
	}
	s.ClientInit()

	return s
}
