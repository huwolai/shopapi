-- +migrate Up
-- 地址信息
CREATE TABLE IF NOT EXISTS address(
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  app_id VARCHAR(255) DEFAULT '' COMMENT 'APPID',
  open_id VARCHAR(100) DEFAULT '' COMMENT '用户ID',
  longitude NUMERIC(14,10) COMMENT '经度',
  latitude NUMERIC(14,10) COMMENT '维度',
  name VARCHAR(255) DEFAULT '' COMMENT '姓名',
  mobile VARCHAR(255) DEFAULT '' COMMENT '手机号',
  address VARCHAR(255) DEFAULT '' COMMENT '地址',
  weight int COMMENT '地址权重 1.权重越高 越优先',
  flag VARCHAR(100) DEFAULT '' COMMENT '标记',
  json VARCHAR(1000) DEFAULT '' COMMENT '附加字段',
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间戳',
  KEY open_id (open_id)

)CHARACTER SET utf8mb4;

-- 账户信息
CREATE TABLE IF NOT EXISTS account(
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  app_id VARCHAR(255) DEFAULT '' COMMENT 'APPID',
  open_id VARCHAR(255) DEFAULT '' COMMENT '用户ID',
  mobile VARCHAR(255) DEFAULT '' COMMENT '手机号',
  money NUMERIC(12,2) COMMENT '账户金额',
  password VARCHAR(255) DEFAULT '' COMMENT '账户密码',
  status int COMMENT '状态 1.正常 0.锁定 2.等待开通支付',
  flag VARCHAR(100) DEFAULT '' COMMENT '标记',
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间戳',
  KEY open_id (open_id)
);


-- 商户
CREATE TABLE IF NOT EXISTS merchant(

  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(255) DEFAULT '' COMMENT '商户名称',
  app_id VARCHAR(255) DEFAULT '' COMMENT 'APPID',
  open_id VARCHAR(255) DEFAULT '' COMMENT '商户open_id',
  longitude NUMERIC(14,10)  COMMENT '经度',
  latitude NUMERIC(14,10) COMMENT '维度',
  address VARCHAR(255) DEFAULT '' COMMENT '商户地址',
  cover_distance INT COMMENT '覆盖距离 单位米',
  weight int COMMENT '商户权重',
  status INT COMMENT '商户状态 1.正常 0.关闭',
  flag VARCHAR(255) DEFAULT '' COMMENT '商户标记',
  json VARCHAR(1000) DEFAULT '' COMMENT '附加数据',
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间戳'

)CHARACTER SET utf8mb4;

-- 商户营业时间
CREATE TABLE IF NOT EXISTS merchant_open(
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  app_id VARCHAR(255) DEFAULT '' COMMENT 'APPID',
  merchant_id INT NOT NULL unique COMMENT '商户ID',
  is_open int COMMENT '是否营业',
  open_time_start VARCHAR(30) DEFAULT '' COMMENT '营业开始时间',
  open_time_end VARCHAR(30) DEFAULT '' COMMENT '营业结束时间',
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间戳'

)CHARACTER SET utf8mb4;

-- 商户每天服务时间(特殊表,不是标准电商表)
CREATE TABLE IF NOT EXISTS merchant_service_time(
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  merchant_id BIGINT COMMENT '商户ID',
  stime VARCHAR(20) COMMENT '时间'
) CHARACTER SET utf8mb4;

-- 商户图片
CREATE TABLE IF NOT EXISTS merchant_imgs(
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  app_id VARCHAR(255) DEFAULT '' COMMENT 'APPID',
  open_id VARCHAR(255) DEFAULT '' COMMENT 'open_id',
  merchant_id BIGINT COMMENT '商户ID',
  flag VARCHAR(100) DEFAULT '' COMMENT '图片标记',
  url VARCHAR(400) DEFAULT '' COMMENT '图片URL',
  json VARCHAR(1000) DEFAULT '' COMMENT '附加字段'
) CHARACTER SET utf8mb4;


-- select getDistance(116.3899,39.91578,116.3904,39.91576);

-- 商户产品
CREATE TABLE IF NOT EXISTS merchant_prod(
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  app_id VARCHAR(255) DEFAULT '' COMMENT 'APPID',
  merchant_id BIGINT COMMENT '商户ID',
  prod_id BIGINT COMMENT '产品ID',
  flag  VARCHAR(255 ) DEFAULT '' COMMENT '标记',
  json VARCHAR(1000) DEFAULT '' COMMENT '附加字段',
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间戳'
) CHARACTER SET utf8mb4;

