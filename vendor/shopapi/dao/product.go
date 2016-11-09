package dao

import (
	"github.com/gocraft/dbr"
	"gitlab.qiyunxin.com/tangtao/utils/db"
	"gitlab.qiyunxin.com/tangtao/utils/log"
)

type Product struct  {
	Id int64
	AppId string
	//商品标题
	Title string
	//子标题
	SubTitle string
	//商品描述
	Description string
	//商品价格
	Price float64
	//折扣价格
	DisPrice float64
	//是否推荐
	IsRecom int
	//已售数量
	SoldNum int
	//商品状态
	Status int
	Flag string
	//附加数据
	Json string
}

type ProductDetail struct {
	//商品ID
	Id int64
	AppId string
	//商品描述
	Description string
	//分类ID
	CategoryId int64
	//分类名
	CategoryName string
	//商品标题
	Title string
	//子标题
	SubTitle string
	//商品价格
	Price float64
	//折扣价格
	DisPrice float64
	//商品状态
	Status int
	//是否推荐
	IsRecom int
	//已售数量
	SoldNum int
	//商户ID
	MerchantId int64
	//商户名称
	MerchantName string
	Flag string
	//附加数据
	Json string
	//商品图片集合
	ProdImgs []*ProdImgsDetail
}

type ProductSearch struct {
	Keyword		 string
	Category  	 uint64
	Status  	 uint64
	IsRecom  	 uint64
	PriceUp	 	 float64
	PriceDown 	 float64
}


func NewProductDetail() *ProductDetail {

	return &ProductDetail{}
}

func NewProduct() *Product  {

	return &Product{}
}

//详情集合
func (self *ProductDetail) ProdDetailListWith(keywords interface{} ,merchantId int64,flags []string,noflags []string,orderBy string,pageIndex uint64,pageSize uint64,appId string) ([]*ProductDetail,error)  {
	
	search:=keywords.(ProductSearch)	
	
	var prodList []*ProductDetail
	buider :=db.NewSession().Select("product.*,IFNULL(merchant.id,0) merchant_id,IFNULL(merchant.name,'') merchant_name,IFNULL(category.id,0) category_id,IFNULL(category.title,'') category_name").From("product").LeftJoin("merchant_prod","product.id=merchant_prod.prod_id").LeftJoin("merchant","merchant_prod.merchant_id=merchant.id").LeftJoin("prod_category","prod_category.prod_id=product.id").LeftJoin("category","prod_category.category_id=category.id")
	if flags!=nil{
		buider = buider.Where("product.flag in ?",flags)
	}

	if search.Keyword!="" {
		buider = buider.Where("product.title like ?","%"+search.Keyword+"%")
	}
	//分类类别  ( 1 水果 2干货 3海鲜 4食材 5家常用餐 6经典家宴 7私人订制)
	if search.Category>0 {
		buider = buider.Where("product.id in (select prod_id from prod_category where category_id =?)",search.Category)
	}
	//是否上架状态 ( 1 上架 2 下架)
	if search.Status>0 {
		if search.Status==1 {
			buider = buider.Where("product.status = ?",1)
		}else{
			buider = buider.Where("product.status = ?",0)
		}		
	}
	//是否推荐 ( 1 是 2 否)
	if search.IsRecom>0 {
		if search.IsRecom==1 {
			buider = buider.Where("product.is_recom = ?",1)
		}else{
			buider = buider.Where("product.is_recom = ?",0)
		}		
	}	
	//价格区间 左边(包含)
	if search.PriceUp>0 {
		buider = buider.Where("product.price >= ?",search.PriceUp)
	}
	//价格区间  右边(包含)
	if search.PriceDown>0 {
		buider = buider.Where("product.price <= ?",search.PriceDown)
	}

	if noflags!=nil {
		buider = buider.Where("product.flag not in ?",noflags)
	}
	/* if isRecomm!="" {
		buider = buider.Where("product.is_recomm=?",isRecomm)
	} */
	if merchantId!=0{
		buider = buider.Where("merchant.id=?",merchantId)
	}
	if orderBy!="" {
		buider =buider.OrderDir(orderBy,false)
	}
	_,err :=buider.Limit(pageSize).Offset((pageIndex-1)*pageSize).LoadStructs(&prodList)
	if err!=nil{
		return nil,err
	}
	err = FillProdImgs(appId,prodList)

	return prodList,err
}

