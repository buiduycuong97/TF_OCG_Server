ALTER TABLE `e-commerce`.`products`
    CHANGE COLUMN `price` `price` INT NULL DEFAULT NULL ;
ALTER TABLE `e-commerce`.`products`
    ADD COLUMN `image` TEXT NULL AFTER `updated_at`;
