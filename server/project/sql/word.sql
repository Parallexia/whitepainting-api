CREATE TABLE IF NOT EXISTS `words`(
    `data_id` INT UNSIGNED AUTO_INCREMENT,
    `content` VARCHAR(192) NOT NULL,
    PRIMARY KEY (`data_id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO words
(content)
VALUES
("");
