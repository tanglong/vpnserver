package main

import (
	"base/common"
	l4g "base/log4go"
	"flag"
	"fmt"
	"github.com/go-redis/redis"
	"net/http"
	_ "net/http/pprof"
)

var (
	g_config = new(xmlConfig)
	g_redis  *redis.Client
)

func init() {
}

var configFile = flag.String("config", "../config/vpnserver_config.xml", "")

func main() {

	flag.Parse()

	if err := common.LoadConfig(*configFile, g_config); err != nil {
		panic(fmt.Sprintf("load config %v fail: %v", *configFile, err))
	}

	l4g.LoadConfiguration(g_config.Log.Config)
	defer l4g.Close()

	rd := redis.NewClient(&redis.Options{
		Addr:     g_config.Redis.Address,
		Password: "", // no password set
		DB:       g_config.Redis.Index,
	})
	g_redis = rd
	pong, err := g_redis.Ping().Result()
	l4g.Info("pong :[%s] err:[%v]", pong, err)

	defer g_redis.Close()

	httpRouter()
	l4g.Error(http.ListenAndServe(g_config.Http, nil))
}
