package commute

import (
	"encoding/json"
	"github.com/Unknwon/com"
	"io/ioutil"
	"net/http"
)

type smsInfo struct {
	account   string
	password  string
	serverUrl string
}

var smsConf *smsInfo

type SMSOpt struct {
	Phone string
	Msg   string
}

func SendSMS(phone, msg string) error {
	if smsConf == nil {
		panic("init sms conf first!")
	}

	p := map[string]interface{}{
		"account":  smsConf.account,
		"password": smsConf.password,
		"phone":    phone,
		"msg":      msg,
		"report":   false,
	}

	body, err := json.Marshal(&p)
	if err != nil {
		panic(err)
	}

	client := http.Client{}
	header := http.Header{}
	header.Set("Content-Type", "application/json;charset=UTF-8")
	reader, err := com.HttpPost(&client, smsConf.serverUrl, header, body)
	if err != nil {
		return err
	}
	defer reader.Close()

	respData, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	var res struct {
		Time     string `json:"time"`
		MsgId    string `json:"msgId"`
		Code     int    `json:"code,string"`
		ErrorMsg string `json:"errorMsg"`
	}

	//{"code":"0","msgId":"19041619345725959","time":"20190416193457","errorMsg":""}
	if err := json.Unmarshal(respData, &res); err != nil {
		return err
	}

	return nil
}

func SetupSMSConfig(account, password, serverUrl string) {
	smsConf = &smsInfo{
		account:   account,
		password:  password,
		serverUrl: serverUrl,
	}
}
