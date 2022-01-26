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


CREATE TABLE `warehouses`
(
    `id`     int(11)      NOT NULL,
    `name`   varchar(40)  NOT NULL,
    `adress` varchar(150) NOT NULL
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

INSERT INTO `warehouses` (`id`, `name`, `adress`)
VALUES (1, 'Main Warehouse', '221b Baker Street');
-- Volcado de datos para la tabla `warehouses`
ALTER TABLE `warehouses`
    ADD PRIMARY KEY (`id`);
ALTER TABLE `warehouses`
    MODIFY `id` int(11) NOT NULL AUTO_INCREMENT,
    AUTO_INCREMENT = 2;


ALTER TABLE `users`
    ADD `id_warehouse` INT NOT NULL AFTER `created_date`;

UPDATE `users`
SET `id_warehouse` = '1';

COMMIT;
