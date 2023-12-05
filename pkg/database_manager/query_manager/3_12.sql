ALTER TABLE `e-commerce`.`reviews`
DROP FOREIGN KEY `reviews_ibfk_2`;
ALTER TABLE `e-commerce`.`reviews`
DROP COLUMN `product_id`,
ADD COLUMN `variant_id` INT NULL DEFAULT NULL AFTER `comment`,
DROP INDEX `product_id` ;
;

ALTER TABLE reviews
    ADD CONSTRAINT fk_review_variant
        FOREIGN KEY (variant_id)
            REFERENCES variants(variant_id)
            ON UPDATE CASCADE
            ON DELETE CASCADE;

ALTER TABLE `e-commerce`.`reviews`
    ADD COLUMN `create_at` DATETIME NULL DEFAULT NULL AFTER `variant_id`;


ALTER TABLE `e-commerce`.`reviews`
    CHANGE COLUMN `create_at` `created_at` DATETIME NULL DEFAULT NULL ;

ALTER TABLE order_details
    ADD COLUMN is_review BOOLEAN;

ALTER TABLE `e-commerce`.`orders`
    ADD COLUMN `created_at` DATETIME NULL AFTER `grand_total`;
