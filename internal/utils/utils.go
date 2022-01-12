package utils

import (
	"time"
)

type Functions interface {
	AddMock(mock *Mock)
	ClearMock()
	GetNowAsString() string
}

type utilFunctions struct {
	Mock *Mock
}

type Mock struct {
	Date string
}

func New() Functions {
	return &utilFunctions{}
}

func (util *utilFunctions) AddMock(mock *Mock) {
	util.Mock = mock
}

func (util *utilFunctions) ClearMock() {
	util.Mock = nil
}

func (util *utilFunctions) GetNowAsString() string {
	if util.Mock != nil {
		if util.Mock.Date == "" {
			return "02/01/2006 15:04:05"
		} else {
			return util.Mock.Date
		}
	}

	t := time.Now()
	return t.Format("02/01/2006 15:04:05")
}
