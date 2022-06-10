package db_version

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const ext = ".sql"

type (
	DbVersion struct {
		DbName string
		SqlDir string
		Db     *sql.DB
	}
)

func NewDbVersion(dbUser string, dbPasswd string, dbHost string, dbPort int, dbName string, sqlDir string) (*DbVersion, error) {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", dbUser, dbPasswd, dbHost, dbPort, "")

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	dbVersion := &DbVersion{
		DbName: dbName,
		SqlDir: sqlDir,
		Db:     db,
	}
	return dbVersion, nil
}

func (d *DbVersion) Close() error {
	return d.Db.Close()
}

// Init 初始化
func (d *DbVersion) Init() error {
	rows := d.Db.QueryRow("SELECT `SCHEMA_NAME` FROM information_schema.SCHEMATA WHERE `SCHEMA_NAME` = '" + d.DbName + "';")
	name := ""
	_ = rows.Scan(&name)
	//if err != nil {
	//	return err
	//}
	needInit := name == ""

	if needInit {
		fmt.Println("创建数据库: ", d.DbName)
		_, err := d.Db.Exec("CREATE DATABASE `" + d.DbName + "` CHARACTER SET 'utf8' COLLATE 'utf8_general_ci';\n")
		if err != nil {
			return errors.Wrap(err, "创建数据库失败")
		}

		_, err = d.Db.Exec("use `" + d.DbName + "` ;")
		if err != nil {
			return errors.Wrap(err, "切换数据库失败")
		}

		fmt.Println("创建表: db_version")
		_, err = d.Db.Exec("CREATE TABLE `db_version` ( `version` INT, PRIMARY KEY ( `version`));\n")

		if err != nil {
			return errors.Wrap(err, "创建表 db_version 失败")
		}

		_, err = d.Db.Exec("INSERT INTO `db_version` VALUES (0);\n")
		if err != nil {
			return errors.Wrap(err, "初始化表 db_version 失败")
		}

	}

	_, err := d.Db.Exec("use `" + d.DbName + "` ;")
	if err != nil {
		return errors.Wrap(err, "切换数据库失败")
	}

	return nil
}

// Update 更新数据库版本
func (d *DbVersion) Update() error {
	nowDbVersion, err := d.Version()
	if err != nil {
		return err
	}

	fmt.Println("更新前数据库版本:", nowDbVersion)

	files, _ := d.getSQLFileList(d.SqlDir, ext)
	sort.Strings(files)
	for _, v := range files {
		err := d.execSqlFile(v, nowDbVersion)
		if err != nil {
			return err
		}
	}

	newDbVersion, err := d.Version()
	if err != nil {
		return err
	}

	fmt.Println("更新后数据库版本:", newDbVersion)

	return nil
}

// 执行一个SQL文件到数据库
func (d *DbVersion) execSqlFile(fileName string, nowDbVersion int) error {
	thisSqlFileVersion, err := d.fileName2Version(fileName)
	if err != nil {
		return err
	}

	if thisSqlFileVersion <= nowDbVersion {
		return nil
	}

	baseName := filepath.Base(fileName)
	fmt.Print("执行文件: ", baseName, " ", strings.Repeat(".", 50-len(baseName)))

	file, err := os.Open(fileName)
	if err != nil {
		return errors.Wrap(err, "打开文件失败")
	}

	context, err := ioutil.ReadAll(file)
	if err != nil {
		return errors.Wrap(err, "读取文件失败")
	}

	sqlList := strings.Split(string(context), ";")
	for _, s := range sqlList {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		_, err = d.Db.Exec(s)
		if err != nil {
			fmt.Println("\nFail SQL:", s)
			if err != nil {
				return errors.Wrap(err, "SQL执行失败")
			}
		}
	}

	err = d.setVersion(thisSqlFileVersion)
	if err != nil {
		return err
	}

	fmt.Println(" [OK]")
	return nil
}

// 设置数据库版本
func (d *DbVersion) setVersion(v int) error {
	_, err := d.Db.Exec("UPDATE `db_version` SET `version` = " + strconv.Itoa(v))
	if err != nil {
		return errors.Wrap(err, "设置数据库版本失败")
	}

	return nil
}

// Drop 删除数据库
func (d *DbVersion) Drop() error {
	fmt.Println("删除数据库:", d.DbName)

	_, err := d.Db.Exec("drop database `" + d.DbName + "`;")
	if err != nil {
		return errors.Wrap(err, "删除数据库失败")
	}

	return nil
}

// Version 获取数据库版本
func (d *DbVersion) Version() (int, error) {
	rows := d.Db.QueryRow("SELECT `version` FROM `db_version`")

	dbVersion := -1
	err := rows.Scan(&dbVersion)

	return dbVersion, err
}

// 文件名 转换成版本号
func (d *DbVersion) fileName2Version(fileName string) (int, error) {
	baseName := filepath.Base(fileName)
	versionStr := strings.Split(baseName, ".")[0]
	versionInt, err := strconv.Atoi(versionStr)
	if err != nil {
		return 0, errors.Wrap(err, "文件名转换版本号失败:"+fileName)
	}
	return versionInt, nil
}

// 获取文件夹下的所有SQL文件
func (d *DbVersion) getSQLFileList(dirPth, suffix string) (files []string, err error) {
	files = make([]string, 0, 50)
	suffix = strings.ToUpper(suffix)

	err = filepath.Walk(dirPth, func(filename string, fi os.FileInfo, err error) error {
		//遍历目录
		if fi.IsDir() {
			// 忽略目录
			return nil
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
			baseName := filepath.Base(filename)
			if len(baseName) == 14 {
				files = append(files, filename)
			} else {
				fmt.Println("忽略文件:", filename)
			}
		}
		return nil
	})
	return files, err
}
