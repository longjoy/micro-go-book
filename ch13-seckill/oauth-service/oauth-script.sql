create database oauth;

drop table if exists `client_details`;
CREATE TABLE `client_details`
(
    `client_id`     VARCHAR(128)      NOT NULL,
    `client_secret`    VARCHAR(128)  NOT NULL,
    `access_token_validity_seconds`   INT(10) NOT NULL,
    `refresh_token_validity_seconds`  INT(10) NOT NULL,
    `registered_redirect_uri` VARCHAR(128)      NOT NULL,
    `authorized_grant_types` VARCHAR(128)      NOT NULL,

    PRIMARY KEY (`client_id`)
);

insert into client_details (`client_id`, `client_secret`, `access_token_validity_seconds`, `refresh_token_validity_seconds`, `registered_redirect_uri`, `authorized_grant_types`) values ('clientId', 'clientSecret', 1800, 18000, 'http://127.0.0.1','["password","refresh_token"]');