func (self *ProductDetail) ProdDetailListCountWith(keywords interface{},merchantId int64,flags []string,noflags []string)  (int64,error) {

	search:=keywords.(ProductSearch)	

	var count int64
	buider :=db.NewSession().Select("count(*)").From("product").LeftJoin("merchant_prod","product.id=merchant_prod.prod_id").LeftJoin("merchant","merchant_prod.merchant_id=merchant.id")


	if flags!=nil{
		buider = buider.Where("product.flag in ?",flags)
	}

	if noflags!=nil {
		buider = buider.Where("product.flag not in ?",noflags)
	}

	if search.Keyword!="" {
		buider = buider.Where("product.title like ?","%"+search.Keyword+"%")
	}
	//分类类别  ( 1 水果 2干货 3海鲜 4食材 5家常用餐 6经典家宴 7私人订制)
	if search.Category>0 {
		buider = buider.Where("product.id in (select prod_id from prod_category where category_id =?)",search.Category)
	}
	//是否上架状态 ( 1 上架 2 下架)
	if search.Status>0 {
		if search.Status==1 {
			buider = buider.Where("product.status = ?",1)
		}else{
			buider = buider.Where("product.status = ?",0)
		}		
	}
	//是否推荐 ( 1 是 2 否)
	if search.IsRecom>0 {
		if search.IsRecom==1 {
			buider = buider.Where("product.is_recom = ?",1)
		}else{
			buider = buider.Where("product.is_recom = ?",0)
		}		
	}	
	//价格区间 左边(包含)
	if search.PriceUp>0 {
		buider = buider.Where("product.price >= ?",search.PriceUp)
	}
	//价格区间  右边(包含)
	if search.PriceDown>0 {
		buider = buider.Where("product.price <= ?",search.PriceDown)
	}

	/* if isRecomm!="" {
		buider = buider.Where("product.is_recomm=?",isRecomm)
	} */

	if merchantId!=0{
		buider = buider.Where("merchant.id=?",merchantId)
	}

	err :=buider.LoadValue(&count)

	return count,err
}

func (self *Product) SoldNumInc(num int,prodId int64,tx *dbr.Tx) error  {

	_,err :=tx.UpdateBySql("update product set sold_num=sold_num+? where id=?",num,prodId).Exec()

	return err
}

func (self *Product) InsertTx(tx *dbr.Tx) (int64,error)  {

	result,err :=tx.InsertInto("product").Columns("title","sub_title","app_id","description","sold_num","price","dis_price","json","flag","status","is_recom").Record(self).Exec()
	if err !=nil {

		return 0,err
	}
	pid,err :=  result.LastInsertId()
	
	_,err =tx.Update("product").Set("order",pid).Where("id=?",pid).Exec()
	
	
	return pid,err
}

func (self *Product) UpdateTx(tx *dbr.Tx) error {
	_,err :=tx.Update("product").Set("title",self.Title).Set("sub_title",self.SubTitle).Set("description",self.Description).Set("price",self.Price).Set("dis_price",self.DisPrice).Set("json",self.Json).Where("id=?",self.Id).Exec()
	return err
}

func (self *Product) WithFlag(flag string,merchantId int64)  ([]*Product,error)  {
	var products []*Product
	_,err :=db.NewSession().SelectBySql("select * from product pt,merchant_prod mp where pt.id = mp.prod_id and pt.status=1 and pt.flag=? and mp.merchant_id=?",flag,merchantId).LoadStructs(&products)

	return products,err
}

//商品推荐列表
func (self *ProductDetail) ProductListWithRecomm(appId string,pageIndex uint64,pageSize uint64) ([]*ProductDetail,int,error)  {
	session := db.NewSession()
	
	var count int
	_,err :=session.SelectBySql("select count(id) from product  where is_recom=1 and status=1 and app_id=? order by `order` desc limit ?,?",appId,(pageIndex-1)*pageSize,pageSize).LoadStructs(&count)
	
	
	var prodList []*ProductDetail
	_,err =session.SelectBySql("select * from product  where is_recom=1 and status=1 and app_id=? order by `order` desc limit ?,?",appId,(pageIndex-1)*pageSize,pageSize).LoadStructs(&prodList)
	if err!=nil{
		log.Debug("----err",err)
		return nil,0,err
	}

	err = FillProdImgs(appId,prodList)

	return prodList,count,err
}

