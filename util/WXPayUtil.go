package util

import (
	"crypto/md5"
	"crypto/sha256"
	"egg_backend/def"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sort"
	"strings"
	"time"
)

/******************************-----------下面是获取转发的token和ticker与下面的登陆的tokenk不一样-----------*************************/
//获取转发的token
func GetForwardToken() (token string) {
	//获取微信转发token
	response_token, _ := http.Get("https://api.weixin.qq.com/cgi-bin/token?appid=wx53d52d70ccd6439f&secret=dfb513840c45e387cd869af3887e69cb&grant_type=client_credential")
	defer response_token.Body.Close()
	token_body, _ := ioutil.ReadAll(response_token.Body)
	p := *JsonDecode([]byte(string(token_body)))
	token = p["access_token"].(string)
	fmt.Println("这是从转发获取拿的token")
	return
}

//根据token来获取ticker
func GetTicket(token string) string {
	//从token获取微信ticket
	response_ticket, _ := http.Get("https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=" + token + "&type=jsapi")
	defer response_ticket.Body.Close()
	ticket_body, _ := ioutil.ReadAll(response_ticket.Body)
	p := *JsonDecode([]byte(string(ticket_body)))
	ticket := p["ticket"].(string)
	AddCache("ticket", ticket)
	fmt.Println("ticket 是从重新拿的")
	return string(ticket)
}

/******************************-----------下面是获取登陆的token-----------*************************/
//获取登陆的token和openid
func GetTokenAndOpenid(code string) (access_token, openid string) {
	//获取微信登陆的token
	response_token, _ := http.Get("https://api.weixin.qq.com/sns/oauth2/access_token?appid=" + def.WEIXINAPPID + "&secret=" + def.WEIXINKEY + "&code=" + code + "&grant_type=authorization_code")
	//关闭链接
	defer response_token.Body.Close()

	token_body, _ := ioutil.ReadAll(response_token.Body)

	p := *JsonDecode([]byte(string(token_body)))
	if p["errcode"] != nil {
		return "1", ""
	}
	refresh_token := p["refresh_token"].(string)
	//直接通过获取的token获取刷新token
	refresh_token_token, _ := http.Get("https://api.weixin.qq.com/sns/oauth2/refresh_token?appid=" + def.WEIXINAPPID + "&grant_type=refresh_token&refresh_token=" + refresh_token)
	defer refresh_token_token.Body.Close()
	ticket_body, _ := ioutil.ReadAll(refresh_token_token.Body)
	p = *JsonDecode([]byte(string(ticket_body)))
	access_token = p["access_token"].(string)
	openid = p["openid"].(string)

	if checkToken(access_token, openid) {
		return
	} else {
		return "token is error", "openid is error"
	}

}

//验证token和openid是否有效
func checkToken(access_token, openid string) bool {
	checkToken, _ := http.Get("https://api.weixin.qq.com/sns/auth?access_token=" + access_token + "&openid=" + openid)
	defer checkToken.Body.Close()
	checkToken_body, _ := ioutil.ReadAll(checkToken.Body)
	p := *JsonDecode([]byte(string(checkToken_body)))
	errmsg := p["errmsg"].(string)
	if errmsg == "ok" {
		//把token放到缓存里，k=token,v=openid
		AddCache(access_token, openid)
		AddCache("token", access_token)
		return true

	} else {
		return false
	}
}

//获取微信登陆用户信息
func GetUserInfo(code string) (p *map[string]interface{}) {
	access_token, openid := GetTokenAndOpenid(code)
	if access_token == "1" {

		return
	}
	userInfo, _ := http.Get("https://api.weixin.qq.com/sns/userinfo?access_token=" + access_token + "&openid=" + openid + "&lang=zh_CN")
	defer userInfo.Body.Close()
	userInfo_body, _ := ioutil.ReadAll(userInfo.Body)
	p = JsonDecode([]byte(string(userInfo_body)))

	return

}

