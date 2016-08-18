package service

import (
	"shopapi/dao"
	"errors"
)


type AddressDto struct  {
	Id int64 `json:"id"`
	OpenId string `json:"open_id"`
	Longitude float64 `json:"longitude"`
	Latitude float64 `json:"latitude"`
	Address string `json:"address"`
	Weight int `json:"weight"`
	Json string `json:"json"`
	AppId string `json:"app_id"`

}

func AddressWithId(id int64) (*dao.Address,error) {
	address := dao.NewAddress()
	address,err := address.WithId(id)
	return address,err
}

func AddressWithRecom(openId string,appId string) (*dao.Address,error) {

	address :=dao.NewAddress()
	address,err := address.AddressWithRecom(openId,appId)
	return address,err
}

func AddressWithOpenId(openId,appId string) ([]*dao.Address,error)  {
	address :=dao.NewAddress()
	addressList,err := address.AddressWithOpenId(openId,appId)
	return addressList,err
}

func AddressAdd(dto *AddressDto) (*AddressDto,error)  {

	address :=AddressDtoToModel(dto)
	address.Weight=0
	aid,err :=address.Insert()
	address.Id = aid

	return AddressToDto(address),err
}

func AddressUpdate(dto *AddressDto) (*AddressDto,error) {

	address := dao.NewAddress()
	address,err := address.WithId(dto.Id)
	if err!=nil {
		return nil,err
	}

	if address==nil{
		return nil,errors.New("地址不存在!")
	}

	fillAddress(address,dto)

	return dto,nil
}

func AddressDelete(id int64) error  {

	address :=dao.NewAddress()
	address.Id = id
	err :=address.Delete()

	return err
}

func fillAddress(model *dao.Address,dto *AddressDto)  {
	model.Address = dto.Address
	model.Longitude = dto.Longitude
	model.Latitude = dto.Latitude
	model.AppId = dto.AppId
	model.Json = dto.Json
	model.OpenId = dto.OpenId

}


func AddressToDto(model *dao.Address) *AddressDto  {
	dto := &AddressDto{}
	dto.Json = model.Json
	dto.Address = model.Address
	dto.Id = model.Id
	dto.Latitude =model.Latitude
	dto.Longitude = model.Longitude
	dto.OpenId = model.OpenId
	dto.Weight = model.Weight
	dto.Id = model.Id
	dto.AppId  = model.AppId

	return dto
}

func AddressDtoToModel(dto *AddressDto) *dao.Address {
	model :=&dao.Address{}
	model.Address = dto.Address
	model.Json = dto.Json
	model.OpenId = dto.OpenId
	model.AppId = dto.AppId
	model.Latitude = dto.Latitude
	model.Longitude = dto.Longitude
	model.Weight = dto.Weight
	model.Id = dto.Id
	return model
}