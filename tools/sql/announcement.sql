CREATE TABLE `announcement` (
	`id` bigint NOT NULL,
	`type` int DEFAULT NULL,
	`module` int DEFAULT NULL,
	`title` text DEFAULT NULL,
	`content` text DEFAULT NULL,
	`tag` text DEFAULT NULL,
	`image` text DEFAULT NULL,
	`start_time` int DEFAULT NULL,
	`end_time` int DEFAULT NULL,
	`show_start_time` int DEFAULT NULL,
	`show_end_time` int DEFAULT NULL,
	`priority` int DEFAULT NULL,
	PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;

CREATE TABLE `banner` (
	`id` bigint NOT NULL,
	`type` int DEFAULT NULL,
	`module` int DEFAULT NULL,
	`jump` text DEFAULT NULL,
	`image` text DEFAULT NULL,
	`start_time` int DEFAULT NULL,
	`end_time` int DEFAULT NULL,
	`show_start_time` int DEFAULT NULL,
	`show_end_time` int DEFAULT NULL,
	`priority` int DEFAULT NULL,
	PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;

CREATE TABLE `caution` (
	`id` bigint NOT NULL,
	`content` text DEFAULT NULL,
	`start_time` int DEFAULT NULL,
	`end_time` int DEFAULT NULL,
	PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;