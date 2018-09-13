package zz253

import (
	"errors"
	"github.com/bububa/ljson"
	"strings"
)

type SendRequest struct {
	BaseRequest

	Phones   []string `json:"phones" codec:"phones"`       // 手机号码。多个手机号码使用英文逗号分隔
	Content  string   `json:"content" codec:"content"`     // 短信内容。长度不能超过536个字符
	SendTime string   `json:"send_time" codec:"send_time"` // 定时发送短信时间。格式为yyyyMMddHHmm，值小于或等于当前时间则立即发送，默认立即发送，选填
	Report   bool     `json:"report" codec:"report"`       // 是否需要状态报告（默认false），选填
	Extend   int      `json:"extend" codec:"extend"`       // 下发短信号码扩展码，纯数字，建议1-3位，选填
	Uid      string   `json:"uid" codec:"uid"`             // 场景名（英文或者拼音）-批次编号" //自助通系统内使用UID判断短信使用的场景类型，可重复使用，可自定义场景名称，示例如 VerificationCode（选填参数）
}

type SendResponse struct {
	Code     string `json:"code" codec:"code"`         // 状态码（详细参考提交响应状态码）
	ErrorMsg string `json:"errorMsg" codec:"errorMsg"` // 状态码说明（成功返回空）

	MsgId string `json:"msgId" codec:"msgId"` // 消息id
	Time  string `json:"time" codec:"time"`   // 响应时间
}

func Send(apiReq *SendRequest) (*SendResponse, error) {
	if len(apiReq.Phones) == 0 {
		return nil, errors.New("no phones.")
	}
	if len(apiReq.Content) == 0 {
		return nil, errors.New("no sms content.")
	}

	client := NewClient(apiReq.Account, apiReq.Password, apiReq.IsTest)
	req := NewRequest(URL_SEND)
	req.Params["phone"] = strings.Join(apiReq.Phones, ",")
	req.Params["msg"] = apiReq.Content
	if apiReq.SendTime != "" {
		req.Params["sendtime"] = apiReq.SendTime
	}
	if apiReq.Extend > 0 {
		req.Params["extend"] = apiReq.Extend
	}
	if apiReq.Report {
		req.Params["report"] = "true"
	}
	if apiReq.Uid != "" {
		req.Params["uid"] = apiReq.Uid
	}

	response, err := client.Execute(req)
	if err != nil {
		return nil, err
	}
	j := SendResponse{}
	err = ljson.Unmarshal([]byte(response), &j)
	if err != nil {
		err = &Error{Code: "-1", Msg: err.Error()}
		logger.Error(err)
		return nil, err
	}
	if apiReq.IsDebug {
		logger.Infof("response: %s", Json(j))
	}
	if j.Code != "0" {
		err = &Error{Code: j.Code, Msg: j.ErrorMsg}
		logger.Error(err)
		return nil, err
	}

	return &j, nil
}
