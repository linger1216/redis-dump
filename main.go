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

	conf := core.DumpConfig{}
	err := configor.Load(&conf, *configFileName)
	if err != nil {
		panic(fmt.Errorf("need config yaml file"))
	}

	buf, _ := json.Marshal(conf)
	fmt.Println(string(buf))

	// make tasks

	// 集群是不支持db number除了0以外的
	//xx := redis.NewClusterClient(nil)
	//client := redis.NewClient(&redis.Options{
	//	network:  conf.Src.network,
	//	Addr:     conf.Src.url[0],
	//	password: "",
	//	DB:       0, // use default DB
	//})
	//
	//pong, err := client.Ping().Result()
	//fmt.Println(pong, err)

	//xx := client.Scan(0,"",100)
	//xx.String()
	//xx.Val()

	//core.dumpKeys(client, []string{"35f1ab794", "36d083e44"}, true, core.RedisCmdSerializer)

}
