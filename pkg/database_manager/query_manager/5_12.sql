ALTER TABLE `e-commerce`.`variants`
    CHANGE COLUMN `price` `price` INT NULL DEFAULT NULL ,
    CHANGE COLUMN `compare_price` `compare_price` INT NULL DEFAULT NULL ;
