package main

import (
	"db_version"
	"flag"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"time"
)

var (
	dbHost     string
	dbPort     int
	dbUser     string
	dbPassword string
	dbName     string
	sqlDir     string
	action     string
)

func main() {
	flag.StringVar(&dbHost, "h", "127.0.0.1", "数据库地址")
	flag.IntVar(&dbPort, "port", 3306, "数据库端口")
	flag.StringVar(&dbUser, "u", "root", "数据库账号")
	flag.StringVar(&dbPassword, "p", "", "数据库密码")
	flag.StringVar(&dbName, "n", "", "数据库名字")
	flag.StringVar(&sqlDir, "d", "./", "sql文件目录")

	flag.StringVar(&action, "c", "", "操作: update|version|drop")

	flag.Parse()

	dbVersion, err := db_version.NewDbVersion(dbUser, dbPassword, dbHost, dbPort, dbName, sqlDir)
	if err != nil {
		panic(errors.Wrap(err, "连接数据库失败"))
	}
	defer dbVersion.Close()

	err = dbVersion.Init()
	if err != nil {
		panic(errors.Wrap(err, "初始化失败"))
	}

	t1 := time.Now()
	switch action {
	case "update":
		err = dbVersion.Update()
		if err != nil {
			panic(err)
		}
	case "drop":
		err = dbVersion.Drop()
		if err != nil {
			panic(err)
		}
	case "version":
		nowDbVersion, err := dbVersion.Version()
		if err != nil {
			panic(err)
		}
		fmt.Println("当前数据库版本:", nowDbVersion)
	default:
		fmt.Println("参数错误:", action)
		os.Exit(1)
	}
	cost := time.Since(t1)
	fmt.Println("耗时:", cost)
}