-- 类别
CREATE TABLE IF NOT EXISTS category(
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  app_id VARCHAR(255) DEFAULT '' COMMENT 'APPID',
  title VARCHAR(255) DEFAULT '' COMMENT '标题',
  description VARCHAR(255) DEFAULT '' COMMENT '描述',
  icon VARCHAR(255) DEFAULT '' COMMENT '图标',
  flag VARCHAR(255) DEFAULT '' COMMENT '标记',
  json VARCHAR(1000) DEFAULT '' COMMENT '附加字段',
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间戳'
) CHARACTER SET utf8mb4;

 -- INSERT INTO category(app_id, title, description, flag,json) VALUES ('hwl','默认分类','默认分类描述','default','');

-- 商品
CREATE TABLE IF NOT EXISTS product(
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  app_id VARCHAR(255) DEFAULT '' COMMENT 'APPID',
  title VARCHAR(255) DEFAULT '' COMMENT '标题',
  description VARCHAR(1000) DEFAULT '' COMMENT '描述',
  price NUMERIC(14,2) COMMENT '原价',
  dis_price NUMERIC(14,2) COMMENT '折扣价格',
  status int COMMENT '商品状态',
  is_recom int COMMENT '是否推荐 1.推荐 0.不推荐',
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间戳',
  flag VARCHAR(255) DEFAULT '' COMMENT '标记',
  json VARCHAR(1000) DEFAULT '' COMMENT '附加字段'

) CHARACTER SET utf8mb4;

-- 商品图片
CREATE TABLE IF NOT EXISTS prod_imgs(
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  app_id VARCHAR(255) DEFAULT '' COMMENT 'APPID',
  prod_id BIGINT COMMENT '产品ID',
  flag VARCHAR(100) DEFAULT '' COMMENT '图片标记',
  url VARCHAR(400) DEFAULT '' COMMENT '图片URL',
  json VARCHAR(1000) DEFAULT '' COMMENT '附加字段'
) CHARACTER SET utf8mb4;


--  商品属性key
CREATE TABLE IF NOT EXISTS prod_attr_key(
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  prod_id BIGINT COMMENT '商品ID',
  attr_key VARCHAR(255) DEFAULT '' COMMENT '属性唯一key',
  attr_name VARCHAR(255) DEFAULT '' COMMENT '属性名',
  status int COMMENT '1.正常 0.关闭',
  flag VARCHAR(100) DEFAULT '' COMMENT '标记',
  json VARCHAR(1000) DEFAULT '' COMMENT '附加字段'

) CHARACTER SET utf8mb4;

-- 商品属性值
CREATE TABLE IF NOT EXISTS prod_attr_val(
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  prod_id BIGINT COMMENT '商品ID',
  attr_key VARCHAR(255) DEFAULT ''  COMMENT '属性key',
  attr_value VARCHAR(255) DEFAULT '' COMMENT '属性值',
  flag VARCHAR(100) DEFAULT '' COMMENT '标记',
  json VARCHAR(1000) DEFAULT '' COMMENT '附加字段'

) CHARACTER SET utf8mb4;

-- 订单折扣
CREATE TABLE IF NOT EXISTS order_discount(
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  t_coupon_code VARCHAR(100) DEFAULT ''COMMENT '第三方券代号',
  order_no VARCHAR(100) DEFAULT '' COMMENT '订单号',
  act_price NUMERIC(14,2)  COMMENT '实际价格',
  dis_price NUMERIC(14,2) COMMENT '折扣价格',
  dis_money NUMERIC(10,2) COMMENT '折掉的金额',
  flag VARCHAR(100) DEFAULT '' COMMENT '标记',
  json VARCHAR(1000)  DEFAULT '' COMMENT '附加字段',
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间戳'

) CHARACTER SET utf8mb4;

--  商品sku
CREATE TABLE IF NOT EXISTS prod_sku(

  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  sku_no VARCHAR(255) DEFAULT '' COMMENT '唯一编号',
  prod_id BIGINT COMMENT '商品ID',
  app_id VARCHAR(255) DEFAULT '' COMMENT 'APPID',
  price NUMERIC(14,2) COMMENT '原价',
  dis_price NUMERIC(14,2) COMMENT '折扣价格',
  attr_symbol_path VARCHAR(255) DEFAULT '' COMMENT '属性组合出的规格路径',
  stock int COMMENT '库存量',
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间戳',
  flag VARCHAR(100) DEFAULT '' COMMENT '标记',
  json VARCHAR(1000)  DEFAULT '' COMMENT '附加字段'

) CHARACTER SET utf8mb4;

