-- +migrate Up
-- 标记管理
CREATE TABLE IF NOT EXISTS flags(
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  app_id VARCHAR(100) DEFAULT '' COMMENT 'APPID',
  name VARCHAR(100) DEFAULT '' COMMENT '标记名称',
  flag VARCHAR(100) DEFAULT '' COMMENT 'flag',
  type VARCHAR(20) COMMENT '类型',
  status int COMMENT '状态 0.未启用 1.启用',
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间戳',
  KEY flag (flag),
  KEY type (type)

)CHARACTER SET utf8mb4;