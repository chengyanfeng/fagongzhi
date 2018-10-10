package controllers

import (
	"github.com/astaxie/beego"
	"fagongzhi/models"
	"fagongzhi/redis"
	"fagongzhi/util"
	"fmt"
	"time"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {

}

func (c *MainController) GetMessage() {
	returnMessege := models.UpMessageList{}
	phoneNumber := c.GetString("phoneNumber")
	if (util.ToInt(phoneNumber) == -1) {

		lng := c.GetString("lng")
		lat := c.GetString("lat")
		curWord := c.GetString("curWord")
		liebie := c.GetString("liebie")
		fmt.Print(lng)
		fmt.Print(lat)

		//获取地理位置集合
		n := redis.ReGetRediusLoction(lat,lng , "5000")
		nlist := n.([]interface{})
		mdlist := []string{}
		messagelist := []models.UpMessage{}
		maplist := map[string]string{}
		for _, v := range nlist {
			a := v.([]interface{})[0]
			b := v.([]interface{})[1]
			c := v.([]interface{})[2]
			c0 := c.([]interface{})[0]
			c1 := c.([]interface{})[1]

			fmt.Println(util.ToString(a))
			fmt.Println(util.ToString(b))
			fmt.Println(util.ToString(c0))
			fmt.Println(util.ToString(c1))
			mdlist = append(mdlist, util.ToString(a))
			maplist[util.ToString(a)] = util.ToString(b)

		}

		if lng == "" && lat == "" {
			models.DB.Order("time desc").Find(&messagelist).Limit(50)

		} else {

			if len(curWord) > 0 {
				if util.ToInt(liebie) >= 0 {
					models.DB.Where("md5_id in (?) AND key_word= ? AND liebie = ?", mdlist, curWord, liebie).Or("md5_id in (?) AND message like ? AND liebie = ?", mdlist, "%"+curWord+"%", liebie).Order("time desc").Find(&messagelist)

				} else {
					models.DB.Where("md5_id in (?) AND key_word= ?", mdlist, curWord).Or("md5_id in (?) AND message like ?", mdlist, "%"+curWord+"%").Order("time desc").Find(&messagelist)

				}
			} else {
				if util.ToInt(liebie) >= 0 {
					models.DB.Where("md5_id in (?) AND liebie = ?", mdlist, liebie).Order("time desc").Find(&messagelist)
				} else {
					models.DB.Where("md5_id in (?)", mdlist).Order("time desc").Find(&messagelist)

				}

			}
		}
		for k, _ := range messagelist {
			messagelist[k].Distance = maplist[messagelist[k].Md5Id]
		}
		fmt.Print(nlist)
		returnMessege.UserList = &messagelist
		returnMessege.TotalPage = 10
		c.Data["json"] = returnMessege
		c.ServeJSON()
	} else {
		messagelist := []models.UpMessage{}
		models.DB.Where("phone_number = ?", phoneNumber).Order("time desc").Find(&messagelist)
		returnMessege.UserList = &messagelist
		returnMessege.TotalPage = 1
		c.Data["json"] = returnMessege
		c.ServeJSON()
	}
}

func (c *MainController) Index() {
	c.TplName = "index.html"
}

func (c *MainController) UpMessage() {
	lng:="40.00699"
	lat:="116.48349"
	UpMessage := models.UpMessage{}
	leibie := c.GetString("leibie")
	if leibie == "0" {
		local := c.GetString("local")
		destion := c.GetString("destion")
		UpMessage.Local = local
		UpMessage.Destion = destion
	} else {
		keyword := c.GetString("keyword")
		UpMessage.KeyWord = keyword
	}
	lng = c.GetString("lng")
	lat = c.GetString("lat")
	province := c.GetString("province")
	city := c.GetString("city")
	street := c.GetString("street")
	address := c.GetString("address")
	prize := c.GetString("prize")

	UpMessage.Province = province
	UpMessage.City = city
	UpMessage.Street = street
	UpMessage.Address = address
	UpMessage.Lng = util.ToFloat(lng)
	UpMessage.Lat = util.ToFloat(lat)
	message := c.GetString("message")
	UpMessage.Prize = prize
	UpMessage.Message = message
	UpMessage.Liebie = leibie
	ti := util.GetCurTime()
	UpMessage.Time = ti
	md5Id := util.Md5(util.ToString(ti) + message)
	UpMessage.Md5Id = md5Id
	redis.ReSetLoction(md5Id, lng, lat)
	//存储到mysql

		models.DB.Create(&UpMessage)


	c.Data["json"] = message
	c.ServeJSON()
}

func (c *MainController) SendCode() {
	i := 0
	code := models.Code{}
	returnMessage:=models.LoginMessage{}
	phoneNumber := c.GetString("phoneNumber")
	if phoneNumber==""{
		returnMessage.Message="号码为空"
		c.Data["json"] = returnMessage
		c.ServeJSON()
		return
	}

	//先去查询今天发送几次验证码了，如果超过三次，则暂停发送
	models.DB.Where("phone_number = ? AND date = ? ", phoneNumber, util.GetCurDayTime()).Find(&code).Count(&i)
	returnMessage.Count=i
	if i <2 {
		//调取验证码的接口
		sendcode:="1234"
		//把获取的验证码放入到redis 和mysql 中
		//存储到mysql中
		code.Code = sendcode
		code.PhoneNumber = phoneNumber
		code.Date = util.GetCurDayTime()
		models.DB.Create(&code)
		//放到redis中,缓存时间为102秒
		redis.ReAdd(phoneNumber, "1234", 120)
		}


	c.Data["json"] = returnMessage
	c.ServeJSON()
	return
}

func (c *MainController) Login() {
	phoneNumber := c.GetString("phoneNumber")
	password := c.GetString("password")
	user := models.UserPhone{}
	returnMessage:=models.LoginMessage{}
	if phoneNumber==""||password==""{
		returnMessage.Message="信息为空"
		c.Data["json"] = returnMessage
		c.ServeJSON()
		return
	}
	//获取账号密码
	models.DB.Where("phone_number = ?", phoneNumber, util.GetCurDayTime()).Find(&user)
	if user.PhoneNumber==password{
		md5:=phoneNumber+util.ToString(util.GetYesDayTime())
		token:=util.Md5(md5)
		returnMessage.Message="ok"
		returnMessage.Token=token
		//像redis里面放信息
		redis.ReAdd("token",phoneNumber,600000)
	}else {
		returnMessage.Message="err"
		returnMessage.Token=""
	}
	c.Data["json"] = returnMessage
	c.ServeJSON()
	return
}

func (c *MainController) Register() {
	phoneNumber := c.GetString("phoneNumber")
	password := c.GetString("password")
	code := c.GetString("code")
	user := models.UserPhone{}
	returnMessage:=models.LoginMessage{}
	if phoneNumber==""||password==""||code==""{
		returnMessage.Message="信息为空"
		c.Data["json"] = returnMessage
		c.ServeJSON()
		return
	}
	//获取redis 的缓存
	getcode,_:=redis.ReGet(phoneNumber)
	if getcode==code{
		//验证码一致，去数据库里存储user
		user.PhoneNumber=phoneNumber
		user.Time=util.ToInt(time.Now().Unix())
		user.PassWord=password
		models.DB.Create(&user)
		returnMessage.Message="ok"
	}else {
		returnMessage.Message="code is err "
	}
	c.Data["json"] = returnMessage
	c.ServeJSON()
	return
}


func (c *MainController) TestGetPersion(){
	Id:=c.GetString("lastId")
	lastId:=0
	if Id!=""{
		lastId=util.ToInt(Id)
	}

	count:=10
	returnMessege := models.UpMessageList{}
	phoneNumber := c.GetString("phoneNumber")
	if (util.ToInt(phoneNumber) == -1) {

		lng := c.GetString("lng")

		lat := c.GetString("lat")
		curWord := c.GetString("curWord")
		liebie := c.GetString("liebie")
		fmt.Print(lng)
		fmt.Print(lat)
		messagelist := []models.UpMessage{}
		messageCount:=models.UpMessage{}


		if lng == "" && lat == "" {
			if lastId==0{
				models.DB.Order("time desc").Limit(10).Find(&messagelist)
				models.DB.Model(&messageCount).Count(&count)
			}else{
				models.DB.Where("id < ? ",lastId).Order("time desc").Limit(10).Find(&messagelist)

				models.DB.Model(&messageCount).Where("id < ? ",lastId).Count(&count)

			}



		} else {
			maxlat,minlat,minlog,maxlog:=util.GetMinMax(util.ToFloat(lat),util.ToFloat(lng),1000)

			if len(curWord) > 0 {
				if util.ToInt(liebie) >= 0 {
					if lastId==0{
						models.DB.Where("lng > ? AND lng < ? AND lat > ? AND lat < ? AND key_word= ? AND liebie = ?", minlog,maxlog,minlat,maxlat, curWord, liebie).Or("lng > ? AND lng < ? AND lat > ? AND lat < ? AND  message like ? AND liebie = ?", minlog,maxlog,minlat,maxlat, "%"+curWord+"%", liebie).Order("time desc").Limit(10).Find(&messagelist)
						models.DB.Model(&messageCount).Where("lng > ? AND lng < ? AND lat > ? AND lat < ? AND key_word= ? AND liebie = ?", minlog,maxlog,minlat,maxlat, curWord, liebie).Or("lng > ? AND lng < ? AND lat > ? AND lat < ? AND  message like ? AND liebie = ?", minlog,maxlog,minlat,maxlat, "%"+curWord+"%", liebie).Count(&count)

					}else {
						models.DB.Where("lng > ? AND lng < ? AND lat > ? AND lat < ? AND key_word= ? AND liebie = ? AND id < ?", minlog,maxlog,minlat,maxlat, curWord, liebie,lastId).Or("lng > ? AND lng < ? AND lat > ? AND lat < ? AND  message like ? AND liebie = ? AND id < ?", minlog,maxlog,minlat,maxlat, "%"+curWord+"%", liebie,lastId).Order("time desc").Limit(10).Find(&messagelist)

						models.DB.Model(&messageCount).Where("lng > ? AND lng < ? AND lat > ? AND lat < ? AND key_word= ? AND liebie = ? AND id < ?", minlog,maxlog,minlat,maxlat, curWord, liebie,lastId).Or("lng > ? AND lng < ? AND lat > ? AND lat < ? AND  message like ? AND liebie = ? AND id < ?", minlog,maxlog,minlat,maxlat, "%"+curWord+"%", liebie,lastId).Count(&count)
						}
					} else {
					if lastId==0{
						models.DB.Where("lng > ? AND lng < ? AND lat > ? AND lat < ? AND key_word= ?", minlog,maxlog,minlat,maxlat, curWord).Or("lng > ? AND lng < ? AND lat > ? AND lat < ? AND message like ?", minlog,maxlog,minlat,maxlat, "%"+curWord+"%").Order("time desc").Limit(10).Find(&messagelist)
						models.DB.Model(&messageCount).Where("lng > ? AND lng < ? AND lat > ? AND lat < ? AND key_word= ?", minlog,maxlog,minlat,maxlat, curWord).Or("lng > ? AND lng < ? AND lat > ? AND lat < ? AND message like ?", minlog,maxlog,minlat,maxlat, "%"+curWord+"%").Count(&count)

					}else {
						models.DB.Model(&messageCount).Where("lng > ? AND lng < ? AND lat > ? AND lat < ? AND key_word= ? AND id < ?", minlog,maxlog,minlat,maxlat, curWord,lastId).Or("lng > ? AND lng < ? AND lat > ? AND lat < ? AND message like ? AND id < ?", minlog,maxlog,minlat,maxlat, "%"+curWord+"%",lastId).Count(&count)

						models.DB.Where("lng > ? AND lng < ? AND lat > ? AND lat < ? AND key_word= ? AND id < ?", minlog,maxlog,minlat,maxlat, curWord,lastId).Or("lng > ? AND lng < ? AND lat > ? AND lat < ? AND message like ? AND id < ?", minlog,maxlog,minlat,maxlat, "%"+curWord+"%",lastId).Order("time desc").Limit(10).Find(&messagelist)

					}




				}
			} else {
				if util.ToInt(liebie) >= 0 {
					if lastId==0{
						models.DB.Where("lng > ? AND lng < ? AND lat > ? AND lat < ? AND liebie = ?", minlog,maxlog,minlat,maxlat, liebie).Order("time desc").Limit(10).Find(&messagelist)
						models.DB.Model(&messageCount).Where("lng > ? AND lng < ? AND lat > ? AND lat < ? AND liebie = ?", minlog,maxlog,minlat,maxlat, liebie).Count(&count)

					} else {

						models.DB.Where("lng > ? AND lng < ? AND lat > ? AND lat < ? AND liebie = ? AND id < ?", minlog,maxlog,minlat,maxlat, liebie,lastId).Order("time desc").Limit(10).Find(&messagelist)
						models.DB.Model(&messageCount).Where("lng > ? AND lng < ? AND lat > ? AND lat < ? AND liebie = ? AND id < ?", minlog,maxlog,minlat,maxlat, liebie,lastId).Count(&count)

					}

				} else {
					if lastId==0{
						models.DB.Where("lng > ? AND lng < ? AND lat > ? AND lat < ?", minlog,maxlog,minlat,maxlat).Order("time desc").Limit(10).Find(&messagelist)
						models.DB.Model(&messageCount).Where("lng > ? AND lng < ? AND lat > ? AND lat < ?", minlog,maxlog,minlat,maxlat).Count(&count)

					}else {
						models.DB.Where("lng > ? AND lng < ? AND lat > ? AND lat < ? AND id < ?", minlog,maxlog,minlat,maxlat,lastId).Order("time desc").Limit(10).Find(&messagelist)
						models.DB.Model(&messageCount).Where("lng > ? AND lng < ? AND lat > ? AND lat < ? AND id < ?", minlog,maxlog,minlat,maxlat,lastId).Count(&count)

					}


				}

			}
		}
		if lng!=""||lat!=""{
			for k, _ := range messagelist {

				distances:=util.ToString(util.GetDistance(util.ToFloat(lat),util.ToFloat(lng),messagelist[k].Lat,messagelist[k].Lng))
				//获取两点之间的距离
				if len(distances)>5{
					messagelist[k].Distance =distances[0:5]
				}else {

					messagelist[k].Distance =distances
				}


			}
		}


		returnMessege.UserList = &messagelist
		returnMessege.TotalPage = util.GetPage(count,10)
		c.Data["json"] = returnMessege
		c.ServeJSON()
	} else {
		messagelist := []models.UpMessage{}
		models.DB.Where("phone_number = ?", phoneNumber).Order("time desc").Find(&messagelist)
		returnMessege.UserList = &messagelist
		returnMessege.TotalPage =util.GetPage(count,10)
		c.Data["json"] = returnMessege
		c.ServeJSON()
	}



}