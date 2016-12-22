package dao

import (
	"gitlab.qiyunxin.com/tangtao/utils/db"
	"github.com/gocraft/dbr"
	//"errors"
	"fmt"
	//"github.com/gocraft/dbr"
)

type Cashout struct  {
	Id 			string `json:"id"`
	AppId 		string `json:"app_id"`
	OpenId 		string `json:"open_id"`
	Amount 		int64  `json:"amount"`
	Title	 	string `json:"title"`
	Remark	 	string `json:"remark"`
	Status	 	int64  `json:"status"`
	Mobile	 	string `json:"mobile"`
	Cashoutcode	string `json:"cashoutcode"`
	Name		string `json:"name"`
}
func MakeCashout(appId string,cashOut interface{},tx *dbr.Tx) (int64,error)  {
	cashout:=cashOut.(Cashout)	
	
	cashout.Status	=0
	cashout.AppId	=appId
	cashout.Amount	=cashout.Amount
	
	result,err :=tx.InsertInto("account_cash_out").Columns("app_id","open_id","amount","title","remark","status").Record(cashout).Exec()

	lastId,err := result.LastInsertId()

	return lastId,err
}
func CashoutRecordTx(cashoutId string,tx *dbr.Tx) (*Cashout,error){
	var cashout *Cashout
	_,err :=db.NewSession().Select("*").From("account_cash_out").Where("id=?",cashoutId).LoadStructs(&cashout)
	return cashout,err
}
func CashoutcodeUpdateTx(cashoutId int64,code string,tx *dbr.Tx) error {
	_,err :=tx.Update("account_cash_out").Set("cashoutcode",code).Where("id=?",cashoutId).Exec()
	return err
}
func CashoutStatusUpdateTx(cashoutId string,cashout map[string]interface{},tx *dbr.Tx) error {
	_,err :=db.NewSession().Update("account_cash_out").Set("status",1).Set("sub_trade_no",cashout["sub_trade_no"].(string)).Where("id=?",cashoutId).Exec()
	return err
}
func CashoutRecord(appId string,pageIndex uint64,pageSize uint64,mobile string,openId string,status string) ([]*Cashout,error) {
	var cashout []*Cashout
	
	sql:=fmt.Sprintf("select c.*,a.mobile,if(m.name is null,'',m.name) as name from account_cash_out c left join account as a on c.open_id=a.open_id left join merchant as m on c.open_id=m.open_id where c.app_id = '%s'",appId)
	
	if mobile!="" {
		sql=sql+fmt.Sprintf(" and a.mobile like '%s%%'",mobile)
	}
	if openId!="" {
		sql=sql+fmt.Sprintf(" and c.open_id = '%s'",openId)
	}
	if status!="" {
		sql=sql+fmt.Sprintf(" and c.status = '%s'",status)
	}
	_,err :=db.NewSession().SelectBySql(sql+" order by id desc limit ?,?",(pageIndex-1)*pageSize,pageSize).LoadStructs(&cashout)	

	return cashout,err
}
func CashoutRecordCount(appId string,mobile string,openId string,status string) int64 {
	sql:=fmt.Sprintf("select count(c.id) from account_cash_out c left join account as a on c.open_id=a.open_id where c.app_id = '%s'",appId)
	
	if mobile!="" {
		sql=sql+fmt.Sprintf(" and a.mobile like '%s%%'",mobile)
	}
	if openId!="" {
		sql=sql+fmt.Sprintf(" and c.open_id = '%s'",openId)
	}
	if status!="" {
		sql=sql+fmt.Sprintf(" and c.status = '%s'",status)
	}
	
	count, _ :=db.NewSession().SelectBySql(sql).ReturnInt64()

	return count
}




























