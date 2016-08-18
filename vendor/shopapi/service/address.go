package service

import "shopapi/dao"


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

func AddressUpdate(dto *AddressDto)  {

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