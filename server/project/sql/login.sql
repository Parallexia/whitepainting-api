CREATE TABLE IF NOT EXISTS `users`(
    `data_id` INT UNSIGNED AUTO_INCREMENT,
    `username` VARCHAR(20) NOT NULL,
    `email` VARCHAR(64) NOT NULL,
    `salt` BINARY(4) NOT NULL,
    `password` BINARY(32) NOT NULL,
    PRIMARY KEY (`data_id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
