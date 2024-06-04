//创建表
CREATE TABLE IF NOT EXISTS `t_user` (
                                    `id` INT NOT NULL AUTO_INCREMENT,
                                    `name` VARCHAR(100) NOT NULL,
                                    `age` INT NOT NULL,
                                    `gender` VARCHAR(30) NOT NULL,
                                    `password` VARCHAR(255) NOT NULL DEFAULT '',
                                    `nickname` VARCHAR(100) NOT NULL DEFAULT '',
                                    `create_time` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                    `creator` VARCHAR(100) NOT NULL DEFAULT '',
                                    `modify_time` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后一次修改时间',
                                    `modifier` VARCHAR(100) NOT NULL DEFAULT '',
                                  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci