CREATE TABLE
  IF NOT EXISTS users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    nick_name VARCHAR(15) NOT NULL COMMENT '用户昵称',
    user_name VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名（唯一）',
    password_hash VARCHAR(255) NOT NULL COMMENT '密码哈希（加密存储）',
    phone VARCHAR(20) UNIQUE COMMENT '手机号（唯一，可选）',
    email VARCHAR(100) UNIQUE COMMENT '邮箱（唯一）',
    status TINYINT DEFAULT 1 COMMENT '用户状态：1=正常，0=禁用',
    addrs_count INT NOT NULL DEFAULT 0 COMMENT '地址数量',
    avatar VARCHAR(200) DEFAULT NULL COMMENT '头像路径',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
  ) COMMENT = '用户基本信息表';

CREATE TABLE
  IF NOT EXISTS user_wallets (
    user_id BIGINT PRIMARY KEY COMMENT '用户ID，外键关联 users(id)',
    balance DECIMAL(18, 2) DEFAULT 0.00 COMMENT '可用余额',
    frozen_balance DECIMAL(18, 2) DEFAULT 0.00 COMMENT '冻结金额（支付中或锁定资金）',
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    FOREIGN KEY (user_id) REFERENCES users (id)
  ) COMMENT = '用户钱包表';

CREATE TABLE
  IF NOT EXISTS user_addresses (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    user_id BIGINT NOT NULL COMMENT '关联的用户ID，外键关联 users(id)',
    receiver_name VARCHAR(50) COMMENT '收货人姓名',
    receiver_phone VARCHAR(20) COMMENT '收货人手机号',
    province VARCHAR(50) COMMENT '省份',
    city VARCHAR(50) COMMENT '城市',
    district VARCHAR(50) COMMENT '区/县',
    detail_address VARCHAR(255) COMMENT '详细地址（街道/门牌号）',
    is_default TINYINT DEFAULT 0 COMMENT '是否为默认地址：0=否，1=是',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    FOREIGN KEY (user_id) REFERENCES users (id)
  ) COMMENT = '用户收货地址表';

CREATE TABLE
  IF NOT EXISTS frozen_log (
    id BIGINT PRIMARY KEY COMMENT 'snowflake order no',
    price DECIMAL(18, 2) CHECK (price > 0) NOT NULL,
    status TINYINT NOT NULL DEFAULT 0 COMMENT '0==已冻结，1==已扣减',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间'
  );