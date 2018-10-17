package models

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

type Gorm struct {
}
func (UpAdvice) TableName() string {
	return "upadvice"
}
func (UserPhone) TableName() string {
	return "userphone"
}

func (Code) TableName() string {
	return "code"
}
func (UpMessage) TableName() string {
	return "upmessage"
}
func (UpMessageLogin) TableName() string {
	return "upmessagelogin"
}
//数据库初始化
func init() {
	var err error
	conn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=True&loc=Local",
		Mysqlconn["username"],
		Mysqlconn["password"],
		Mysqlconn["host"],
		Mysqlconn["port"],
		Mysqlconn["name"],
	)
	fmt.Print("mysqlconn:")
	fmt.Println(Mysqlconn)
	DB, err = gorm.Open("mysql", conn)
	if err != nil {
		panic(err.Error())
	}
	if DB.HasTable("upmessage") {
		//自动添加模式
		DB.AutoMigrate(&UpMessage{})
		fmt.Println("数据表已经存在")
	} else {
		DB.CreateTable(&UpMessage{})
	}
	if DB.HasTable("code") {
		//自动添加模式
		DB.AutoMigrate(&Code{})
		fmt.Println("code表已经存在")
	} else {
		DB.CreateTable(&Code{})
	}
	if DB.HasTable("userphone") {
		//自动添加模式
		DB.AutoMigrate(&UserPhone{})
		fmt.Println("数据表已经存在")
	} else {
		DB.CreateTable(&UserPhone{})
	}

	if DB.HasTable("upmessagelogin") {
		//自动添加模式
		DB.AutoMigrate(&UpMessageLogin{})
		fmt.Println("数据表已经存在")
	} else {
		DB.CreateTable(&UpMessageLogin{})
	}
	if DB.HasTable("upadvice") {
		//自动添加模式
		DB.AutoMigrate(&UpAdvice{})
		fmt.Println("数据表已经存在")
	} else {
		DB.CreateTable(&UpAdvice{})
	}
}

