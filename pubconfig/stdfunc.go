package pubconfig

import (
	"encoding/json"
	"io/ioutil"
)

type JsonCode struct{}

func NewJsonCode() *JsonCode {
	return &JsonCode{}
}

func (jc *JsonCode) Load(fileName string, v interface{}) error {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	err2 := json.Unmarshal(data, v)
	if err2 != nil {
		return err2
	}
	return nil
}
