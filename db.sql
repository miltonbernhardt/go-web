CREATE DATABASE IF NOT EXISTS `storage` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;

USE `storage`;

CREATE TABLE `users`
(
    `id`           int(11)     NOT NULL,
    `firstname`    varchar(30) NOT NULL,
    `lastname`     varchar(30) NOT NULL,
    `email`        varchar(60) NOT NULL,
    `age`          int(3)      NOT NULL,
    `height`       int(3)      NOT NULL,
    `active`       bool        NOT NULL,
    `created_date` varchar(60) NOT NULL
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

ALTER TABLE `users`
    ADD PRIMARY KEY (`id`);

ALTER TABLE `users`
    MODIFY `id` int(11) NOT NULL AUTO_INCREMENT,
    AUTO_INCREMENT = 4;

COMMIT;
