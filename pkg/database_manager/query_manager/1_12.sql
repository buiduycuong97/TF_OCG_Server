ALTER TABLE `e-commerce`.`variants`
DROP FOREIGN KEY `variants_ibfk_3`,
DROP FOREIGN KEY `variants_ibfk_2`;
ALTER TABLE `e-commerce`.`variants`
DROP COLUMN `option_product_2`,
DROP COLUMN `option_product_1`,
DROP INDEX `variants_ibfk_3_idx` ,
DROP INDEX `variants_ibfk_2_idx` ;
;

ALTER TABLE `e-commerce`.`variants`
    ADD COLUMN `option_value_1` INT NULL AFTER `image`,
ADD COLUMN `option_value_2` INT NULL AFTER `option_value_1`,
ADD INDEX `variants_ibfk_2_idx` (`option_value_1` ASC) VISIBLE,
ADD INDEX `variants_ibfk_3_idx` (`option_value_2` ASC) VISIBLE;
;
ALTER TABLE `e-commerce`.`variants`
    ADD CONSTRAINT `variants_ibfk_2`
        FOREIGN KEY (`option_value_1`)
            REFERENCES `e-commerce`.`option_values` (`option_value_id`)
            ON DELETE NO ACTION
            ON UPDATE NO ACTION,
ADD CONSTRAINT `variants_ibfk_3`
  FOREIGN KEY (`option_value_2`)
  REFERENCES `e-commerce`.`option_values` (`option_value_id`)
  ON DELETE NO ACTION
  ON UPDATE NO ACTION;

--CuongLD
ALTER TABLE `e-commerce`.`variants`
DROP FOREIGN KEY `variants_ibfk_2`,
DROP FOREIGN KEY `variants_ibfk_3`;
ALTER TABLE `e-commerce`.`variants`
    CHANGE COLUMN `option_value_1` `option_value1` INT NULL DEFAULT NULL ,
    CHANGE COLUMN `option_value_2` `option_value2` INT NULL DEFAULT NULL ;
ALTER TABLE `e-commerce`.`variants`
    ADD CONSTRAINT `variants_ibfk_2`
        FOREIGN KEY (`option_value1`)
            REFERENCES `e-commerce`.`option_values` (`option_value_id`),
ADD CONSTRAINT `variants_ibfk_3`
  FOREIGN KEY (`option_value2`)
  REFERENCES `e-commerce`.`option_values` (`option_value_id`);

-- Thêm cột VariantID và thiết lập nó làm khóa ngoại trong bảng Cart
ALTER TABLE carts
    ADD COLUMN variant_id INT,
ADD FOREIGN KEY (variant_id) REFERENCES variants(variant_id);

-- Thêm cột VariantID và thiết lập nó làm khóa ngoại trong bảng OrderDetail
ALTER TABLE order_details
    ADD COLUMN variant_id INT,
ADD FOREIGN KEY (variant_id) REFERENCES variants(variant_id);

ALTER TABLE `e-commerce`.`carts`
DROP FOREIGN KEY `carts_ibfk_2`;
ALTER TABLE `e-commerce`.`carts`
DROP COLUMN `product_id`,
DROP INDEX `product_id` ;
;

ALTER TABLE `e-commerce`.`products`
DROP COLUMN `quantity_remaining`;

ALTER TABLE `e-commerce`.`variants`
    CHANGE COLUMN `price` `price` INT NULL DEFAULT NULL ;

ALTER TABLE `e-commerce`.`variants`
    CHANGE COLUMN `compare_price` `compare_price` INT NULL DEFAULT NULL ;

ALTER TABLE `e-commerce`.`order_details`
DROP FOREIGN KEY `order_details_ibfk_2`;
ALTER TABLE `e-commerce`.`order_details`
DROP COLUMN `product_id`,
DROP INDEX `product_id` ;
;
