package zz253

import (
	"fmt"
	"github.com/bububa/ljson"
)

type BalanceRequest struct {
	BaseRequest
}

type BalanceResponse struct {
	Code     int    `json:"code" codec:"code"`         // 状态码（详细参考提交响应状态码）
	ErrorMsg string `json:"errorMsg" codec:"errorMsg"` // 状态码说明（成功返回空）

	Balance string `json:"balance" codec:"balance"` // 剩余可用余额条数
	Time    string `json:"time" codec:"time"`       // 响应时间
}

func Balance(apiReq *BalanceRequest) (*BalanceResponse, error) {
	client := NewClient(apiReq.Account, apiReq.Password, apiReq.IsTest)
	req := NewRequest(URL_BALANCE)

	response, err := client.Execute(req)
	if err != nil {
		return nil, err
	}
	j := BalanceResponse{}
	err = ljson.Unmarshal([]byte(response), &j)
	if err != nil {
		err = &Error{Code: "-1", Msg: err.Error()}
		logger.Error(err)
		return nil, err
	}
	if j.Code != 0 {
		err = &Error{Code: fmt.Sprintf("%v", j.Code), Msg: j.ErrorMsg}
		logger.Error(err)
		return nil, err
	}

	return &j, nil
}
