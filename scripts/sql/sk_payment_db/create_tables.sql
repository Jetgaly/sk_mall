CREATE TABLE
    IF NOT EXISTS payment_logs (
        id BIGINT PRIMARY KEY AUTO_INCREMENT,
        order_no BIGINT NOT NULL COMMENT '订单号',
        dtm_gid VARCHAR(128) NOT NULL COMMENT 'DTM事务ID',
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        UNIQUE KEY uk_order_no (order_no),
        INDEX idx_dtm_gid (dtm_gid)
    ) COMMENT '支付日志表';