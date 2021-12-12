CREATE TABLE IF NOT EXISTS `concentrate`(
    `data_id` INT UNSIGNED AUTO_INCREMENT,
    `session_id` VARCHAR(32) NOT NULL,
    `user_id` INT,
    `timestart` DATETIME NOT NULL,
    `timesub` DATETIME NOT NULL,
    PRIMARY KEY (`data_id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `community`(
    `data_id` INT UNSIGNED AUTO_INCREMENT,
    `user_id` INT NOT NULL,
    `content` VARCHAR(384) NOT NULL,
    `pic_url` VARCHAR(384) ,
    `submit_time`  DATETIME NOT NULL,
    PRIMARY KEY (`data_id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `concentrate_plan`(
    `data_id` INT UNSIGNED AUTO_INCREMENT,
    `user_id` INT,
    `subject` VARCHAR(32),
    `description` VARCHAR(192)
    `timestart` DATETIME NOT NULL,
    `timesub` DATETIME NOT NULL,
    PRIMARY KEY (`data_id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `message`(
    `data_id` INT UNSIGNED AUTO_INCREMENT,
    `like` INT UNSIGNED,
)