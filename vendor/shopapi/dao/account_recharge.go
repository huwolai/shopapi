package dao

import (
	"github.com/gocraft/dbr"
	"gitlab.qiyunxin.com/tangtao/utils/db"
	//"gitlab.qiyunxin.com/tangtao/utils/log"
)

type AccountRecharge struct  {
	Id int64
	No string
	AppId string
	OpenId string
	Amount float64
	Status int	
	Flag string
	Json string
	Froms int
	CreateTimeUnix  int64
	Mobile string
	
	Opt string
	Remark string
	
	BaseDModel
}

type AccountRechargeSearch struct {
	No		 	 string
	YdgyId  	 string
	YdgyName  	 string
	Mobile  	 string
}

func NewAccountRecharge() *AccountRecharge {

	return &AccountRecharge{}
}

func (self *AccountRecharge) InsertTx(tx *dbr.Tx) error  {

	_,err :=tx.InsertInto("account_recharge").Columns("no","app_id","open_id","amount","status","flag","json","froms","opt","remark").Record(self).Exec()

	return err
}

func (self *AccountRecharge) Insert() error  {

	_,err :=db.NewSession().InsertInto("account_recharge").Columns("no","app_id","open_id","amount","status","flag","json","froms").Record(self).Exec()

	return err
}


func (self *AccountRecharge) WithNo(no string,appId string) (*AccountRecharge,error)  {
	var model *AccountRecharge
	_,err :=db.NewSession().Select("*").From("account_recharge").Where("no=?",no).Where("app_id=?",appId).LoadStructs(&model)

	return model,err
}

func (self *AccountRecharge) UpdateStatusWithNo(status int,no string,appId string) error {

	_,err :=db.NewSession().Update("account_recharge").Set("status",status).Where("no=?",no).Where("app_id=?",appId).Exec()

	return err
}
//账户充值记录
func (self *AccountRecharge) WithOpenId(openId string,appId string,froms int64) ([]*AccountRecharge,error)  {
	var model []*AccountRecharge
	
	buider :=db.NewSession().Select("*").From("account_recharge").Where("open_id=?",openId).Where("app_id=?",appId)
	
	if froms>=0 {
		buider = buider.Where("froms=?",froms)
	}
	
	_,err :=buider.OrderDir("id",false).Limit(68).LoadStructs(&model)
	return model,err
}
func (self *AccountRecharge) RecordWithUser(appId string,froms int64,pageIndex uint64,pageSize uint64, search AccountRechargeSearch) ([]*AccountRecharge,error)  {
	var model []*AccountRecharge
	
	/* buider :=db.NewSession().Select("account_recharge.*,account.mobile,UNIX_TIMESTAMP(account_recharge.create_time) as create_time_unix").From("account_recharge").LeftJoin("account","account_recharge.open_id=account.open_id").Where("account_recharge.app_id=?",appId).Where("account.mobile is not null")
	
	if froms>=0 {
		buider = buider.Where("account_recharge.froms=?",froms)
	} 
	
	//log.Error( buider.ToSql() )
	_,err :=buider.OrderDir("account_recharge.id",false).LoadStructs(&model)*/
	
	buider :=db.NewSession().Select("*,UNIX_TIMESTAMP(create_time) as create_time_unix").From("account_recharge").Where("app_id=?",appId)
	
	if froms>=0 {
		buider = buider.Where("froms=?",froms)
	}
	
	if search.No!="" {
		buider = buider.Where("no like ?",search.No+"%")
	}
	if search.YdgyId!="" {
		buider = buider.Where("open_id in (select open_id from account where ydgy_id like ?)",search.YdgyId+"%")
	}
	if search.YdgyName!="" {
		buider = buider.Where("open_id in (select open_id from account where ydgy_name like ?)",search.YdgyName+"%")
	}
	if search.Mobile!="" {
		buider = buider.Where("open_id in (select open_id from account where mobile like ?)",search.Mobile+"%")	
	}
	
	_,err :=buider.Limit(pageSize).Offset((pageIndex-1)*pageSize).OrderDir("id",false).LoadStructs(&model)
	
	return model,err
}
func (self *AccountRecharge) RecordWithUserCount(appId string,froms int64,pageIndex uint64,pageSize uint64, search AccountRechargeSearch) (int64,error)  {
	buider :=db.NewSession().Select("count(id)").From("account_recharge").Where("app_id=?",appId)
	
	if froms>=0 {
		buider = buider.Where("froms=?",froms)
	}
	
	if search.No!="" {
		buider = buider.Where("no like ?",search.No+"%")
	}
	if search.YdgyId!="" {
		buider = buider.Where("open_id in (select open_id from account where ydgy_id like ?)",search.YdgyId+"%")
	}
	if search.YdgyName!="" {
		buider = buider.Where("open_id in (select open_id from account where ydgy_name like ?)",search.YdgyName+"%")
	}
	if search.Mobile!="" {
		buider = buider.Where("open_id in (select open_id from account where mobile like ?)",search.Mobile+"%")	
	}
	
	var count int64
	_,err :=buider.LoadStructs(&count)

	return count,err
}














