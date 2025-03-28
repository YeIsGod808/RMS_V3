-- Active: 1740461445112@@192.168.80.128@3309@resource
-- +goose Up

-- 创建视频表
-- 创建视频表
CREATE TABLE IF NOT EXISTS `videos` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '视频ID',
    `title` VARCHAR(100) NOT NULL COMMENT '视频标题',
    `play_url` VARCHAR(255) NOT NULL COMMENT '视频播放URL',
    `cover_url` VARCHAR(255) NOT NULL COMMENT '视频封面URL',
    `description` TEXT COMMENT '视频描述',
    `knowledge_point_id` BIGINT NOT NULL COMMENT '关联的知识点ID(Neo4j)',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    INDEX `idx_knowledge_point_id` (`knowledge_point_id`) COMMENT '知识点ID索引'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='视频表';

-- 创建练习题表
CREATE TABLE IF NOT EXISTS `exercises` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '练习题ID',
    `title` VARCHAR(100) NOT NULL COMMENT '练习题标题',
    `exercise_url` VARCHAR(255) NOT NULL COMMENT '练习题URL',
    `difficulty` ENUM('easy', 'medium', 'hard') NOT NULL COMMENT '练习题难度',
    `description` TEXT COMMENT '练习题描述',
    `knowledge_point_id` BIGINT NOT NULL COMMENT '关联的知识点ID(Neo4j)',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    INDEX `idx_knowledge_point_id` (`knowledge_point_id`) COMMENT '知识点ID索引'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='练习题表';

-- 创建课件表
CREATE TABLE IF NOT EXISTS `coursewares` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '课件ID',
    `title` VARCHAR(100) NOT NULL COMMENT '课件标题',
    `courseware_url` VARCHAR(255) NOT NULL COMMENT '课件URL',
    `description` TEXT COMMENT '课件描述',
    `knowledge_point_id` BIGINT NOT NULL COMMENT '关联的知识点ID(Neo4j)',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    INDEX `idx_knowledge_point_id` (`knowledge_point_id`) COMMENT '知识点ID索引'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='课件表';


-- +goose Down


-- 删除课件表
DROP TABLE IF EXISTS `coursewares`;

-- 删除练习题表
DROP TABLE IF EXISTS `exercises`;

-- 删除视频表
DROP TABLE IF EXISTS `videos`;

