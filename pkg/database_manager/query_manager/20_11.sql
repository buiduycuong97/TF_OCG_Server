
ALTER TABLE category
    CHANGE category_name name VARCHAR(255);

ALTER TABLE category
    ADD COLUMN value VARCHAR(255)

ALTER TABLE `e-commerce`.`category`
    ADD COLUMN `image` TEXT NULL AFTER `handle`,
CHANGE COLUMN `value` `handle` VARCHAR(255) NULL DEFAULT NULL ;

ALTER TABLE `e-commerce`.`variant`
DROP COLUMN `role`;
