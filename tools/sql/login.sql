CREATE TABLE `login`
(
    `id`          bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `open_id`     varchar(255) DEFAULT NULL,
    `third_party` varchar(255) DEFAULT NULL,
    `password`    tinytext,
    `account`     tinytext,
    `token`       tinytext,
    `mtime`       timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `ctime`       bigint(20) DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `id` (`id`),
    UNIQUE KEY `idx_open_id_thiry_party` (`open_id`,`third_party`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4;