func (self *ProductDetail) ProductListWithMerchant(merchantId int64,appId string,flags []string,noflags []string) ([]*ProductDetail,error)  {
	session := db.NewSession()
	var prodList []*ProductDetail
	var builder *dbr.SelectBuilder
	if flags!=nil&&len(flags)>0&&(noflags==nil||len(noflags)==0) {
		builder = session.SelectBySql("select pt.id,pt.app_id,pt.title,pt.sub_title,pt.price,pt.dis_price,pt.flag,pt.`status`,mt.id merchant_id,mt.`name` merchant_name,pt.json from merchant_prod md,merchant mt,product pt where md.prod_id=pt.id and pt.status=1  and md.merchant_id=mt.id  and mt.id=? and pt.app_id=? and pt.flag in ?",merchantId,appId,flags)
	}

	if noflags!=nil&&len(noflags)>0&&(flags==nil||len(flags)==0) {
		builder = session.SelectBySql("select pt.id,pt.app_id,pt.title,pt.sub_title,pt.price,pt.dis_price,pt.flag,pt.`status`,mt.id merchant_id,mt.`name` merchant_name,pt.json from merchant_prod md,merchant mt,product pt where md.prod_id=pt.id and pt.status=1  and md.merchant_id=mt.id  and mt.id=? and pt.app_id=? and pt.flag not in ?",merchantId,appId,noflags)
	}

	if noflags==nil&&len(noflags)==0&&flags==nil&&len(flags)==0 {
		builder = session.SelectBySql("select pt.id,pt.app_id,pt.title,pt.sub_title,pt.price,pt.dis_price,pt.flag,pt.`status`,mt.id merchant_id,mt.`name` merchant_name,pt.json from merchant_prod md,merchant mt,product pt where md.prod_id=pt.id and pt.status=1  and md.merchant_id=mt.id  and mt.id=? and pt.app_id=?",merchantId,appId)
	}

	if noflags!=nil&&len(noflags)>0&&flags!=nil&&len(flags)>0 {
		builder = session.SelectBySql("select pt.id,pt.app_id,pt.title,pt.sub_title,pt.price,pt.dis_price,pt.flag,pt.`status`,mt.id merchant_id,mt.`name` merchant_name,pt.json from merchant_prod md,merchant mt,product pt where md.prod_id=pt.id and pt.status=1  and md.merchant_id=mt.id  and mt.id=? and pt.app_id=? flag in ? and pt.flag not in ?",merchantId,appId,flags,noflags)
	}
	_,err :=builder.LoadStructs(&prodList)
	if err!=nil{
		return nil,err
	}
	err = FillProdImgs(appId,prodList)

	return prodList,err
}

func (self *ProductDetail) ProductListWithCategory(appId string,categoryId int64,flags []string,noflags []string) ([]*ProductDetail,error)  {
	session := db.NewSession()
	var prodList []*ProductDetail
	builder :=session.Select("product.*,merchant.id merchant_id,merchant.name merchant_name").From("product").Join("prod_category","product.id = prod_category.prod_id").Join("merchant_prod","product.id = merchant_prod.prod_id").Join("merchant","merchant.id = merchant_prod.merchant_id").Where("prod_category.category_id=?",categoryId).Where("product.status=?",1).Where("product.app_id=?",appId)
	if flags!=nil&&len(flags)>0{

		builder = builder.Where("product.flag in ?",flags)
	}
	if noflags!=nil&&len(noflags) >0 {
		builder = builder.Where("product.flag not in ?",noflags)
	}
	_,err := builder.LoadStructs(&prodList)
	if err!=nil{
		return nil,err
	}
	if prodList!=nil&&len(prodList)>0 {
		err = FillProdImgs(appId,prodList)
	}


	return prodList,err
}

//填充商品图片数据
func FillProdImgs(appId string,prodList []*ProductDetail) error {
	prodids := make([]int64,0)
	if prodList!=nil{
		for _,prod :=range prodList {
			prodids = append(prodids,prod.Id)
		}
	}

	if len(prodids)<=0 {
		return nil
	}

	prodImgDetail := NewProdImgsDetail()
	prodImgDetails,err := prodImgDetail.ProdImgsWithProdIds(prodids)
	if err!=nil{
		return err
	}
	prodimgsMap := make(map[int64][]*ProdImgsDetail)
	if prodImgDetails!=nil{
		for _,prodimgd :=range prodImgDetails {
			key := prodimgd.ProdId
			pdimgdetails :=prodimgsMap[key]
			if pdimgdetails==nil{
				pdimgdetails = make([]*ProdImgsDetail,0)
			}


			pdimgdetails= append(pdimgdetails,prodimgd)

			prodimgsMap[key] = pdimgdetails
			log.Debug(prodimgsMap)
		}
	}
	for _,prod :=range prodList {
		key := prod.Id
		prodimgs := prodimgsMap[key]
		prod.ProdImgs = prodimgs
	}

	return nil
}

func (self *Product) ProductDetailWithId(id int64) (*ProductDetail,error) {
	var prodDetail *ProductDetail
	_,err :=db.NewSession().SelectBySql("select pt.*,mt.`name` merchant_name,mt.id merchant_id,pct.category_id from product pt left join merchant_prod md on pt.id=md.prod_id LEFT JOIN merchant mt on md.merchant_id=mt.id left join prod_category pct on pct.prod_id=pt.id WHERE pt.id=md.prod_id and md.merchant_id=mt.id and pt.id=?",id).LoadStructs(&prodDetail)

	return prodDetail,err
}

func (self *Product) ProductWithId(id int64,appId string) (*Product,error)  {
	sess :=db.NewSession()
	var prod *Product
	_,err :=sess.Select("*").From("product").Where("app_id=?",appId).Where("id=?",id).LoadStructs(&prod)

	return prod,err
}

func (self *Product) UpdateStatusWithProdId(status int,prodId int64) error  {

	_,err :=db.NewSession().Update("product").Set("status",status).Where("id=?",prodId).Exec()

	return err
}

func (self *Product) UpdateRecomWithProdId(isRecom int,prodId int64) error  {
	_,err :=db.NewSession().Update("product").Set("is_recom",isRecom).Where("id=?",prodId).Exec()

	return err
}