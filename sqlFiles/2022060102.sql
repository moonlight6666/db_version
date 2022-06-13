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