//微信服务器获取上传的文件和图片
func GetImageFromCould(mediaId, url string) (imagePath string) {
	token := ToString(GetCache("forword_token"))
	retrnbody, _ := http.Get("https://api.weixin.qq.com/cgi-bin/media/get?access_token=" + token + "&media_id=" + mediaId)
	defer retrnbody.Body.Close()

	imageName := retrnbody.Header.Get("Content-Disposition")
	if imageName == "" {
		return "fail"
	}
	imageName = strings.Split(imageName, "=")[1]
	imageName = strings.Replace(imageName, "\"", "", -1)
	token_body, _ := ioutil.ReadAll(retrnbody.Body)
	URL := url + imageName
	flag := WriteFile(URL, token_body)
	if flag {
		return URL
	} else {
		return "保存图片失败"
	}
}

//短信验证码
func SendMessage(mobile string) bool {
	randVcode := GetRandomString()
	tim := time.Now().Unix()
	sendString := "appkey=" + def.MESSAPPKEY + "&random=7226249334&time=" + ToString(tim) + "&mobile=" + mobile
	sign := fmt.Sprintf("%x", sha256.Sum256([]byte(sendString)))
	senda := `{"params":["` + randVcode + `"],"sig":"` + sign + `","tel":{"mobile":"` + mobile + `", "nationcode":"86" } ,"time":` + ToString(tim) + `,"tpl_id":122011}`
	response, _ := http.Post("https://yun.tim.qq.com/v5/tlssmssvr/sendsms?sdkappid=1400091286&random=7226249334", "application/json;charset=utf-8", strings.NewReader(senda))
	defer response.Body.Close()
	token_body, _ := ioutil.ReadAll(response.Body)
	p := *JsonDecode([]byte(string(token_body)))
	result := p["result"].(float64)
	if result == 0 {
		//添加到缓存里
		AddCache(mobile, randVcode)
		ac := GetCache(mobile)
		fmt.Print(ac)
		return true
	} else {
		return false
	}

}

//微信支付
func GetWXpay_id(openid string) (xml string) {
	userMap := &StringMap{}
	(*userMap)["appid"] = def.WEIXINAPPID
	(*userMap)["mch_id"] = def.WEIXINMCH_ID
	(*userMap)["nonce_str"] = GetRandomString()
	(*userMap)["body"] = "1212121"
	(*userMap)["out_trade_no"] = "123456"
	(*userMap)["total_fee"] = "1"
	(*userMap)["spbill_create_ip"] = "123.12.12.123"
	(*userMap)["trade_type"] = "JSAPI"
	(*userMap)["notify_url"] = "http://www.weixin.qq.com/wxpay/pay.php"
	(*userMap)["sign_type"] = "MD5"
	(*userMap)["openid"] = openid

	xml = MapToxml(userMap)
	response, _ := http.Post("https://api.mch.weixin.qq.com/sandbox/pay/unifiedorder", "application/xml;charset=utf-8", strings.NewReader(xml))
	defer response.Body.Close()
	token_body, _ := ioutil.ReadAll(response.Body)
	xml = string(token_body)
	return xml
}

/******************************-----------公共方法----------*************************/

//Map转xml
func MapToxml(userMap *StringMap) string {
	(*userMap)["sign"] = GetSign(userMap)
	buf, _ := xml.Marshal(StringMap(*userMap))
	xml := string(buf)
	xml = strings.Replace(xml, "StringMap", "xml", -1)
	return xml
}

//获取签名
func GetSign(p *StringMap) string {
	sign := ""
	md := md5.New()
	strs := []string{}
	for k := range *p {
		strs = append(strs, k)
	}
	sort.Strings(strs)
	for _, v := range strs {
		sign = sign + v + "=" + (*p)[v] + "&"
	}
	sign = sign + "key=" + def.WEIXINKEY
	fmt.Print(sign)
	md.Write([]byte(sign))
	sign = fmt.Sprintf("%x", md5.Sum([]byte(sign)))
	return strings.ToUpper(sign)

}

//生成随机字符串
func GetRandomString() string {
	bytes := []byte(def.WEIXINRANDSTR)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 4; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
