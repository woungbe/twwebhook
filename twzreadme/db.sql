-- MySQL Workbench Forward Engineering

SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION';

-- -----------------------------------------------------
-- Schema mydb
-- -----------------------------------------------------

-- -----------------------------------------------------
-- Schema mydb
-- -----------------------------------------------------
CREATE SCHEMA IF NOT EXISTS `mydb` DEFAULT CHARACTER SET utf8 ;
USE `mydb` ;

-- -----------------------------------------------------
-- Table `mydb`.`users`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `mydb`.`users` (
  `useridx` INT NOT NULL AUTO_INCREMENT,
  `accesskey` VARCHAR(100) NOT NULL,
  `scriptkey` VARCHAR(100) NOT NULL,
  `createdate` DATETIME NOT NULL DEFAULT now(),
  `userid` VARCHAR(45) NOT NULL,
  PRIMARY KEY (`useridx`))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `mydb`.`strategy`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `mydb`.`strategy` (
  `stsrl` INT NOT NULL,
  `stname` VARCHAR(45) NULL,
  `stcontent` TEXT NULL,
  PRIMARY KEY (`stsrl`))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `mydb`.`mapping`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `mydb`.`mapping` (
  `useridx` INT NOT NULL AUTO_INCREMENT,
  `users_useridx` INT NOT NULL,
  `strategy_stsrl` INT NOT NULL,
  `profitflg` TINYINT NULL,
  `profitval` VARCHAR(45) NULL,
  `losscutflg` TINYINT NULL,
  `losscutval` VARCHAR(45) NULL,
  `rateflg` TINYINT NULL COMMENT '자산 대비 비율 : 0 \n고정금 차임 : 1 ',
  `rateval` INT NULL COMMENT '0% ~ 100%\n',
  `fixedval` INT NULL COMMENT 'USDT 기준 고정 ',
  `counterflg` INT NULL COMMENT '0: 알아서 익절 나것지… 뭔가 안함 \n1: 롱 매수일때 숏청산,  숏매수 일때 롱 청산',
  INDEX `fk_mapping_users_idx` (`users_useridx` ASC) VISIBLE,
  INDEX `fk_mapping_strategy1_idx` (`strategy_stsrl` ASC) VISIBLE,
  INDEX `PRIMARIKEY` (`users_useridx` ASC, `strategy_stsrl` ASC) VISIBLE,
  PRIMARY KEY (`useridx`),
  CONSTRAINT `fk_mapping_users`
    FOREIGN KEY (`users_useridx`)
    REFERENCES `mydb`.`users` (`useridx`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_mapping_strategy1`
    FOREIGN KEY (`strategy_stsrl`)
    REFERENCES `mydb`.`strategy` (`stsrl`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `mydb`.`ordersetting`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `mydb`.`ordersetting` (
  `settingsrl` INT NOT NULL AUTO_INCREMENT,
  `useridx` INT NOT NULL,
  `setname` VARCHAR(45) NOT NULL,
  `profitflg` TINYINT NOT NULL DEFAULT 0,
  `profitval` VARCHAR(45) NOT NULL DEFAULT '',
  `losscutflg` TINYINT NOT NULL DEFAULT 0,
  `losscutval` VARCHAR(45) NOT NULL DEFAULT '',
  `rateflg` TINYINT NOT NULL DEFAULT 0,
  `rateval` INT NULL DEFAULT 30,
  `fixedval` INT NULL,
  `counterflg` INT NOT NULL DEFAULT 0,
  INDEX `fk_ordersetting_users1_idx` (`useridx` ASC) VISIBLE,
  PRIMARY KEY (`settingsrl`),
  CONSTRAINT `fk_ordersetting_users1`
    FOREIGN KEY (`useridx`)
    REFERENCES `mydb`.`users` (`useridx`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `mydb`.`signnalLog`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `mydb`.`signnalLog` (
  `idx` INT NOT NULL AUTO_INCREMENT,
  `stsrl` INT NOT NULL,
  `jsonData` TEXT NULL,
  `createdate` DATETIME NULL DEFAULT now(),
  `actionCnt` INT NULL,
  `mapping` VARCHAR(255) NULL,
  PRIMARY KEY (`idx`))
ENGINE = InnoDB;


SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;
