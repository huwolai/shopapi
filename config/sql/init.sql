-- +migrate Up
-- 地址信息
CREATE TABLE IF NOT EXISTS address(
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  app_id VARCHAR(255) COMMENT 'APPID',
  open_id VARCHAR(255) COMMENT '用户ID',
  longitude NUMERIC(14,10) COMMENT '经度',
  latitude NUMERIC(14,10) COMMENT '维度',
  name VARCHAR(255) COMMENT '姓名',
  mobile VARCHAR(255) COMMENT '手机号',
  address VARCHAR(255) COMMENT '地址',
  weight int COMMENT '地址权重 1.权重越高 越优先',
  flag VARCHAR(100) COMMENT '标记',
  json VARCHAR(1000) COMMENT '附加字段',
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间戳',
  KEY open_id (open_id)

)CHARACTER SET utf8mb4;

-- 账户信息
CREATE TABLE IF NOT EXISTS account(
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  app_id VARCHAR(255) COMMENT 'APPID',
  open_id VARCHAR(255) COMMENT '用户ID',
  money NUMERIC(12,2) COMMENT '账户金额',
  password VARCHAR(255) COMMENT '账户密码',
  status int COMMENT '状态 1.正常 0.锁定 2.等待开通支付',
  flag VARCHAR(100) COMMENT '标记',
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间戳',
  KEY open_id (open_id)
);


-- 商户
CREATE TABLE IF NOT EXISTS merchant(

  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(255) COMMENT '商户名称',
  app_id VARCHAR(255) COMMENT 'APPID',
  open_id VARCHAR(255) COMMENT '商户open_id',
  longitude NUMERIC(14,10) COMMENT '经度',
  latitude NUMERIC(14,10) COMMENT '维度',
  address VARCHAR(255) COMMENT '商户地址',
  cover_distance INT COMMENT '覆盖距离 单位米',
  weight int COMMENT '商户权重',
  status INT COMMENT '商户状态 1.正常 0.关闭',
  flag VARCHAR(255) COMMENT '商户标记',
  json VARCHAR(1000) COMMENT '附加数据',
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间戳'

)CHARACTER SET utf8mb4;

-- 商户图片
CREATE TABLE IF NOT EXISTS merchant_imgs(
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  app_id VARCHAR(255) COMMENT 'APPID',
  open_id VARCHAR(255) COMMENT 'open_id',
  merchant_id BIGINT COMMENT '商户ID',
  flag VARCHAR(100) COMMENT '图片标记',
  url VARCHAR(400) COMMENT '图片URL',
  json VARCHAR(1000) COMMENT '附加字段'
) CHARACTER SET utf8mb4;


-- select getDistance(116.3899,39.91578,116.3904,39.91576);

-- 商户产品
CREATE TABLE IF NOT EXISTS merchant_prod(
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  app_id VARCHAR(255) COMMENT 'APPID',
  merchant_id BIGINT COMMENT '商户ID',
  prod_id BIGINT COMMENT '产品ID',
  flag  VARCHAR(255) COMMENT '标记',
  json VARCHAR(1000) COMMENT '附加字段',
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间戳'
) CHARACTER SET utf8mb4;

-- 类别
CREATE TABLE IF NOT EXISTS category(
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  app_id VARCHAR(255) COMMENT 'APPID',
  title VARCHAR(255) COMMENT '标题',
  description VARCHAR(255) COMMENT '描述',
  icon VARCHAR(255) COMMENT '图标',
  flag VARCHAR(255) COMMENT '标记',
  json VARCHAR(1000) COMMENT '附加字段',
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间戳'
) CHARACTER SET utf8mb4;

-- INSERT INTO category(app_id, title, description, flag,json) VALUES ('hwl','默认分类','默认分类描述','default','');

