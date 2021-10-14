# MySQLExport

[![LICENSE](https://img.shields.io/badge/license-Apache%20v2-blue.svg)](https://github.com/itnotebooks/MySQLExport/blob/master/LICENSE)
[![LICENSE](https://img.shields.io/badge/version-Go%20v1.16-blue.svg)](https://golang.org/)
[![GitHub Releases (latest SemVer)](https://img.shields.io/github/v/release/itnotebooks/MySQLExport)](https://github.com/itnotebooks/MySQLExport/releases/latest)
[![GitHub All Releases](https://img.shields.io/github/downloads/itnotebooks/MySQLExport/total)](https://github.com/itnotebooks/MySQLExport/releases)

MySQLExport工具会执行指定的一条或多条SQL，并将查询结果生成CSV文件上传到指定存储空间中

## 功能

-   SQL 查询自动生成CSV文件
-   CSV 文件压缩
-   CSV 文件压缩加密，支持以下加密算法
    -   Standard         ZIP标准，安全性最低
    -   AES128
    -   AES192
    -   AES256           （默认）
-   结果上传到共享储存空间
    -   SFTP
    -   OSS（开发中...）

## 编译

进入项目根目录，执行以下命令
```bash
go build
```

## 配置

```yaml
# 是否开启Debug模式，默认false(关闭)
debug: false

# MySQL数据库连接信息
mysql:
  enable: true
  host:  127.0.0.1
  port: 3306
  user: root
  password:
  database:

# 将结果上传到指定的服务器上
uploads:
  # 上传到SFTP服务器
  - enable: true
    engine: sftp
    host:
    port: 22
    user:
    password:
    # 接收日期变量（分钟需要用 FF 代替），只支持24小时制，变量字符全双位大写
    # YYYY-MM 2021-03
    # YYYY-MM-DD 2021-03-15
    # YYYY-MM-DD_HHFF 2021-03-15_1342
    # YYYY-MM-DD_HHFFSS 2021-03-15_134233
    # 特殊格式可能会导致文件创建失败，请按OS系统的全名规则创建
    # 注：不支持特殊字符使用双引号括起的做法
    target: /tmp/YYYYMM

# 只接受所有CSV压缩到一个包中
archive:
  # 是否对CSV文件进行压缩
  enable: false
  # 是否需要对压缩包进行加密，默认为false（关闭）
  encrypt: false
  # 加密算法，支持以下加密方式
  # Standard         ZIP标准，安全性最低
  # AES128           AES128位，安全性高
  # AES192           AES192位，安全性高
  # AES256           AES256位，安全性最高，本程序默认采用此加密方式
  encryptionMethod: AES256
  # 密码如果为空，将会自动生成一个12位的随机密码，随程序执行完成后打印出来
  password:
  # 接收日期变量（分钟需要用 FF 代替），只支持24小时制，变量字符全双位大写
  # YYYY-MM 2021-03
  # YYYY-MM-DD 2021-03-15
  # YYYY-MM-DD_HHFF 2021-03-15_1342
  # YYYY-MM-DD_HHFFSS 2021-03-15_134233
  # 特殊格式可能会导致文件创建失败，请按OS系统的全名规则创建
  # 注：不支持特殊字符使用双引号括起的做法
  zipFile: YYYYMM.zip

# 需要执行的语句及存入的Sheet页名称
# 如果SQL语句分行，请使用双引号包起来
queries:
  - sql: SELECT * FROM TABLES
    # 文件名以.csv结尾，需要注意带文件后缀，Windows平台可能会有影响
    fileName: XXXX.csv
```

## 运行

执行成功后会在程序所在目录中创建 target目录，保存所有CSV文件
```shell
./MySQLExport -c ./config.yaml
```

