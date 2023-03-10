CREATE TABLE `server_group_mails` (
  `id` bigint NOT NULL,
  `template_id` int DEFAULT NULL,
  `title` text,
  `title_args` text,
  `content` text,
  `content_args` text,
  `attachment` text,
  `sender` varchar(64) DEFAULT NULL,
  `send_time` int DEFAULT NULL,
  `expire_time` int DEFAULT NULL,
  `end_time` int DEFAULT NULL,
  `users` longtext,
  `send_all` tinyint DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;

CREATE TABLE `server_personal_mails` (
  `id` bigint NOT NULL,
  `template_id` int DEFAULT NULL,
  `title` text,
  `title_args` text,
  `content` text,
  `content_args` text,
  `attachment` text,
  `sender` varchar(64) DEFAULT NULL,
  `send_time` int DEFAULT NULL,
  `expire_time` int DEFAULT NULL,
  `end_time` int DEFAULT NULL,
  `user_id` bigint DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;
