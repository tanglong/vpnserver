package main

import (
	l4g "base/log4go"
	//"bytes"
	"encoding/json"
	"fmt"
	//"github.com/go-redis/redis"
	"math/rand"
	"net/http"
	"os/exec"
	"time"
)

func httpRouter() {
	http.HandleFunc("/queryaccount", QueryAccountHandler)
}

func QueryAccountHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	deviceid := r.Form.Get("deviceid")

	ret, pd, et := GetUserInfo(deviceid)
	if ret {
		l4g.Info("QueryAccountCodeHandler ret:[%d] pd:[%s] et:[%d]", ret, pd, et)
		var result = make(map[string]interface{})
		result["username"] = deviceid
		result["password"] = "bgm"
		result["psk"] = "psk"
		result["serverip"] = g_config.OutIp
		show_result, _ := json.Marshal(result)
		w.Write(show_result)
		return
	}

	cmd1 := fmt.Sprintf(`docker exec -t -e "ACCOUNT=%s" -e"PASSWORD=%s" cf4439b109b4 gen-account`, deviceid, "123456")
	cmd2 := `docker exec -t  cf4439b109b4 ipsec reload`
	output, err := exec.Command("/bin/sh", "-c", cmd1).Output()
	if err != nil {
		l4g.Error("exec command fail: %s %v %s", cmd1, err, output)
		return
	}
	output, err = exec.Command("/bin/sh", "-c", cmd2).Output()
	if err != nil {
		l4g.Error("exec command fail: %s %v %s", cmd2, err, output)
		return
	}

	key := "user:" + deviceid

	var letterRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, 18)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	passwd := string(b)

	err = g_redis.HSet(key, "passwd", passwd).Err()
	if err != nil {
		l4g.Error("client hset passwd fail: %v ", err)
		return
	}

	now := uint32(time.Now().Unix())
	expiretime := now + 7*24*60*60

	err = g_redis.HSet(key, "expiretime", expiretime).Err()
	if err != nil {
		l4g.Error("client hset expiretime fail: %v ", err)
		return
	}

	err = g_redis.SAdd("users", deviceid).Err()
	if err != nil {
		l4g.Error("users sadd fail: %v ", err)
		return
	}

	var result = make(map[string]interface{})
	result["username"] = deviceid
	result["password"] = passwd
	result["psk"] = "654321"
	result["serverip"] = g_config.OutIp
	show_result, _ := json.Marshal(result)
	w.Write(show_result)
}

// Ascii numbers 0-9
const (
	ascii_0 = 48
	ascii_9 = 57
)

func ParseUint64(d []byte) (uint64, bool) {
	var n uint64
	d_len := len(d)
	if d_len == 0 {
		return 0, false
	}
	for i := 0; i < d_len; i++ {
		j := d[i]
		if j < ascii_0 || j > ascii_9 {
			return 0, false
		}
		n = n*10 + (uint64(j - ascii_0))
	}
	return n, true
}

func GetUserInfo(username string) (bool, string, uint64) {
	key := "user:" + username
	fields, err := g_redis.HMGet(key, "passwd", "expiretime").Result()
	if err != nil {
		return false, "", 0
	}
	l4g.Info("GetUserInfo %v", fields)
	var expiretime uint64
	str, _ := fields[1].(string)
	if v, e := ParseUint64([]byte(str)); !e {
		l4g.Error("get user info parse expiretime error: %v", e)
		return false, "", 0
	} else {
		expiretime = v
	}

	passwd, _ := fields[0].(string)
	return true, passwd, expiretime
}
