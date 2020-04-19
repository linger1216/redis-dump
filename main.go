package main

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/configor"
	"github.com/linger1216/redis-dump/core"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	configFileName = kingpin.Flag("file", "config yaml file").
		Short('f').Required().String()
)

func main() {

	kingpin.Parse()
	conf := &core.DumpConfig{}
	err := configor.Load(conf, *configFileName)
	if err != nil {
		//panic(fmt.Errorf("need config.json file"))
	}

	buf, _ := json.Marshal(conf)
	fmt.Println(string(buf))

	err = core.Exec(conf)
	if err != nil {
		panic(err)
	}
}
