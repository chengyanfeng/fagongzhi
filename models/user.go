package models


//用户表
type UserPhone struct {
	ID           int  `gorm:"primary_key"`
	PhoneNumber   string `json:"phone_number"`
	PassWord	  string `json:"password"`
	Time    int `json:"time"`
}
//用户表
type Test struct {
	UserList *[]UserPhone
	TotalPage string `json:"totalPage"`
}

type UpMessageList struct {
	UserList *[]UpMessage
	UserListLogin *[]UpMessageLogin
	TotalPage int `json:"totalPage"`
}


//信息表
type UpMessage struct {
	ID           int `gorm:"primary_key"`
	Md5Id           string `json:"md5_id"`
	Liebie   string `json:"liebie"`
	Lng    float64 `json:"lng"`
	Lat   float64    `json:"lat"`
	Local   string `json:"local"`
	Destion  string `json:"destion"`
	KeyWord  string `json:"keyword"`
	Message string `json:"message"`
	Time    int `json:"time"`
	Province  string `json:"province"`
	City     string `json:"city"`
	Street  string `json:"street"`
	Address string `json:"address"`
	Distance string `json:"distance"`
	Prize string `json:"prize"`


}
//信息表
type UpMessageLogin struct {
	ID           int `gorm:"primary_key"`
	Md5Id           string `json:"md5_id"`
	Liebie   string `json:"liebie"`
	Lng    float64 `json:"lng"`
	Lat   float64    `json:"lat"`
	Local   string `json:"local"`
	Destion  string `json:"destion"`
	KeyWord  string `json:"keyword"`
	Message string `json:"message"`
	Time    int `json:"time"`
	Province  string `json:"province"`
	City     string `json:"city"`
	Street  string `json:"street"`
	Address string `json:"address"`
	Distance string `json:"distance"`
	Prize string `json:"prize"`
	PhoneNumber string `json:"phone_number"`

}
//验证码表
type Code struct {
	ID           int `gorm:"primary_key"`
	PhoneNumber string `json:"phone_number"`
	Code     string `json:"code"`
	Date    int `json:"date"` //当天的零点时间

}

type LoginMessage struct {
	Count   int `gorm:"count"`
	Message string `gorm:"message"`
	Token   string `gorm:"message"`
}


type UpAdvice struct {
	ID           int `gorm:"primary_key"`
	Mesaage string `json:"message"`
	Date    string `json:"date"` //time
}