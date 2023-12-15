ALTER TABLE option_products
DROP FOREIGN KEY option_products_ibfk_1;

ALTER TABLE option_products
    ADD CONSTRAINT option_products_ibfk_1
        FOREIGN KEY (product_id)
            REFERENCES products(product_id)
            ON DELETE CASCADE;

ALTER TABLE option_values
DROP FOREIGN KEY option_values_ibfk_1;

ALTER TABLE option_values
    ADD CONSTRAINT option_values_ibfk_1
        FOREIGN KEY (option_product_id)
            REFERENCES option_products(option_product_id)
            ON DELETE CASCADE;

ALTER TABLE variants
DROP FOREIGN KEY variants_ibfk_1,
DROP FOREIGN KEY variants_ibfk_2,
DROP FOREIGN KEY variants_ibfk_3;

ALTER TABLE variants
    ADD CONSTRAINT variants_ibfk_1
        FOREIGN KEY (product_id)
            REFERENCES products(product_id)
            ON DELETE CASCADE;