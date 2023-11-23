
ALTER TABLE category
    CHANGE category_name name VARCHAR(255);

ALTER TABLE category
    ADD COLUMN value VARCHAR(255)

ALTER TABLE `e-commerce`.`category`
    ADD COLUMN `image` TEXT NULL AFTER `handle`,
CHANGE COLUMN `value` `handle` VARCHAR(255) NULL DEFAULT NULL ;

ALTER TABLE `e-commerce`.`variant`
DROP COLUMN `role`;


ALTER TABLE cart
    ADD COLUMN product_detail JSON;

ALTER TABLE cart
    ADD COLUMN total_price DECIMAL(10, 2);

ALTER TABLE product
    ADD COLUMN quantity_remaining INT;

ALTER TABLE `e-commerce`.`category`
    RENAME TO  `e-commerce`.`categories` ;

ALTER TABLE `e-commerce`.`product`
    RENAME TO  `e-commerce`.`products` ;


ALTER TABLE `e-commerce`.`products`
    ADD COLUMN `created_at` DATETIME NULL AFTER `quantity_remaining`,
ADD COLUMN `updated_at` DATETIME NULL AFTER `create_at`;


ALTER TABLE `e-commerce`.`cart`
    RENAME TO  `e-commerce`.`carts` ;

ALTER TABLE `e-commerce`.`users`
    ADD COLUMN `reset` VARCHAR(45) NULL AFTER `reset_token`;

ALTER TABLE `e-commerce`.`order_detail`
    RENAME TO  `e-commerce`.`order_details` ;
