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
	userid := c.GetString("userid")
	if (len(userid) == 0) {

		lng := c.GetString("lng")
		lat := c.GetString("lat")
		curWord := c.GetString("curWord")
		liebie := c.GetString("liebie")
		fmt.Print(lng)
		fmt.Print(lat)

		//获取地理位置集合
		n := redis.ReGetRediusLoction(lat, lng, "5000")
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
		models.DB.Where("phone_number = ?", userid).Order("time desc").Find(&messagelist)
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
	lngdefault := "116.48349"
	latdefault := "40.00699"
	UpMessage := models.UpMessage{}
	leibie := c.GetString("leibie")
	local := c.GetString("local")
	destion := c.GetString("destion")
	keyword := c.GetString("keyword")
	if leibie == "0" {

		UpMessage.Local = local
		UpMessage.Destion = destion
	} else {

		UpMessage.KeyWord = keyword
	}
	lng := c.GetString("lng")
	lat := c.GetString("lat")
	if lng==""{
		lng=lngdefault
	}
	if lat==""{
		lat=latdefault
	}
	province := c.GetString("province")
	city := c.GetString("city")
	street := c.GetString("street")
	address := c.GetString("address")
	prize := c.GetString("prize")
	userid := c.GetString("userid")
	message := c.GetString("message")
	ti := util.GetCurTime()
	md5Id := util.Md5(util.ToString(ti) + message)
	UpMessage.Province = province
	UpMessage.City = city
	UpMessage.Street = street
	UpMessage.Address = address
	UpMessage.Lng = util.ToFloat(lng)
	UpMessage.Lat = util.ToFloat(lat)
	UpMessage.Prize = prize
	UpMessage.Message = message
	UpMessage.Liebie = leibie
	UpMessage.Time = ti
	UpMessage.Md5Id = md5Id
	redis.ReSetLoction(md5Id, lng, lat)
	//存储到mysql,这张临时表是不插入电话号码的
	models.DB.Create(&UpMessage)
	if (userid!=""){
		phoneNumber,_:=	redis.ReGet(userid)
		if len(phoneNumber)>0{
			//已经登陆的，插入到不删除的表
			loginmessage:=	models.UpMessageLogin{}
			loginmessage.Province = province
			loginmessage.City = city
			loginmessage.Street = street
			loginmessage.Address = address
			loginmessage.Local=local
			loginmessage.Destion=destion
			loginmessage.KeyWord=keyword
			loginmessage.Lng = util.ToFloat(lng)
			loginmessage.Lat = util.ToFloat(lat)
			loginmessage.Prize = prize
			loginmessage.Message = message
			loginmessage.Liebie = leibie
			loginmessage.Time = ti
			loginmessage.Md5Id = md5Id
			loginmessage.PhoneNumber=phoneNumber
			//插入用不删除的表中
			models.DB.Create(&loginmessage)
		}
	}

	c.Data["json"] = message
	c.ServeJSON()
}

