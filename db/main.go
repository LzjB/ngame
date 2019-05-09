package main

import (
	"fmt"
	"log"
	"ngame/db/initconn"
	std "ngame/pubconfig"
	"os"
	"path"
)

func main() {
	dir, _ := os.Getwd()
	cfg := path.Join(dir, "db.json")

	jsonParse := std.NewJsonCode()
	v := &std.DbConfig{}
	err := jsonParse.Load(cfg, v)
	if err != nil {
		log.Fatal("Load err : ", err)
	}
	fmt.Println(v.Listen)
	initconn.InitDb(v)
}
