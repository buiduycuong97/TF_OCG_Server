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
ADD INDEX `variants_ibfk_2_idx` (`option_value_1` ASC) VISIBLE;
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
