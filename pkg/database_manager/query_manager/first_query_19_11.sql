
CREATE SCHEMA IF NOT EXISTS `e-commerce`;
USE `e-commerce`;

-- Tạo bảng User
CREATE TABLE users (
                       user_id INT PRIMARY KEY AUTO_INCREMENT not null,
                       user_name VARCHAR(255),
                       email VARCHAR(255),
                       password VARCHAR(255), -- Lưu ý: Cần hash mật khẩu.
                       role varchar(50),
                       user_type VARCHAR(50),
                       refresh_token varchar(255),
                       reset_token varchar(255),
                       created_at datetime,
                       updated_at datetime
);

-- Tạo bảng Category
CREATE TABLE category (
                          category_id INT PRIMARY KEY auto_increment not null,
                          category_name VARCHAR(255)
);

-- Tạo bảng Product
CREATE TABLE product (
                         product_id INT PRIMARY KEY auto_increment not null,
                         handle VARCHAR(255), -- url vi du: le-duy-cuong
                         title VARCHAR(255),
                         description TEXT,
                         price DECIMAL(10, 2),
                         category_id INT,
                         FOREIGN KEY (category_id) REFERENCES category(category_id)
);

-- Tạo bảng Cart
CREATE TABLE cart (
                      cart_id INT PRIMARY KEY auto_increment not null,
                      user_id INT,
                      product_id INT,
                      quantity INT,
                      FOREIGN KEY (user_id) REFERENCES users(user_id),
                      FOREIGN KEY (product_id) REFERENCES product(product_id)
);

-- Tạo bảng Order
CREATE TABLE orders (
                        order_id INT PRIMARY KEY auto_increment not null,
                        user_id INT,
                        order_date DATE,
                        shipping_address text,
                        status VARCHAR(50),
                        FOREIGN KEY (user_id) REFERENCES users(user_id)
);

-- Tạo bảng OrderDetail
CREATE TABLE order_detail (
                              order_detail_id INT PRIMARY KEY auto_increment not null,
                              order_id INT,
                              product_id INT,
                              quantity INT,
                              price DECIMAL(10, 2),
                              FOREIGN KEY (order_id) REFERENCES orders(order_id),
                              FOREIGN KEY (product_id) REFERENCES product(product_id)
);

-- Tạo bảng ShippingOption
CREATE TABLE shipping_option (
                                 shipping_option_id int PRIMARY KEY auto_increment not null,
                                 option_name VARCHAR(255),
                                 cost DECIMAL(10, 2)
);

-- Tạo bảng Review
CREATE TABLE review (
                        review_id INT PRIMARY KEY auto_increment not null,
                        user_id INT,
                        product_id INT,
                        rating INT,
                        title text,
                        comment text,
                        status varchar(50),
                        FOREIGN KEY (user_id) REFERENCES users(user_id),
                        FOREIGN KEY (product_id) REFERENCES product(product_id)
);

-- Tạo bảng ReviewApproval
CREATE TABLE review_approval (
                                 review_id INT PRIMARY KEY auto_increment not null,
                                 admin_id INT,
                                 approved BOOLEAN,
                                 FOREIGN KEY (review_id) REFERENCES review(review_id)
);

-- Tạo bảng Discount
CREATE TABLE discount (
                          discount_id INT PRIMARY KEY auto_increment not null,
                          discount_code VARCHAR(50),
                          discount_type VARCHAR(50),
                          value DECIMAL(10, 2),
                          start_date DATE,
                          end_date DATE
);

-- Tạo bảng OrderDiscount //
CREATE TABLE order_discount (
                                order_discount_id INT PRIMARY KEY auto_increment not null,
                                order_id INT,
                                discount_id INT,
                                FOREIGN KEY (order_id) REFERENCES orders(order_id),
                                FOREIGN KEY (discount_id) REFERENCES discount(discount_id)
);

-- Tạo bảng UserDiscount
CREATE TABLE user_discount (
                               user_discount_id INT PRIMARY KEY auto_increment not null,
                               user_id INT,
                               discount_id INT,
                               FOREIGN KEY (user_id) REFERENCES users(user_id),
                               FOREIGN KEY (discount_id) REFERENCES discount(discount_id)
);

-- Tạo bảng TechnicalSupport
CREATE TABLE technical_support (
                                   technical_support_id INT PRIMARY KEY auto_increment not null,
                                   user_id INT,
                                   issue_description TEXT,
                                   resolution TEXT,
                                   FOREIGN KEY (user_id) REFERENCES users(user_id)
);

-- Tạo bảng UserPolicy
CREATE TABLE user_policy (
                             user_policy_id INT PRIMARY KEY auto_increment not null,
                             policy_type VARCHAR(50),
                             content TEXT
);

-- Tạo bảng UserShippingChoice
CREATE TABLE user_shipping_choice (
                                      user_shipping_choice_id INT PRIMARY KEY auto_increment not null,
                                      user_id INT,
                                      shipping_option_id INT,
                                      FOREIGN KEY (user_id) REFERENCES users(user_id),
                                      FOREIGN KEY (shipping_option_id) REFERENCES shipping_option(shipping_option_id)
);

-- Tạo bảng ProductRecommendation
CREATE TABLE product_recommendation (
                                        product_recommendation_id INT PRIMARY KEY auto_increment not null,
                                        user_id INT,
                                        product_id INT,
                                        FOREIGN KEY (user_id) REFERENCES users(user_id),
                                        FOREIGN KEY (product_id) REFERENCES product(product_id)
);

-- Tạo bảng CustomerRewards
CREATE TABLE customer_rewards (
                                  customer_rewards_id INT PRIMARY KEY auto_increment not null,
                                  user_id INT,
                                  reward_type VARCHAR(50),
                                  reward_value DECIMAL(10, 2),
                                  FOREIGN KEY (user_id) REFERENCES users(user_id)
);

-- Tạo bảng AutomaticNotification
CREATE TABLE automatic_notification (
                                        automatic_notification_id INT PRIMARY KEY auto_increment not null,
                                        user_id INT,
                                        product_id INT,
                                        notification_type VARCHAR(50),
                                        notification_content TEXT,
                                        FOREIGN KEY (user_id) REFERENCES users(user_id),
                                        FOREIGN KEY (product_id) REFERENCES product(product_id)
);

-- sql2
CREATE TABLE option_product (
                                option_product_id INT PRIMARY KEY auto_increment not null,
                                product_id INT,
                                option_type VARCHAR(50),
                                FOREIGN KEY (product_id) REFERENCES product(product_id)
);

CREATE TABLE option_value (
                              option_value_id INT PRIMARY KEY auto_increment not null,
                              option_product_id INT,
                              value VARCHAR(255),
                              FOREIGN KEY (option_product_id) REFERENCES option_product(option_product_id)
);

CREATE TABLE variant (
                         variant_id INT PRIMARY KEY auto_increment not null,
                         product_id int,
                         title VARCHAR(255),
                         price DECIMAL(10, 2),
                         compare_price DECIMAL(10, 2),
                         count_in_stock int,
                         image VARCHAR(255),
                         role VARCHAR(50),
                         option_product_1 INT,
                         option_product_2 INT,
                         FOREIGN KEY (product_id) REFERENCES product(product_id),
                         FOREIGN KEY (option_product_1) REFERENCES option_product(option_product_id),
                         FOREIGN KEY (option_product_2) REFERENCES option_product(option_product_id)
);

