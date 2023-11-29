package query_manager
ALTER TABLE `e-commerce`.`orders`
    ADD COLUMN `grand_total` FLOAT NULL AFTER `discount_amount`;