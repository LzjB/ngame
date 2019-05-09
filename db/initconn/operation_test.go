package initconn

import (
	"fmt"
	std "ngame/pubconfig"
	"testing"
)

type testingT struct {
	Name string
}

func TestInitDb(t *testing.T) {
	con := &std.DbConfig{Redis: "127.0.0.1:6379"}
	InitDb(con)

	id := uint32(1)
	for ; id < 100; id++ {
		setToRedis(id, &testingT{Name: string("Name" + string(id))})
	}

	for id = 1; id < 100; id++ {
		//time.Sleep(time.Second * 2)
		_, err := getInRedis(id)
		if err != nil {
			fmt.Println("getInRedis err : ", err)
		} else {
			//fmt.Println("id : ",id,"  data : ",data)
		}
	}

	for id = 1; id < 100; id++ {
		delInRedis(id)
	}
}
