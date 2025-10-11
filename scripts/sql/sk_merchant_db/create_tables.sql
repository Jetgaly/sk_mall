CREATE TABLE
    IF NOT EXISTS merchants (
        id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
        user_id BIGINT NOT NULL UNIQUE COMMENT '关联用户ID（店主账号，来自用户服务）',
        name VARCHAR(100) NOT NULL COMMENT '商户/店铺名称',
        logo VARCHAR(255) COMMENT '店铺Logo',
        description VARCHAR(500) COMMENT '店铺简介',
        type TINYINT NOT NULL DEFAULT 1 COMMENT '商户类型：1-个人，2-企业',
        status TINYINT NOT NULL DEFAULT 1 COMMENT '状态：0-待审核，1-正常，2-冻结，3-关闭',
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    ) COMMENT = '商户表';

CREATE TABLE
    IF NOT EXISTS merchant_auths (
        merchant_id BIGINT PRIMARY KEY COMMENT '商户ID 主键，关联merchants的外键',
        legal_name VARCHAR(50) COMMENT '联系人姓名',
        id_card VARCHAR(30) COMMENT '身份证号',
        license_img VARCHAR(255) COMMENT '营业执照图片',
        status TINYINT NOT NULL DEFAULT 1 COMMENT '认证状态：0-待审核，1-已通过，2-未通过',
        remark VARCHAR(255) COMMENT '审核备注',
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        FOREIGN KEY (merchant_id) REFERENCES merchants (id)
    ) COMMENT = '商户认证表';

CREATE TABLE
    IF NOT EXISTS merchant_accounts (
        merchant_id BIGINT PRIMARY KEY COMMENT '商户ID',
        balance DECIMAL(18, 2) DEFAULT 0.00,
        account_type TINYINT COMMENT '账户类型：1-银行卡，2-支付宝，3-微信',
        account_no VARCHAR(100) COMMENT '收款账号',
        account_name VARCHAR(50) COMMENT '账户姓名/开户名',
        bank_name VARCHAR(100) COMMENT '银行名称（如果是银行卡）',
        status TINYINT DEFAULT 1 COMMENT '状态：0-无效，1-有效',
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        FOREIGN KEY (merchant_id) REFERENCES merchants (id)
    ) COMMENT = '商户结算账户表';

CREATE TABLE
    IF NOT EXISTS pay_logs (
        id BIGINT PRIMARY KEY COMMENT 'snowflake order no',
        price DECIMAL(18, 2) CHECK (price > 0) NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间'
    )