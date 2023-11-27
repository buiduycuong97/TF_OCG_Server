ALTER TABLE `e-commerce`.`review`
    RENAME TO  `e-commerce`.`reviews` ;

ALTER TABLE `e-commerce`.`reviews`
DROP COLUMN `status`,
DROP COLUMN `title`;


ALTER TABLE `e-commerce`.`reviews`
DROP COLUMN `approve`;

ALTER TABLE discount
    ADD COLUMN available_quantity INT DEFAULT 0;

ALTER TABLE orders
    ADD COLUMN total_quantity INT,
ADD COLUMN total_price FLOAT,
ADD COLUMN discount_amount FLOAT;

ALTER TABLE `e-commerce`.`discount`
    RENAME TO  `e-commerce`.`discounts` ;

ALTER TABLE `e-commerce`.`user_discount`
    RENAME TO  `e-commerce`.`user_discounts` ;