-- 商品
CREATE TABLE IF NOT EXISTS product(
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  app_id VARCHAR(255) COMMENT 'APPID',
  title VARCHAR(255) COMMENT '标题',
  description VARCHAR(1000) COMMENT '描述',
  price NUMERIC(14,2) COMMENT '原价',
  dis_price NUMERIC(14,2) COMMENT '折扣价格',
  status int COMMENT '商品状态',
  is_recom int COMMENT '是否推荐 1.推荐 0.不推荐',
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间戳',
  flag VARCHAR(255) COMMENT '标记',
  json VARCHAR(1000) COMMENT '附加字段'

) CHARACTER SET utf8mb4;

-- 商品图片
CREATE TABLE IF NOT EXISTS prod_imgs(
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  app_id VARCHAR(255) COMMENT 'APPID',
  prod_id BIGINT COMMENT '产品ID',
  flag VARCHAR(100) COMMENT '图片标记',
  url VARCHAR(400) COMMENT '图片URL',
  json VARCHAR(1000) COMMENT '附加字段'
) CHARACTER SET utf8mb4;

--  商品属性key
CREATE TABLE IF NOT EXISTS prod_attr_key(
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  prod_id BIGINT COMMENT '商品ID',
  attr_key VARCHAR(255) COMMENT '属性唯一key',
  attr_name VARCHAR(255) COMMENT '属性名',
  status int COMMENT '1.正常 0.关闭',
  flag VARCHAR(100) COMMENT '标记',
  json VARCHAR(1000) COMMENT '附加字段'

) CHARACTER SET utf8mb4;

-- 商品属性值
CREATE TABLE IF NOT EXISTS prod_attr_val(
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  prod_id BIGINT COMMENT '商品ID',
  attr_key VARCHAR(255)  COMMENT '属性key',
  attr_value VARCHAR(255) COMMENT '属性值',
  flag VARCHAR(100) COMMENT '标记',
  json VARCHAR(1000) COMMENT '附加字段'

) CHARACTER SET utf8mb4;

--  商品sku
CREATE TABLE IF NOT EXISTS prod_sku(

  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  sku_no VARCHAR(255) COMMENT '唯一编号',
  prod_id BIGINT COMMENT '商品ID',
  app_id VARCHAR(255) COMMENT 'APPID',
  price NUMERIC(14,2) COMMENT '原价',
  dis_price NUMERIC(14,2) COMMENT '折扣价格',
  attr_symbol_path VARCHAR(255) COMMENT '属性组合出的规格路径',
  stock int COMMENT '库存量',
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间戳',
  flag VARCHAR(100) COMMENT '标记',
  json VARCHAR(1000) COMMENT '附加字段'

) CHARACTER SET utf8mb4;

-- 商品分类
CREATE TABLE IF NOT EXISTS prod_category (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  category_id BIGINT COMMENT '类别ID',
  app_id VARCHAR(255) COMMENT 'APPID',
  prod_id BIGINT COMMENT '商品ID',
  json VARCHAR(1000) COMMENT '附加字段',
  flag VARCHAR(100) COMMENT '标记',
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间戳'
) CHARACTER SET utf8mb4;

-- 订单
CREATE TABLE IF NOT EXISTS `order` (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  no VARCHAR(255)  COMMENT '订单编号',
  code VARCHAR(255) COMMENT '预付款编号',
  address_id VARCHAR(255) COMMENT '地址ID',
  address VARCHAR(255) COMMENT '配送地址',
  payapi_no VARCHAR(255) COMMENT '支付中心的订单号',
  m_open_id VARCHAR(255) COMMENT '商户ID',
  open_id VARCHAR(255) COMMENT '用户ID',
  app_id VARCHAR(255) COMMENT 'APPID',
  title VARCHAR(255) COMMENT '订单标题',
  act_price NUMERIC(14,2) COMMENT '订单实际金额',
  omit_money NUMERIC(10,4) COMMENT '省略金额',
  price NUMERIC(14,2) COMMENT '订单应付金额',
  order_status int COMMENT '订单状态 0，未确认；1，已确认；2，已取消；3，无效；4，退货',
  pay_status int COMMENT '付款状态 支付状态；0，未付款；2，付款中；1，已付款',
  shipping_fee  decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '配送费用',
  json VARCHAR(1000) COMMENT '附加字段',
  flag VARCHAR(100) COMMENT '标记',
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间戳',
  UNIQUE (`no`),
  KEY `order_status` (`order_status`),
  KEY `pay_status` (`pay_status`)
) CHARACTER SET utf8mb4;

