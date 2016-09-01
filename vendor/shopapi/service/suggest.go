package service

import "shopapi/dao"

func SuggestAdd(content string,contact string,openId string) error  {

	suggest := dao.NewSuggest()
	suggest.Content = content
	suggest.Contact = contact
	suggest.OpenId = openId
	return suggest.Insert()
}
