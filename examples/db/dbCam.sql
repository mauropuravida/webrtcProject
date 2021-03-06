-- MySQL Script generated by MySQL Workbench
-- Sun Dec 20 15:56:20 2020
-- Model: New Model    Version: 1.0
-- MySQL Workbench Forward Engineering

SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION';

-- -----------------------------------------------------
-- Schema dbcam
-- -----------------------------------------------------

-- -----------------------------------------------------
-- Schema dbcam
-- -----------------------------------------------------
CREATE SCHEMA IF NOT EXISTS `dbcam` DEFAULT CHARACTER SET utf8 ;
USE `dbcam` ;

-- -----------------------------------------------------
-- Table `dbcam`.`users`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `dbcam`.`Users` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(50) NOT NULL,
  `surname` VARCHAR(50) NOT NULL,
  `age` INT NOT NULL,
  `email` VARCHAR(50) NOT NULL,
  `created` DATE NOT NULL,
  `password` VARCHAR(50) NOT NULL,
  PRIMARY KEY (`id`))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `dbcam`.`cameras`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `dbcam`.`Cameras` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `active` TINYINT(1) NOT NULL,
  `created` DATE NOT NULL,
  `loc` VARCHAR(50) NOT NULL,
  `url` VARCHAR(100) NULL,
  `token_session_camera` VARCHAR(10000) NULL,
  `token_session_consumer` VARCHAR(10000) NULL,
  `id_camera` INT NOT NULL,
  `users_id` INT NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_cameras_users_idx` (`users_id` ASC) VISIBLE,
  CONSTRAINT `fk_cameras_users`
    FOREIGN KEY (`users_id`)
    REFERENCES `dbcam`.`Users` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;