-- 商品分类
CREATE TABLE IF NOT EXISTS prod_category (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  category_id BIGINT COMMENT '类别ID',
  app_id VARCHAR(255) DEFAULT '' COMMENT 'APPID',
  prod_id BIGINT COMMENT '商品ID',
  json VARCHAR(1000) DEFAULT '' COMMENT '附加字段',
  flag VARCHAR(100) DEFAULT '' COMMENT '标记',
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间戳'
) CHARACTER SET utf8mb4;

-- 订单
CREATE TABLE IF NOT EXISTS `order` (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  no VARCHAR(30)  DEFAULT '' COMMENT '订单编号',
  code VARCHAR(255) DEFAULT '' COMMENT '预付款编号',
  address_id VARCHAR(255) DEFAULT '' COMMENT '地址ID',
  address VARCHAR(255) DEFAULT '' COMMENT '配送地址',
  payapi_no VARCHAR(255) DEFAULT '' COMMENT '支付中心的订单号',
  merchant_id VARCHAR(255) DEFAULT '' COMMENT '商户ID',
  m_open_id VARCHAR(255) DEFAULT '' COMMENT '商户OpenId',
  open_id VARCHAR(255) DEFAULT '' COMMENT '用户ID',
  app_id VARCHAR(255) DEFAULT '' COMMENT 'APPID',
  title VARCHAR(255) DEFAULT '' COMMENT '订单标题',
  act_price NUMERIC(14,2)  COMMENT '订单实际金额(此金额为实际付款金额)',
  omit_money NUMERIC(10,4) COMMENT '省略金额',
  price NUMERIC(14,2) COMMENT '订单金额',
  order_status int COMMENT '订单状态 0，未确认；1，已确认；2，已取消；3，无效；4，退货',
  pay_status int COMMENT '付款状态 支付状态；0，未付款；2，付款中；1，已付款',
  shipping_fee  decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '配送费用',
  json VARCHAR(1000) DEFAULT '' COMMENT '附加字段',
  flag VARCHAR(100) DEFAULT '' COMMENT '标记',
  cancel_reason VARCHAR(1000) DEFAULT '' COMMENT '取消订单原因',
  reject_cancel_reason VARCHAR(1000) DEFAULT '' COMMENT '拒绝取消订单的原因(商户方)',
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间戳',
  UNIQUE (`no`),
  KEY `order_status` (`order_status`),
  KEY `pay_status` (`pay_status`)
) CHARACTER SET utf8mb4;

