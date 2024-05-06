create schema `imageshare`;
CREATE TABLE `imageshare`.`users` (
  `username` VARCHAR(100) NOT NULL,
  `password` VARCHAR(256) NOT NULL,
  PRIMARY KEY (`username`, `password`));
CREATE TABLE `imageshare`.`images` (
  `id` VARCHAR(100) NOT NULL,
  `username` VARCHAR(100) NOT NULL,
  `imagepath` VARCHAR(512) NOT NULL,
  `imagetype` VARCHAR(5) NOT NULL,
  UNIQUE INDEX `id_UNIQUE` (`id` ASC) VISIBLE,
  PRIMARY KEY (`id`),
  CONSTRAINT `username`
    FOREIGN KEY (`username`)
    REFERENCES `imageshare`.`users` (`username`)
    ON DELETE CASCADE
    ON UPDATE CASCADE);
INSERT INTO `imageshare`.`users` (`username`, `password`) values ("default", "none");