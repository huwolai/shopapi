package dao

import (
	"gitlab.qiyunxin.com/tangtao/utils/db"
	"github.com/gocraft/dbr"
	//"errors"
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
func CashoutRecord(appId string,pageIndex uint64,pageSize uint64) ([]*Cashout,error) {
	var cashout []*Cashout
	_,err :=db.NewSession().SelectBySql("select c.*,a.mobile from account_cash_out c left join account as a on c.open_id=a.open_id where c.app_id = ? order by id desc limit ?,?",appId,(pageIndex-1)*pageSize,pageSize).LoadStructs(&cashout)

	return cashout,err
}
func CashoutRecordCount(appId string) int64 {
	count, _ := db.NewSession().Select("count(id)").From("account_cash_out").Where("app_id = ?", appId).ReturnInt64()
	return count
}




























