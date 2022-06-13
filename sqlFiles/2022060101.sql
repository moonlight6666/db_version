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