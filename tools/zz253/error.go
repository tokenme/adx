package zz253

import (
	"fmt"
)

type Error struct {
	Code string
	Msg  string
}

func (e Error) Error() string {
	return fmt.Sprintf("CODE:%v, MSG:%v", e.Code, e.Msg)
}
