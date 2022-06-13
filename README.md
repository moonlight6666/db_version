[![Build](https://github.com/logerror/dbupgrade/actions/workflows/maven.yml/badge.svg?branch=release&event=push)](https://github.com/logerror/dbupgrade/actions/workflows/maven.yml)
## 简介

数据库版本控制工具，用于数据库升级和维护，只需编写.sql文件后，自动进行sql的更替。支持DML, DDL语句以及存储过程。

## 安装
    go get github.com/moonlight6666/db_version

## 命令行执行
```shell
go build -o db_version cmd/main.go
./db_version  -h 数据库地址 -u 数据库用户名 -p 密码 -n 数据库名 -d sql文件目录 -c 操作
```
## 代码调用
初始化
```go
dbVersion, err := NewDbVersion("user", "password", "host", "3306", "dbname", "sql文件目录")
if err != nil {
	panic(err)
}
defer dbVersion.Close()
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

## 使用说明
#### 注意: 文件格式必须是:年(4位)月(2位)日(2位)序号(2位).sql
先准备一个文件
sqlFiles/2022060101.sql
```sql
-- 创建user表
CREATE TABLE `user`
(
    `id`      int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `account` varchar(64) NOT NULL COMMENT '账号',
    `name`    varchar(64) NOT NULL COMMENT '名称',
    `mail`    varchar(64) NOT NULL COMMENT '邮件',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT COMMENT='用户';

-- 创建user_log表
CREATE TABLE `user_log`
(
    `id`      int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `account` varchar(64)  NOT NULL COMMENT '账号',
    `action`  varchar(256) NOT NULL COMMENT '操作',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT COMMENT='用户日志';
```
执行
```shell
./db_version  -p password -d "./sqlFiles" -n test -c update
```
输出
```shell
更新前数据库版本: 0
执行文件: 2022060101.sql .................................... [OK]
更新后数据库版本: 2022060101
```
执行后版本号为2022060101
接着再新增一个文件， 对数据库进行结构和数据修改
sqlFiles/2022060102.sql
```sql
-- 修改user_log表
ALTER table user_log add `time` time NOT NULL COMMENT '时间';

-- 创建角色表
CREATE TABLE `role`
(
    `id`   int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `name` varchar(64) NOT NULL COMMENT '角色名称',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT COMMENT='角色';


-- 插入user表数据
INSERT INTO `user`(`account`, `name`, `mail`) VALUES ('test', '测试', 'test@qq.com');

```

执行
```shell
./db_version  -p password -d "./sqlFiles" -n test -c update
```
输出
```shell
更新前数据库版本: 2022060101
执行文件: 2022060102.sql .................................... [OK]
更新后数据库版本: 2022060102
```

数据库版本升级为2022060102


执行
```shell
./db_version  -p password -d "./sqlFiles" -n test -c version
```
输出
```shell
当前数据库版本: 2022060102
```