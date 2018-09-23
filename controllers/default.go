package controllers

import (
	"github.com/astaxie/beego"
	"fagongzhi/models"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	Test:=models.Test{}
	userlist:=make([]models.User,0)
	for i:=0;i<10;i++{
		user:=models.User{}
		user.PdImg="../res/img/pd2.jpg"
		user.PdName="wokao"
		user.PdPrice="1231"
		user.PdSold="12312"
		userlist=	append(userlist, user)
	}
	Test.UserList=&userlist
	Test.TotalPage="10"
	c.Data["json"] = Test
	c.ServeJSON()
}


func (c *MainController) Index() {
	c.TplName="index.html"
}