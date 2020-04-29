create database user;

drop table if exists `user`;
CREATE TABLE `user`
(
  `id`        INT(10)     NOT NULL AUTO_INCREMENT,
  `name`   VARCHAR(64) NULL DEFAULT NULL,
  `habits` VARCHAR(128) NULL DEFAULT NULL,
  `created_time`    DATE        NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
);