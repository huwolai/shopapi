package dao

import (
	"gitlab.qiyunxin.com/tangtao/utils/db"
	//"errors"
	"github.com/gocraft/dbr"
)

type Account struct  {
	Id 			int64
	AppId 		string
	OpenId 		string
	Mobile 		string
	Money 		float64
	Password 	string
	Status 		int
	YdgyId		string
	YdgyName 	string
	YdgyStatus	int64
	Name 		string
	FreezeMoney int64
	Getui 		string
	BaseDModel
}

type GetOnKey struct  {
	Status int
}

func NewAccount() *Account  {

	return &Account{}
}

func NewGetOnKey() *GetOnKey  {

	return &GetOnKey{}
}

func (self *Account) Insert() error {

	_,err :=db.NewSession().InsertInto("account").Columns("app_id","open_id","mobile","money","password","status").Record(self).Exec()

	return err
}

func (self *Account) AccountWithMobile(mobile string,appId string) (*Account,error) {
	var account *Account
	_,err :=db.NewSession().Select("*").From("account").Where("mobile=?",mobile).Where("app_id=?",appId).LoadStructs(&account)

	return account,err
}

func (self *Account) AccountWithOpenId(openId string,appId string) (*Account,error)  {

	var account *Account
	_,err :=db.NewSession().Select("*").From("account").Where("open_id=?",openId).Where("app_id=?",appId).LoadStructs(&account)

	return account,err
}

func (self *Account) AccountUpdatePwd(pwd string,openId,appId string) error {

	_,err :=db.NewSession().Update("account").Set("password",pwd).Where("open_id=?",openId).Where("app_id=?",appId).Exec()

	return err
}

func (self *Account) AccountUpdateStatus(status int,openId string,appId string) error  {
	_,err :=db.NewSession().Update("account").Set("status",status).Where("open_id=?",openId).Where("app_id=?",appId).Exec()

	return err
}

//查询用户
func (self *Account) AccountsWith(pageIndex uint64,pageSize uint64,mobile string,appId string,userName string,ydgyId string,ydgyName string,ydgyStatus string) ([]*Account,error)  {
	var list []*Account
	builder :=db.NewSession().Select("*").From("account").Where("app_id=?",appId).OrderDir("create_time",false)
	if mobile!=""{
		builder = builder.Where("mobile like ?",mobile + "%")
	}
	if ydgyId!=""{
		builder = builder.Where("ydgy_id like ?","%" + ydgyId + "%")
	}
	if ydgyName!=""{
		builder = builder.Where("ydgy_name like ?","%" + ydgyName + "%")
	}
	if userName!=""{
		builder = builder.Where("open_id in (select open_id from merchant where name like ? )","%" + userName + "%")
	}
	if ydgyStatus!=""{
		builder = builder.Where("ydgy_status = ?",ydgyStatus)
	}
	_,err :=builder.Limit(pageSize).Offset((pageIndex-1)*pageSize).LoadStructs(&list)
	
	return list,err
}

func (self *Account) AccountsWithCount(mobile string,appId string,userName string,ydgyId string,ydgyName string,ydgyStatus string) (int64,error) {
	builder :=db.NewSession().Select("count(*)").From("account").Where("app_id=?",appId)
	if mobile!=""{
		builder = builder.Where("mobile like ?",mobile + "%")
	}
	if ydgyId!=""{
		builder = builder.Where("ydgy_id like ?","%" + ydgyId + "%")
	}
	if ydgyName!=""{
		builder = builder.Where("ydgy_name like ?","%" + ydgyName + "%")
	}
	if userName!=""{
		builder = builder.Where("open_id in (select open_id from merchant where name like ? )","%" + userName + "%")
	}
	if ydgyStatus!=""{
		builder = builder.Where("ydgy_status = ?",ydgyStatus)
	}
	var count int64
	_,err :=builder.Load(&count)

	return count,err
}
//配置登入界面
func (self *GetOnKey) GetOnKey() (*GetOnKey,error) {
	var GetOnKey *GetOnKey
	
	builder :=db.NewSession().Select("status").From("flags")
	
	builder = builder.Where("flag = ?","login_type")
	builder = builder.Where("type = ?","ACCOUNT")
	
	_,err :=builder.LoadStructs(&GetOnKey)	

	return GetOnKey,err
}
//冻结金额增加 freeze_money
func (self *Account) AccountAddFreezeMoney(openId string,money int64) error {
	_,err :=db.NewSession().UpdateBySql("update account set freeze_money=freeze_money+? where open_id=? limit 1",money,openId).Exec()
	return err
}
//冻结金额减少 freeze_money
func (self *Account) AccountMinusFreezeMoney(openId string,money int64) error {
	_,err:=db.NewSession().UpdateBySql("update account set freeze_money=? where open_id=? limit 1",money,openId).Exec()
	return err	
}
func (self *Account) AccountMinusFreezeMoneyTx(openId string,money int64,tx *dbr.Tx) error {
	_,err:=tx.UpdateBySql("update account set freeze_money=? where open_id=? limit 1",money,openId).Exec()
	return err	
}
func (self *Account) AccountWithOpenIdTx(openId string,appId string,tx *dbr.Tx) (*Account,error)  {
	var account *Account
	_,err :=tx.Select("*").From("account").Where("open_id=?",openId).Where("app_id=?",appId).LoadStructs(&account)
	return account,err
}
//获取全部用户
func (self *Account) Accounts(appId string) ([]*Account,error)  {
	var list []*Account
	_,err :=db.NewSession().Select("*").From("account").Where("app_id=?",appId).LoadStructs(&list)
	return list,err
}
func (self *Account) UpdateGetui(openId string,json string) error {
	_,err:=db.NewSession().UpdateBySql("update account set getui=? where open_id=? limit 1",json,openId).Exec()
	return err
}












