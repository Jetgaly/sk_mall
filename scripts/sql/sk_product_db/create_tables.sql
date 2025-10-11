CREATE TABLE
    IF NOT EXISTS products (
        id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '商品ID',
        merchant_id BIGINT NOT NULL COMMENT '商户ID',
        name VARCHAR(100) NOT NULL COMMENT '商品名称',
        description TEXT COMMENT '商品描述',
        cover_image VARCHAR(255) COMMENT '商品主图',
        price DECIMAL(10, 2) NOT NULL COMMENT '原价',
        stock INT NOT NULL DEFAULT 0 COMMENT '总库存（非秒杀库存）',
        status TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1=上架，0=下架',
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
        updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
        UNIQUE KEY uk_merchant_name (merchant_id, name)
    ) COMMENT = '商品表';

CREATE TABLE
    IF NOT EXISTS seckill_events (
        id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '秒杀活动ID',
        merchant_id BIGINT NOT NULL COMMENT '商户ID',
        name VARCHAR(100) NOT NULL COMMENT '活动名称',
        start_time DATETIME NOT NULL COMMENT '活动开始时间',
        end_time DATETIME NOT NULL COMMENT '活动结束时间',
        status TINYINT NOT NULL DEFAULT 1 COMMENT '活动状态：1=未开始，2=进行中，3=已结束',
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
        updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
        UNIQUE KEY uk_merchant_name (merchant_id, name)
    ) COMMENT = '秒杀活动表';

CREATE TABLE
    IF NOT EXISTS seckill_products (
        id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '秒杀商品ID',
        event_id BIGINT NOT NULL COMMENT '所属秒杀活动ID',
        product_id BIGINT NOT NULL COMMENT '关联商品ID',
        seckill_price DECIMAL(10, 2) NOT NULL COMMENT '秒杀价格',
        seckill_stock INT NOT NULL COMMENT '秒杀库存',
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
        updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
        UNIQUE KEY uk_event_product (event_id, product_id),
        FOREIGN KEY (event_id) REFERENCES seckill_events (id),
        FOREIGN KEY (product_id) REFERENCES products (id)
    ) COMMENT = '秒杀商品表';