ALTER TABLE `e-commerce`.`review`
    RENAME TO  `e-commerce`.`reviews` ;

ALTER TABLE `e-commerce`.`reviews`
DROP COLUMN `status`,
DROP COLUMN `title`;


ALTER TABLE `e-commerce`.`reviews`
DROP COLUMN `approve`;