-- 对订单操作日志表
CREATE TABLE IF NOT EXISTS order_action(
  id mediumint(8) unsigned  PRIMARY KEY AUTO_INCREMENT,
  order_no VARCHAR(30) DEFAULT '' COMMENT '订单号',
  action_open_id VARCHAR(30) DEFAULT '' COMMENT '操作用户openID',
  order_status int COMMENT '订单状态 0，未确认；1，已确认；2，已取消；3，无效；4，退货',
  pay_status int COMMENT '付款状态 支付状态；0，未付款；2，付款中；1，已付款',
  `action_note` varchar(255)  DEFAULT '' COMMENT '操作备注',
  `action_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP  COMMENT '操作时间',
  key order_no (order_no)
) CHARACTER SET utf8mb4 COMMENT='对订单操作日志表';


-- 订单项
CREATE TABLE IF NOT EXISTS order_item (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  `no` VARCHAR(30) DEFAULT '' COMMENT '订单编号',
  app_id VARCHAR(255) DEFAULT '' COMMENT 'APPID',
  open_id VARCHAR(255)  DEFAULT '' COMMENT '用户ID',
  m_open_id VARCHAR(255) DEFAULT '' COMMENT '商家ID',
  prod_id BIGINT COMMENT '商品ID',
  sku_no VARCHAR(255) DEFAULT '' COMMENT '商品SKU编号',
  num int COMMENT '商品数量',
  offer_unit_price NUMERIC(14,2) COMMENT '单价报价',
  buy_unit_price NUMERIC(14,2) COMMENT '购买单价',
  offer_total_price NUMERIC(14,2) COMMENT '总价格报价',
  buy_total_price NUMERIC(14,2) COMMENT '购买总金额',
  flag VARCHAR(100) DEFAULT '' COMMENT '标记',
  json VARCHAR(1000) DEFAULT '' COMMENT '附加字段',
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间戳'
)CHARACTER SET utf8mb4;


-- 订单地址
CREATE TABLE IF NOT EXISTS order_address (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  order_no VARCHAR(255) DEFAULT '' COMMENT '订单号',
  app_id VARCHAR(255) DEFAULT '' COMMENT 'APPID',
  open_id VARCHAR(255) DEFAULT '' COMMENT '用户ID',
  name VARCHAR(255) DEFAULT '' COMMENT '姓名',
  mobile VARCHAR(255) DEFAULT '' COMMENT '手机号',
  address VARCHAR(255) DEFAULT '' COMMENT '送货地址',
  flag VARCHAR(100) DEFAULT '' COMMENT '标记',
  json VARCHAR(1000) DEFAULT '' COMMENT '附加字段',
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间戳'
) CHARACTER SET utf8mb4;



# INSERT INTO category(app_id, title, description, icon, flag) VALUES ('shopapi','家常用餐','家常菜','../static/area_1.png','home');
# INSERT INTO category(app_id, title, description, icon, flag) VALUES ('shopapi','系列套餐','私人订制','../static/area_2.png','home');
# INSERT INTO category(app_id, title, description, icon, flag) VALUES ('shopapi','私人订制','家常菜','../static/area_3.png','home');

# INSERT INTO category(app_id, title, description, icon, flag) VALUES ('shopapi','优惠','优惠','../../static/mall-icon1.png','');
# INSERT INTO category(app_id, title, description, icon, flag) VALUES ('shopapi','促销','促销','../static/mall-icon2.png','');
# INSERT INTO category(app_id, title, description, icon, flag) VALUES ('shopapi','海鲜','海鲜','../static/mall-icon3.png','');
# INSERT INTO category(app_id, title, description, icon, flag) VALUES ('shopapi','食材','食材','../static/mall-icon4.png','');

# INSERT INTO product(app_id, title, description, price, dis_price, status, is_recom) VALUES ('shopapi','测试商品','测试商品',0.01,0.01,1,1);
# INSERT INTO prod_attr_key(prod_id, attr_key, attr_name, status) VALUES (1,'time','时间',1);
#INSERT INTO product(app_id, title, description, price, dis_price, status, is_recom) VALUES ('shopapi','测试商品2','测试商品2',0.02,0.02,1,1);
#INSERT INTO prod_category(category_id, app_id, prod_id) VALUES (4,'shopapi',2);
#
# INSERT INTO prod_imgs(app_id, prod_id, flag, url, json) VALUES ('shopapi',1,'','http://img3.redocn.com/tupian/20141029/yipinyangpaimeishi_3346599.jpg','');
# INSERT INTO prod_imgs(app_id, prod_id, flag, url, json) VALUES ('shopapi',1,'','http://pic47.nipic.com/20140909/11902156_133459495000_2.jpg','');
#INSERT INTO merchant_prod(app_id, merchant_id, prod_id) VALUES ('shopapi',1,1);
#INSERT INTO prod_sku(sku_no, prod_id, app_id, price, dis_price, attr_symbol_path, stock) VALUES ('1234',1,'shopapi',0.01,0.01,'',1);

# INSERT INTO prod_imgs(app_id, prod_id, flag, url, json) VALUES ('shopapi',5,'','http://img3.imgtn.bdimg.com/it/u=3794806978,249039065&fm=11&gp=0.jpg','');
# INSERT INTO prod_imgs(app_id, prod_id, flag, url, json) VALUES ('shopapi',5,'','http://s.qdcdn.com/cl/11527030,800,450.jpg','');
# INSERT INTO merchant_prod(app_id, merchant_id, prod_id) VALUES ('shopapi',1,5);
# INSERT INTO prod_sku(sku_no, prod_id, app_id, price, dis_price, attr_symbol_path, stock) VALUES ('12345',5,'shopapi',0.02,0.02,'',1);

# INSERT INTO product(app_id, title, description, price, dis_price, status, is_recom) VALUES ('shopapi','家用套餐','家用套餐',0.05,0.05,1,1);
# INSERT INTO prod_imgs(app_id, prod_id, flag, url, json) VALUES ('shopapi',6,'','http://img0.imgtn.bdimg.com/it/u=2088425899,2259442035&fm=21&gp=0.jpg','');
# INSERT INTO prod_imgs(app_id, prod_id, flag, url, json) VALUES ('shopapi',6,'','http://img3.redocn.com/tupian/20141107/gaolushaobingmeishi_3416757.jpg','');
# INSERT INTO merchant_prod(app_id, merchant_id, prod_id) VALUES ('shopapi',3,6);
# INSERT INTO prod_sku(sku_no, prod_id, app_id, price, dis_price, attr_symbol_path, stock) VALUES ('123456',6,'shopapi',0.02,0.02,'',1);
INSERT INTO prod_category(category_id, app_id, prod_id) VALUES (1,'shopapi',6);
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