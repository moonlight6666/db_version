[![Build](https://github.com/logerror/dbupgrade/actions/workflows/maven.yml/badge.svg?branch=release&event=push)](https://github.com/logerror/dbupgrade/actions/workflows/maven.yml)
## 简介

数据库版本控制工具，用于数据库升级和维护，只需编写.sql文件后，自动进行sql的更替。支持DML, DDL语句以及存储过程。

## 安装
    go get github.com/moonlight6666/db_version

## 脚本执行
```shell
go run cmd/main.go -h 数据库地址 -u 数据库用户名 -p 密码 -n 数据库名 -c 操作
```
## 代码调用
初始化
```go
dbVersion, err := NewDbVersion("user", "password", "host", "3306", "dbname", "sql文件目录")
if err != nil {
	panic(err)
}
defer dbVersion.db.Close()
err = dbVersion.Init()
```

更新数据库版本
```go
dbVersion.Update()
```

获取数据库版本
```go
verion, _ := dbVersion.Version()
```

删除数据库
```go
dbVersion.Drop()
```