-- 对订单操作日志表
CREATE TABLE IF NOT EXISTS order_action(
  id mediumint(8) unsigned  PRIMARY KEY AUTO_INCREMENT,
  order_no VARCHAR(30) COMMENT '订单号',
  action_open_id VARCHAR(30) COMMENT '操作用户openID',
  order_status int COMMENT '订单状态 0，未确认；1，已确认；2，已取消；3，无效；4，退货',
  pay_status int COMMENT '付款状态 支付状态；0，未付款；2，付款中；1，已付款',
  `action_note` varchar(255)  COMMENT '操作备注',
  `action_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP  COMMENT '操作时间',
  key order_no (order_no)
) CHARACTER SET utf8mb4 COMMENT='对订单操作日志表';


-- 订单项
CREATE TABLE IF NOT EXISTS order_item (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  `no` VARCHAR(255) COMMENT '订单编号',
  app_id VARCHAR(255) COMMENT 'APPID',
  open_id VARCHAR(255) COMMENT '用户ID',
  m_open_id VARCHAR(255) COMMENT '商家ID',
  prod_id BIGINT COMMENT '商品ID',
  sku_no VARCHAR(255) COMMENT '商品SKU编号',
  num int COMMENT '商品数量',
  offer_unit_price NUMERIC(14,2) COMMENT '单价报价',
  buy_unit_price NUMERIC(14,2) COMMENT '购买单价',
  offer_total_price NUMERIC(14,2) COMMENT '总价格报价',
  buy_total_price NUMERIC(14,2) COMMENT '购买总金额',
  flag VARCHAR(100) COMMENT '标记',
  json VARCHAR(1000) COMMENT '附加字段',
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间戳'
)CHARACTER SET utf8mb4;


-- 订单地址
CREATE TABLE IF NOT EXISTS order_address (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  order_no VARCHAR(255) COMMENT '订单号',
  app_id VARCHAR(255) COMMENT 'APPID',
  open_id VARCHAR(255) COMMENT '用户ID',
  name VARCHAR(255) COMMENT '姓名',
  mobile VARCHAR(255) COMMENT '手机号',
  address VARCHAR(255) COMMENT '送货地址',
  flag VARCHAR(100) COMMENT '标记',
  json VARCHAR(1000) COMMENT '附加字段',
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间戳'
) CHARACTER SET utf8mb4;



-- +migrate StatementBegin
CREATE  FUNCTION  `getDistance`(
   lon1 float(10,7)
  ,lat1 float(10,7)
  ,lon2 float(10,7)
  ,lat2 float(10,7)
) RETURNS double
  begin
    declare d double;
    declare radius int;
    set radius = 6378140; #假设地球为正球形，直径为6378140米
    set d = (2*ATAN2(SQRT(SIN((lat1-lat2)*PI()/180/2)
                          *SIN((lat1-lat2)*PI()/180/2)+
                          COS(lat2*PI()/180)*COS(lat1*PI()/180)
                          *SIN((lon1-lon2)*PI()/180/2)
                          *SIN((lon1-lon2)*PI()/180/2)),
                     SQRT(1-SIN((lat1-lat2)*PI()/180/2)
                            *SIN((lat1-lat2)*PI()/180/2)
                          +COS(lat2*PI()/180)*COS(lat1*PI()/180)
                           *SIN((lon1-lon2)*PI()/180/2)
                           *SIN((lon1-lon2)*PI()/180/2))))*radius;
    return d;
  end;
-- +migrate StatementEnd