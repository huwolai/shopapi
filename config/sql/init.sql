-- +migrate Up
-- 地址信息
CREATE TABLE IF NOT EXISTS address(
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  app_id VARCHAR(255) COMMENT 'APPID',
  open_id VARCHAR(255) COMMENT '用户ID',
  address VARCHAR(255) COMMENT '地址',
  is_default VARCHAR(255) COMMENT '0不是默认,1是默认地址',
  json VARCHAR(1000) COMMENT '附加字段',
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间戳'

)CHARACTER SET utf8mb4;

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
  json VARCHAR(1000) COMMENT '附加数据',
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间戳'

)CHARACTER SET utf8mb4;

drop function getDistance;
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
-- select getDistance(116.3899,39.91578,116.3904,39.91576);

-- 商户产品
CREATE TABLE IF NOT EXISTS merchant_prod(
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  app_id VARCHAR(255) COMMENT 'APPID',
  merchant_id BIGINT COMMENT '商户ID',
  prod_id BIGINT COMMENT '产品ID',
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
  json VARCHAR(1000) COMMENT '附加字段',
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间戳'
) CHARACTER SET utf8mb4;

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


--  商品规格
CREATE TABLE IF NOT EXISTS prod_spec(

  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  prod_id BIGINT COMMENT '商品ID',
  app_id VARCHAR(255) COMMENT 'APPID',
  price_move NUMERIC(14,2) COMMENT '价格波动',
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间戳',
  json VARCHAR(1000) COMMENT '附加字段'

) CHARACTER SET utf8mb4;

-- 商品分类
CREATE TABLE IF NOT EXISTS prod_category (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  category_id BIGINT COMMENT '类别ID',
  app_id VARCHAR(255) COMMENT 'APPID',
  prod_id BIGINT COMMENT '商品ID',
  json VARCHAR(1000) COMMENT '附加字段',
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间戳'
) CHARACTER SET utf8mb4;

-- 订单
CREATE TABLE IF NOT EXISTS `order` (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  no VARCHAR(255)  COMMENT '订单编号',
  payapi_no VARCHAR(255) COMMENT '支付中心的订单号',
  open_id VARCHAR(255) COMMENT '用户ID',
  app_id VARCHAR(255) COMMENT 'APPID',
  title VARCHAR(255) COMMENT '订单标题',
  act_price NUMERIC(14,2) COMMENT '订单实际金额',
  omit_money NUMERIC(10,4) COMMENT '省略金额',
  price NUMERIC(14,2) COMMENT '订单应付金额',
  status int COMMENT '订单状态 0:订单被取消 1:已下单待付款 2:已付款',
  json VARCHAR(1000) COMMENT '附加字段',
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间戳'
) CHARACTER SET utf8mb4;

-- 订单项
CREATE TABLE IF NOT EXISTS order_item (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  `no` VARCHAR(255) COMMENT '订单编号',
  app_id VARCHAR(255) COMMENT 'APPID',
  open_id VARCHAR(255) COMMENT '用户ID',
  m_open_id VARCHAR(255) COMMENT '商家ID',
  prod_id BIGINT COMMENT '商品ID',
  num int COMMENT '商品数量',
  offer_unit_price NUMERIC(14,2) COMMENT '单价报价',
  buy_unit_price NUMERIC(14,2) COMMENT '购买单价',
  offer_total_price NUMERIC(14,2) COMMENT '总价格报价',
  buy_total_price NUMERIC(14,2) COMMENT '购买总金额',
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
  address VARCHAR(255) COMMENT '送货地址',
  json VARCHAR(1000) COMMENT '附加字段',
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间戳'
) CHARACTER SET utf8mb4;


-- 订单事件
CREATE TABLE IF NOT EXISTS order_event (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  order_no VARCHAR(255) COMMENT '订单号',
  app_id VARCHAR(255) COMMENT 'APPID',
  open_id VARCHAR(255) COMMENT '用户ID',
  event_type INT COMMENT '事件类型',
  event_name VARCHAR(100) COMMENT '事件名',
  event_desc VARCHAR(255) COMMENT '事件描述',
  json VARCHAR(1000) COMMENT '附加字段',
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间戳'
) CHARACTER SET utf8mb4;