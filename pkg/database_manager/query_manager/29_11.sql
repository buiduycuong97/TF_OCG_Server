package query_manager
ALTER TABLE `e-commerce`.`orders`
    ADD COLUMN `grand_total` FLOAT NULL AFTER `discount_amount`;

ALTER TABLE `e-commerce`.`users`
    CHANGE COLUMN `phonenumber` `phone_number` VARCHAR(20) NOT NULL ;

ALTER TABLE `e-commerce`.`variants`
DROP FOREIGN KEY `variants_ibfk_2`,
DROP FOREIGN KEY `variants_ibfk_3`;
ALTER TABLE `e-commerce`.`variants`
    CHANGE COLUMN `option_product_1` `option_product1` INT NULL DEFAULT NULL ,
    CHANGE COLUMN `option_product_2` `option_product2` INT NULL DEFAULT NULL ;
ALTER TABLE `e-commerce`.`variants`
    ADD CONSTRAINT `variants_ibfk_2`
        FOREIGN KEY (`option_product1`)
            REFERENCES `e-commerce`.`option_products` (`option_product_id`),
ADD CONSTRAINT `variants_ibfk_3`
  FOREIGN KEY (`option_product2`)
  REFERENCES `e-commerce`.`option_products` (`option_product_id`);
