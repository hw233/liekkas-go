CREATE TABLE `user` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `info` longtext,
  `item_pack` longtext,
  `name` tinytext,
  `character_pack` longtext,
  `equipment_pack` longtext,
  `world_item_pack` longtext,
  `hero_pack` longtext,
  `manual_info` longtext,
  `item_dropped_info` longtext,
  `quest_pack` longtext,
  `graveyard` longtext,
  `levels_info` longtext,
  `explore_info` longtext,
  `chapter_info` longtext,
  `tower_info` longtext,
  `guild` longtext,
  `store_info` longtext,
  `mail_info` longtext,
  `yggdrasil` longtext,
  `gacha_record` longtext,
  `activity_info` longtext,
  `score_pass_info` longtext,
  `formation_info` longtext,
  `mercenary` longtext,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=123213128023213124 DEFAULT CHARSET=utf8mb4;

CREATE TABLE `survey` (
  `user_id` bigint(20) NOT NULL,
  `survey_id` int(8) NOT NULL,
  `survey_answers` text,
  PRIMARY KEY (`user_id`,`survey_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;