func (c *MainController) SendCode() {

	i := 0
	code := models.Code{}
	returnMessage := models.LoginMessage{}
	phoneNumber := c.GetString("phoneNumber")
	if phoneNumber == "" {
		returnMessage.Message = "号码为空"
		c.Data["json"] = returnMessage
		c.ServeJSON()
		return
	}

	//先去查询今天发送几次验证码了，如果超过三次，则暂停发送
	models.DB.Where("phone_number = ? AND date = ? ", phoneNumber, util.GetCurDayTime()).Find(&code).Count(&i)
	returnMessage.Count = i
	if i < 2 {
		//调取验证码的接口
		sendcode := "1234"
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
	returnMessage := models.LoginMessage{}
	if phoneNumber == "" || password == "" {
		returnMessage.Message = "信息为空"
		c.Data["json"] = returnMessage
		c.ServeJSON()
		return
	}
	//获取账号密码
	models.DB.Where("phone_number = ?", phoneNumber).Find(&user)
	if user.PassWord == password {
		md5 := phoneNumber + util.ToString(util.GetYesDayTime())
		token := util.Md5(md5)
		returnMessage.Message = "ok"
		returnMessage.Token = token
		//像redis里面放信息,600000 秒
		redis.ReAdd(token, phoneNumber, 600000)
	} else {
		returnMessage.Message = "err"
		returnMessage.Token = ""
	}
	c.Data["json"] = returnMessage
	c.ServeJSON()
	return
}

func (c *MainController) Register() {
	phoneNumber := c.GetString("phoneNumber")
	password := c.GetString("password")
	code := c.GetString("code")
	userPhone := models.UserPhone{}
	returnMessage := models.LoginMessage{}
	if phoneNumber == "" || password == "" || code == "" {
		returnMessage.Message = "信息为空"
		c.Data["json"] = returnMessage
		c.ServeJSON()
		return
	}
	//查看是否已经注册过
	models.DB.Where("phone_number = ?",phoneNumber).Find(&userPhone)
	if userPhone.Time>0{
		returnMessage.Count=-1
		returnMessage.Message="err"
		c.Data["json"] = returnMessage
		c.ServeJSON()
		return
	}
	//获取redis 的缓存
	getcode, _ := redis.ReGet(phoneNumber)
	if getcode == code {
		//验证码一致，去数据库里存储user
		userPhone.PhoneNumber = phoneNumber
		userPhone.Time = util.ToInt(time.Now().Unix())
		userPhone.PassWord = password
		models.DB.Create(&userPhone)
		returnMessage.Message = "ok"
	} else {
		returnMessage.Message = "code is err "
	}
	c.Data["json"] = returnMessage
	c.ServeJSON()
	return
}

func (c *MainController) TestGetPersion() {
	Id := c.GetString("lastId")
	myselfid:=c.GetString("myselflastid")
	lastId := 0
	myselflastId:=0
	if myselfid!=""{
		myselflastId=util.ToInt(myselfid)
	}

	if Id != "" {
		lastId = util.ToInt(Id)
	}
	//每页显示多少个，默认为30，
	limt := 30
	//总数量默认为30，求分页
	count := 30
	returnMessege := models.UpMessageList{}
	userid := c.GetString("userid")
	//从reids 里获取phonenumber
	phoneNumber,_:=redis.ReGet(userid)
	if (len(phoneNumber) == 0&&userid=="") {

		lng := c.GetString("lng")

		lat := c.GetString("lat")
		curWord := c.GetString("curWord")
		liebie := c.GetString("liebie")
		fmt.Print(lng)
		fmt.Print(lat)
		messagelist := []models.UpMessage{}
		messageCount := models.UpMessage{}

		if lng == "" && lat == "" {
			if lastId == 0 {
				models.DB.Order("time desc").Limit(limt).Find(&messagelist)
				models.DB.Model(&messageCount).Count(&count)
			} else {
				models.DB.Where("id < ? ", lastId).Order("time desc").Limit(limt).Find(&messagelist)
			}

		} else {
			//方圆5000公里。实际是一个方形
			maxlat, minlat, maxlog, minlog := util.GetMinMax(util.ToFloat(lat), util.ToFloat(lng), 5100)

			if len(curWord) > 0 {
				if util.ToInt(liebie) >= 0 {
					if lastId == 0 {
						models.DB.Where("lng > ? AND lng < ? AND lat > ? AND lat < ? AND key_word= ? AND liebie = ?", minlog, maxlog, minlat, maxlat, curWord, liebie).Or("lng > ? AND lng < ? AND lat > ? AND lat < ? AND  message like ? AND liebie = ?", minlog, maxlog, minlat, maxlat, "%"+curWord+"%", liebie).Order("time desc").Limit(limt).Find(&messagelist)
						models.DB.Model(&messageCount).Where("lng > ? AND lng < ? AND lat > ? AND lat < ? AND key_word= ? AND liebie = ?", minlog, maxlog, minlat, maxlat, curWord, liebie).Or("lng > ? AND lng < ? AND lat > ? AND lat < ? AND  message like ? AND liebie = ?", minlog, maxlog, minlat, maxlat, "%"+curWord+"%", liebie).Count(&count)

					} else {
						models.DB.Where("lng > ? AND lng < ? AND lat > ? AND lat < ? AND key_word= ? AND liebie = ? AND id < ?", minlog, maxlog, minlat, maxlat, curWord, liebie, lastId).Or("lng > ? AND lng < ? AND lat > ? AND lat < ? AND  message like ? AND liebie = ? AND id < ?", minlog, maxlog, minlat, maxlat, "%"+curWord+"%", liebie, lastId).Order("time desc").Limit(limt).Find(&messagelist)
					}
				} else {
					if lastId == 0 {
						models.DB.Where("lng > ? AND lng < ? AND lat > ? AND lat < ? AND key_word= ?", minlog, maxlog, minlat, maxlat, curWord).Or("lng > ? AND lng < ? AND lat > ? AND lat < ? AND message like ?", minlog, maxlog, minlat, maxlat, "%"+curWord+"%").Order("time desc").Limit(limt).Find(&messagelist)
						models.DB.Model(&messageCount).Where("lng > ? AND lng < ? AND lat > ? AND lat < ? AND key_word= ?", minlog, maxlog, minlat, maxlat, curWord).Or("lng > ? AND lng < ? AND lat > ? AND lat < ? AND message like ?", minlog, maxlog, minlat, maxlat, "%"+curWord+"%").Count(&count)

					} else {
						models.DB.Model(&messageCount).Where("lng > ? AND lng < ? AND lat > ? AND lat < ? AND key_word= ? AND id < ?", minlog, maxlog, minlat, maxlat, curWord, lastId).Or("lng > ? AND lng < ? AND lat > ? AND lat < ? AND message like ? AND id < ?", minlog, maxlog, minlat, maxlat, "%"+curWord+"%", lastId).Count(&count)

					}
				}
			} else {
				if util.ToInt(liebie) >= 0 {
					if lastId == 0 {
						models.DB.Where("lng > ? AND lng < ? AND lat > ? AND lat < ? AND liebie = ?", minlog, maxlog, minlat, maxlat, liebie).Order("time desc").Limit(limt).Find(&messagelist)
						models.DB.Model(&messageCount).Where("lng > ? AND lng < ? AND lat > ? AND lat < ? AND liebie = ?", minlog, maxlog, minlat, maxlat, liebie).Count(&count)

					} else {
						models.DB.Where("lng > ? AND lng < ? AND lat > ? AND lat < ? AND liebie = ? AND id < ?", minlog, maxlog, minlat, maxlat, liebie, lastId).Order("time desc").Limit(limt).Find(&messagelist)
					}

				} else {
					if lastId == 0 {
						models.DB.Where("lng > ? AND lng < ? AND lat > ? AND lat < ?", minlog, maxlog, minlat, maxlat).Order("time desc").Limit(limt).Find(&messagelist)
						models.DB.Model(&messageCount).Where("lng > ? AND lng < ? AND lat > ? AND lat < ?", minlog, maxlog, minlat, maxlat).Count(&count)
					} else {
						models.DB.Where("lng > ? AND lng < ? AND lat > ? AND lat < ? AND id < ?", minlog, maxlog, minlat, maxlat, lastId).Order("time desc").Limit(limt).Find(&messagelist)

					}

				}

			}
		}
		if lng != "" || lat != "" {
			for k, _ := range messagelist {

				distances := util.ToString(util.GetDistance(util.ToFloat(lat), util.ToFloat(lng), messagelist[k].Lat, messagelist[k].Lng))
				//获取两点之间的距离
				if len(distances) > 5 {
					messagelist[k].Distance = distances[0:5]
				} else {

					messagelist[k].Distance = distances
				}
			}
		}
		returnMessege.UserList = &messagelist
		returnMessege.TotalPage = util.GetPage(count, 30)
		c.Data["json"] = returnMessege
		c.ServeJSON()
	} else {
		messagelist := []models.UpMessageLogin{}
		messagemysl:=models.UpMessageLogin{}
		if myselflastId==0{
			models.DB.Where("phone_number = ?", phoneNumber).Order("time desc").Find(&messagelist).Limit(limt)
			models.DB.Model(&messagemysl).Count(&count)
		}else {
			models.DB.Where("phone_number = ? AND id < ?", phoneNumber,myselflastId).Order("time desc").Find(&messagelist).Limit(limt)

		}
		returnMessege.UserListLogin = &messagelist
		totalPage := util.GetPage(count, limt)
		returnMessege.TotalPage = totalPage
		c.Data["json"] = returnMessege
		c.ServeJSON()
	}

}

func (c *MainController) ResTPassWord(){
	returnMessage:=models.LoginMessage{}
	userPhone:=models.UserPhone{}
	phoneNumber := c.GetString("phoneNumber")
	password := c.GetString("password")
	code := c.GetString("code")
	if len(phoneNumber) >0{
		//查看是否已经注册过
		models.DB.Where("phone_number = ?",phoneNumber).Find(&userPhone)
		//已经注册过
		if userPhone.PhoneNumber==phoneNumber{
			//查看验证码
		redisCode,_:=	redis.ReGet(phoneNumber)
		if redisCode==code{
			//验证码正确修改密码
			models.DB.Model(&userPhone).Update("pass_word", password)
			returnMessage.Count=0
			returnMessage.Message="成功修改密码"
			returnMessage.Token=""
			}else {
				returnMessage.Count=4
				returnMessage.Message="验证码错误"
				returnMessage.Token=""

		}
		}else {
			returnMessage.Count=4
			returnMessage.Message="您还没有注册"
			returnMessage.Token=""

		}
	}else {
		returnMessage.Count=4
		returnMessage.Message="手机号码为空"
		returnMessage.Token=""
	}
	c.Data["json"] = returnMessage
	c.ServeJSON()
	return
}


func (c *MainController) UpAdvice(){
	upMessage:=models.UpAdvice{}
	returnmessage:=models.LoginMessage{}
	upadvice := c.GetString("upadvice")

	if len(upadvice)>0{
		upMessage.Date=util.ToString(time.Now())
		upMessage.Mesaage=upadvice
		models.DB.Create(&upMessage)
		returnmessage.Count=0
		returnmessage.Message="您已经成功上传"
		c.Data["json"] = returnmessage
	}else{
		returnmessage.Count=1
		returnmessage.Message="数据为空"
		c.Data["json"] = returnmessage
	}
	c.ServeJSON()
	return
}