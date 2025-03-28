-- Active: 1740461445112@@192.168.80.128@3309@rms
-- Active: 1740461445112@@192.168.80.128@3309
CREATE DATABASE problem;
CREATE DATABASE user;
-- 公共题库题目较多，id采用自增主键以便后续复杂的分页处理。
-- 而专题下的题目为了不跟公共题库产生主键冲突，而且业务逻辑存在差异，所以用另外的表存储
-- 题目
CREATE TABLE `problem`.`t_problem`(
    `problem_id` INT(10) NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `title` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '',
    `visible` TINYINT(4) NOT NULL DEFAULT 0,
    `description` MEDIUMTEXT NOT NULL COMMENT '需要encode的富文本',
    `tags` TEXT NOT NULL COMMENT '需要encode',
    `author` VARCHAR(32) NOT NULL DEFAULT '' COMMENT '作者id',
    `std_sql` TEXT NOT NULL COMMENT 'encode后标答',
    `spj` TEXT NOT NULL COMMENT 'special judge',
    `test_data` MEDIUMTEXT NOT NULL COMMENT 'encode后生成测试数据的sql',
    `std_ans` MEDIUMTEXT NOT NULL COMMENT 'encode后标答的输出(优化一下自动计算)',
    `std_cost` BIGINT NOT NULL DEFAULT 0 COMMENT 'std执行时长 unix',
    `status` VARCHAR(64) NOT NULL DEFAULT 'ungenerated' COMMENT '题目生成状态',
    `exist` TINYINT NOT NULL DEFAULT 1 COMMENT '0表示被删除',
    INDEX `idx_visible` (`exist`,`visible`)
)ENGINE=InnoDB;

-- 专题列表
CREATE TABLE `problem`.`t_module`(
    `module_id` INT NOT NULL PRIMARY KEY,
    `title` VARCHAR(64) NOT NULL DEFAULT '',
    `group_id` INT NOT NULL DEFAULT 0 COMMENT '0-所有学生在visible时可见，1-group内学生在visible时可见',
    `visible` TINYINT(4) NOT NULL DEFAULT 0,
    `time_limit` TINYINT(4) NOT NULL DEFAULT 0 COMMENT '0-不限时，1-限制时间begin-end',
    `begin` DATETIME not null DEFAULT now(),
    `end` DATETIME not null DEFAULT now(),
    INDEX `idx_visible` (`visible`),
    INDEX `idx_group_id` (`group_id`)
);

-- 专题题目
CREATE TABLE `problem`.`t_module_problem`(
    `module_id` INT NOT NULL,
    `problem_id` INT(10) NOT NULL,
    `title` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '',
    `visible` TINYINT(4) NOT NULL DEFAULT 0,
    `description` MEDIUMTEXT NOT NULL COMMENT '需要encode的富文本',
    `tags` TEXT NOT NULL COMMENT '需要encode',
    `author` VARCHAR(16) NOT NULL DEFAULT '' COMMENT '作者id',
    `std_sql` TEXT NOT NULL COMMENT 'encode后标答',
    `spj` TEXT NOT NULL COMMENT 'special judge',
    `test_data` MEDIUMTEXT NOT NULL COMMENT 'encode后生成测试数据的sql',
    `std_ans` MEDIUMTEXT NOT NULL COMMENT 'encode后标答的输出(优化一下自动计算)',
    `std_cost` BIGINT NOT NULL DEFAULT 0 COMMENT 'std执行时长 unix',
    `status` VARCHAR(64) NOT NULL DEFAULT 'ungenerated' COMMENT '题目生成状态',
    PRIMARY KEY(`module_id`,`problem_id`),
    INDEX `idx_visible` (`visible`)
)ENGINE=InnoDB;



-- 用户提交历史
CREATE TABLE `problem`.`t_submit_history`(
    `submission_id` BIGINT(20) NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `module_id` INT NOT NULL DEFAULT 0 COMMENT '0表示公共题库，其余正数对应专题id',
    `problem_id` INT(10) NOT NULL,
    `user_id` VARCHAR(16) NOT NULL,
    `submit_time` datetime NOT NULL COMMENT '由程序插入，保持一致性',
    `user_sql` TEXT NOT NULL COMMENT 'encode后sql',
    `status` VARCHAR(64) NOT NULL DEFAULT 'In queue',
    `time_cost` BIGINT NOT NULL DEFAULT 0 COMMENT '评测耗时',
    `visible` TINYINT NOT NULL DEFAULT 1,
    INDEX `idx_module` (`module_id`,`submission_id`)
);

-- 用户-账号
CREATE TABLE `user`.`t_account`(
    `user_id` VARCHAR(32) NOT NULL PRIMARY KEY COMMENT '账号',
    `nickname` VARCHAR(128) NOT NULL DEFAULT '' unique COMMENT 'encode后昵称',
    `password` VARCHAR(64) NOT NULL COMMENT '哈希后的密码',
    `user_type` VARCHAR(32) NOT NULL COMMENT '账号类型'
);

-- 群组相关
CREATE TABLE `user`.`t_group`(
    `group_id` INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `owner` VARCHAR (32) NOT NULL COMMENT '建组的teacher id',
    `name` VARCHAR(32) NOT NULL DEFAULT '',
    INDEX `idx_owner` (`owner`)
);

CREATE TABLE `user`.`t_group_user`(
    `group_id` INT NOT NULL,
    `user_id` VARCHAR(32) NOT NULL,
    PRIMARY KEY(`group_id`,`user_id`),
    INDEX `idx_user`(`user_id`)
);
