-- iOS 内购赞助：用户表增加字段 + 订单流水表
-- 执行前请确认 app_users 已存在；若字段/表已存在可跳过对应语句

ALTER TABLE app_users
  ADD COLUMN vip INT DEFAULT 0,
  ADD COLUMN support_total_amount DECIMAL(12,2) DEFAULT 0,
  ADD COLUMN support_level INT DEFAULT 0;

CREATE TABLE IF NOT EXISTS app_support_orders (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  user_id BIGINT NOT NULL,
  platform VARCHAR(20) NOT NULL DEFAULT 'ios',
  product_id VARCHAR(128) NOT NULL,
  transaction_id VARCHAR(128) NOT NULL,
  amount DECIMAL(12,2) NOT NULL DEFAULT 0,
  receipt_data TEXT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY uk_transaction_id (transaction_id),
  KEY idx_user_created (user_id, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
