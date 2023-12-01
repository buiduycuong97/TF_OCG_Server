package query_manager
ALTER TABLE `e-commerce`.`orders`
    ADD COLUMN `grand_total` FLOAT NULL AFTER `discount_amount`;

ALTER TABLE `e-commerce`.`users`
    CHANGE COLUMN `phonenumber` `phone_number` VARCHAR(20) NOT NULL ;
