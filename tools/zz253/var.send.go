package zz253

import (
	"errors"
	"github.com/bububa/ljson"
	"strings"
)

type VarSendRequest struct {
	BaseRequest

	Sign     string     // 短信签名，不含【】
	Tpl      string     // 短信模版，不包含签名, 长度不能超过536个字符
	Params   [][]string // 手机号码和变量参数，多组参数使用英文分号;区分
	SendTime string     //Format: YYYYMMDDHH24MI
	Report   bool       // 是否需要状态报告（默认false），选填
	Extend   int        // 下发短信号码扩展码，纯数字，建议1-3位，选填
	Uid      string     // 场景名（英文或者拼音）-批次编号" //自助通系统内使用UID判断短信使用的场景类型，可重复使用，可自定义场景名称，示例如 VerificationCode（选填参数）
}

type VarSendResponse struct {
	Code     string `json:"code" codec:"code"`         // 状态码（详细参考提交响应状态码）
	ErrorMsg string `json:"errorMsg" codec:"errorMsg"` // 状态码说明（成功返回空）

	MsgId      string `json:"msgId" codec:"msgId"`           // 消息id
	Time       string `json:"time" codec:"time"`             // 响应时间
	FailNum    string `json:"failNum" codec:"failNum"`       // 失败条数
	SuccessNum string `json:"successNum" codec:"successNum"` // 成功条数
}

func VarSend(apiReq *VarSendRequest) (*VarSendResponse, error) {
	if len(apiReq.Params) == 0 {
		return nil, errors.New("no params")
	}
	apiReq.Sign = strings.TrimSpace(apiReq.Sign)
	if apiReq.Sign == "" {
		return nil, errors.New("no sms sign.")
	}
	apiReq.Tpl = strings.TrimSpace(apiReq.Tpl)
	if apiReq.Tpl == "" {
		return nil, errors.New("no sms tpl.")
	}

	client := NewClient(apiReq.Account, apiReq.Password, apiReq.IsTest)
	req := NewRequest(URL_VAR_SEND)

	params := []string{}
	for _, v := range apiReq.Params {
		if len(v) > 0 {
			params = append(params, strings.Join(v, ","))
		}
	}
	if len(params) == 0 {
		return nil, errors.New("no params")
	}
	req.Params["params"] = strings.Join(params, ";")
	req.Params["msg"] = apiReq.Tpl
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
	j := VarSendResponse{}
	err = ljson.Unmarshal([]byte(response), &j)
	if err != nil {
		err = &Error{Code: "-1", Msg: err.Error()}
		logger.Error(err)
		return nil, err
	}
	if j.Code != "0" {
		err = &Error{Code: j.Code, Msg: j.ErrorMsg}
		logger.Error(err)
		return nil, err
	}

	return &j, nil
}
