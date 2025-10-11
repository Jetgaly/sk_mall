-- 订单主表
CREATE TABLE
  `sk_orders` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '订单ID',
    `order_no` VARCHAR(32) NOT NULL COMMENT '订单编号',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `addr_id` BIGINT UNSIGNED NOT NULL COMMENT '地址ID',
    `sk_product_id` BIGINT UNSIGNED NOT NULL COMMENT 'sk商品ID',
    `quantity` INT UNSIGNED NOT NULL DEFAULT 1 COMMENT '购买数量',
    `unit_price` DECIMAL(10, 2) UNSIGNED NOT NULL COMMENT '商品单价',
    `total_amount` DECIMAL(10, 2) UNSIGNED NOT NULL COMMENT '订单总金额',
    `status` TINYINT NOT NULL DEFAULT 0 COMMENT '订单状态：0-待支付，1-已支付，2-已发货，3-已完成，4-已取消，5-超时关闭',
    `pay_status` TINYINT NOT NULL DEFAULT 0 COMMENT '支付状态：0-未支付，1-已支付，2-支付失败',
    `pay_time` DATETIME DEFAULT NULL COMMENT '支付时间',
    `pay_type` TINYINT DEFAULT 0 COMMENT '支付方式：0-平台支付，1-微信，2-支付宝',
    `pay_transaction_id` VARCHAR(64) DEFAULT 'default' COMMENT '支付平台交易号',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `expire_time` DATETIME NOT NULL COMMENT '订单过期时间(用于未支付自动取消)',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_order_no` (`order_no`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_created_at` (`created_at`),
    KEY `idx_status` (`status`),
    KEY `idx_expire_time` (`expire_time`)
  ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = '订单主表';