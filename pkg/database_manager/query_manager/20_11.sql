
ALTER TABLE category
    CHANGE category_name name VARCHAR(255);

ALTER TABLE category
    ADD COLUMN value VARCHAR(255)