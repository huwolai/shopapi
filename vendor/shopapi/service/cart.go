package service

import (
	"shopapi/dao"
	"gitlab.qiyunxin.com/tangtao/utils/db"
	"errors"
)

func CartList(openId string) ([]dao.Cart,error) {
	return dao.NewCart().CartList(openId)
}
func CartAddToList(cart dao.Cart) error {
	cartDao := dao.NewCart()
	//====================
	session := db.NewSession()
	tx,_ :=session.Begin()
	defer func() {
		if err :=recover();err!=nil{
			tx.Rollback()
			panic(err)				
		}
	}()
	//====================
	count,err := cartDao.CartExistInList(cart,tx)
	if err!=nil{
		tx.Rollback()
		return err
	}
	//====================
	if count>0{
		err=cartDao.CartAddNumToList(cart,tx)
	}else{
		err=cartDao.CartAddToList(cart,tx)
	}
	//====================
	if err!=nil{
		tx.Rollback()
		return err
	}	
	//====================	
	err = tx.Commit()
	if err!=nil{
		tx.Rollback()
		return err
	}
	return nil
}
func CartMinusFromList(cart dao.Cart) error {
	cartDao := dao.NewCart()
	//====================
	session := db.NewSession()
	tx,_ :=session.Begin()
	defer func() {
		if err :=recover();err!=nil{
			tx.Rollback()
			panic(err)				
		}
	}()
	//====================
	count,err := cartDao.CartNumInLIst(cart,tx)
	if err!=nil{
		tx.Rollback()
		return err
	}
	//====================
	if count>=cart.Num{
		err=cartDao.CartMinusFromList(cart,tx)
	}else{
		tx.Rollback()
		return errors.New("超出购物车商品数量")
	}
	//====================
	if err!=nil{
		tx.Rollback()
		return err
	}	
	//====================	
	err = tx.Commit()
	if err!=nil{
		tx.Rollback()
		return err
	}
	return nil
}
func CartDelFromList(openId string,id uint64) error {
	return dao.NewCart().CartDelFromList(openId,id)
}
func CartUpdateList(cart dao.Cart) error {
	cartDao := dao.NewCart()
	//====================
	session := db.NewSession()
	tx,_ :=session.Begin()
	defer func() {
		if err :=recover();err!=nil{
			tx.Rollback()
			panic(err)				
		}
	}()
	//====================
	count,err := cartDao.CartExistInList(cart,tx)
	if err!=nil{
		tx.Rollback()
		return err
	}
	//====================
	if count>0{
		err=cartDao.CartUpdateList(cart,tx)
	}else{
		tx.Rollback()
		return errors.New("购物车无此商品")
	}
	//====================
	if err!=nil{
		tx.Rollback()
		return err
	}	
	//====================	
	err = tx.Commit()
	if err!=nil{
		tx.Rollback()
		return err
	}
	return nil
}