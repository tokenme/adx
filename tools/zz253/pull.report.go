package zz253

import (
	"fmt"
	"github.com/bububa/ljson"
)

type PullRptRequest struct {
	BaseRequest

	Count int // 拉取个数（最大100，默认20），选填
}

type SendReport struct {
	Uid        string `json:"uid,omitempty" codec:"uid,omitempty"`               // 场景名（英文或者拼音）-批次编号" //自助通系统内使用UID判断短信使用的场景类型，可重复使用，可自定义场景名称，示例如 VerificationCode（选填参数）
	ReportTime string `json:"reportTime,omitempty" codec:"reportTime,omitempty"` // 状态更新时间，格式yyMMddHHmm，其中yy=年份的最后两位（00-99）
	NotifyTime string `json:"notifyTime,omitempty" codec:"notifyTime,omitempty"` // 253平台收到运营商回复状态报告的时间，格式yyyyMMddHHmmss
	Status     string `json:"status,omitempty" codec:"status,omitempty"`         // 状态（详细参考常见常见状态报告状态码）
	StatusDesc string `json:"statusDesc,omitempty" codec:"statusDesc,omitempty"` // 状态说明
	Msgid      string `json:"msgid,omitempty" codec:"msgid,omitempty"`           // 消息id
	Mobile     string `json:"mobile,omitempty" codec:"mobile,omitempty"`         // 接收短信的手机号码
}

type PullRptResponse struct {
	Code   int           `json:"ret" codec:"ret"`       // 请求状态。0成功，其他状态为失败
	Result []*SendReport `json:"result" codec:"result"` // 状态明细结果，没结果则返回空数组
}

func PullReport(apiReq *PullRptRequest) ([]*SendReport, error) {
	client := NewClient(apiReq.Account, apiReq.Password, apiReq.IsTest)
	req := NewRequest(URL_PULL_REPORT)
	if apiReq.Count > 0 {
		req.Params["count"] = fmt.Sprintf("%v", apiReq.Count)
	}

	response, err := client.Execute(req)
	if err != nil {
		return nil, err
	}
	j := PullRptResponse{}
	err = ljson.Unmarshal([]byte(response), &j)
	if err != nil {
		err = &Error{Code: "-1", Msg: err.Error()}
		logger.Error(err)
		return nil, err
	}
	if apiReq.IsDebug {
		logger.Infof("response: %s", Json(j))
	}
	if j.Code != 0 {
		err = &Error{Code: fmt.Sprintf("%v", j.Code), Msg: fmt.Sprintf("ERRCODE: %v", j.Code)}
		logger.Error(err)
		return nil, err
	}

	return j.Result, nil
}
