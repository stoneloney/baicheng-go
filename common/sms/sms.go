package sms

import (
	"fmt"
	"net/http"
	"time"
	"crypto/sha256"
	"errors"
	"encoding/json"
	"bytes"
	"io/ioutil"

	"common/system"
	"common/method"
)

const QCAPPID = "1400164650"
const QCAPPKEY = "4065a2d3470303f9a74117ccee74eb17"

var QCMsg = map[uint]string {
	1022: "业务短信日下发条数超过设定的上限",
	1023: "单个手机号30秒内下发短信条数超过设定的上限",
	1024: "单个手机号1小时内下发短信条数超过设定的上限",
	1025: "单个手机号日下发短信条数超过设定的上限",
	1026: "单个手机号下发相同内容超过设定的上限",
}

type QcSmsTel struct {
	Nationcode  string  `json:"nationcode"`
	Mobile      string  `json:"mobile"`
}

// 请求包
type QcSmsSingleReq struct {
	Tel    QcSmsTel  `json:"tel"`
	Ext    string    `json:"ext"`
	Extend string    `json:"extend"`
	Params []string  `json:"params"`
	Sig    string    `json:"sig"`
	Sign   string    `json:"sign,omitempty"`
	Time   int64     `json:"time"`
	Tplid  int       `json:"tpl_id"`
}

// 结果包
type QcSmsResult struct {
	Result   uint    `json:"result"`
	Errmsg   string  `json:"errmsg"`
	Ext      string  `json:"ext"`
	Sid      string  `json:"sid,omitempty"`
	Fee      uint    `json:"fee,omitempty"`   
}

// 腾讯云短信
func SendQcloudSmsSingle(phone string) (string, error) {
	if (!method.CheckMobile(phone)) {
		return "", errors.New("请填写正确的手机号")
	}

	// 获取短信配置 
	smsConfig := system.GetSmsConfig()

	// 限制验证码长度
	slen := smsConfig.SmsLength
	if slen < 4 {
		slen = 4
	}
	if slen > 6 {
		slen = 6
	}

	// 多少分钟内有效
	expire := smsConfig.Expire

	random := method.GetRandomString(12)
	number := method.GetRandomNumber(slen)
	t := time.Now().Unix()

	// 生成签名
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("appkey=%s&random=%s&time=%d&mobile=%s", QCAPPKEY, random, t, phone)))
	sig := fmt.Sprintf("%x", h.Sum(nil))
	reqUrl := "https://yun.tim.qq.com/v5/tlssmssvr/sendsms?sdkappid="+QCAPPID+"&random="+random

	// 构建数据格式
	var sm = QcSmsSingleReq{
		Tel: QcSmsTel{Nationcode: "86", Mobile:phone},
		Sig: sig,
		Time: t,
		Tplid: smsConfig.Tplid,
		Params: []string{number, string(expire)},
	}

	d, _ := json.Marshal(sm)
	req, err := http.NewRequest("POST", reqUrl, bytes.NewBuffer([]byte(d)))
	req.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client {
		Timeout: 5*time.Second,
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("请求失败")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var res QcSmsResult
	json.Unmarshal([]byte(body), &res)

	if res.Result == 0 {
		return number, nil
	} 
	if _, ok := QCMsg[res.Result]; ok {
		return "", errors.New(QCMsg[res.Result])
	}
	fmt.Println(res.Errmsg)
	return "", errors.New(res.Errmsg)
}
