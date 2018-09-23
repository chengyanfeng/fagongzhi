package models


//用户表
type User struct {
	PdImg           string `json:"pdImg"`
	PdName   string `json:"pdName"`
	PdPrice    string `json:"pdPrice"`
	PdSold   string    `json:"pdSold"`
}
//用户表
type Test struct {
	UserList *[]User
	TotalPage string `json:"totalPage"`
}
