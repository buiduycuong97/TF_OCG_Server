use `e-commerce`

CREATE TABLE province (
                          province_id INT PRIMARY KEY AUTO_INCREMENT NOT NULL,
                          province_name VARCHAR(255) NOT NULL,
                          shipping_fee DECIMAL(10,2) DEFAULT 0
);

ALTER TABLE orders
ADD COLUMN province_id INT;

-- Thêm ràng buộc khóa ngoại
ALTER TABLE orders
    ADD CONSTRAINT fk_province
        FOREIGN KEY (province_id)
            REFERENCES province(province_id);

ALTER TABLE users
    ADD COLUMN phonenumber VARCHAR(20) NOT NULL;

ALTER TABLE `e-commerce`.`province`
    RENAME TO  `e-commerce`.`provinces` ;
