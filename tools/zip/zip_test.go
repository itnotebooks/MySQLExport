// Copyright 2021 EricWinn
// Author:   Eric Winn
// Email:    eng.eric.winn@gmail.com
// Time:     2021/10/13 13:38
// File:     zip_test.py
// Software: GoLand

package zip

import (
	"MySQLExport/tools/tools"
	"log"
	"testing"
)

func TestZipLib(t *testing.T) {
	pass := tools.RandomString(12)
	//baseDir := tools.GetBaseDir()
	baseDir := "/Users/eric/workspaces/go/src/MySQLExport"
	ZipLib(baseDir+"/test/202110.zip", baseDir+"/target/2021-10-13_114245/adp_data_config_8.csv", pass)
	log.Println("Zip file password:", pass)
}
