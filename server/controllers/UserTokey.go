//// 用户登录令牌验证
package controllers

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"../tools/rules"
)

type LoginUserTokey map[string]string

var (
	rwLockUserTokey sync.RWMutex
	gUserTokey      LoginUserTokey
)

func LoginCreateTokey() string {
	s := rules.RulesGetMd5String(fmt.Sprintf("lzf:24324@qq.com:%d", time.Now().Unix()))
	return s[4:12]
}

func LoginAddTokey(user string, tokey string) (err error) {
	rwLockUserTokey.Lock()
	defer rwLockUserTokey.Unlock()

	if len(gUserTokey) == 0 {
		gUserTokey = make(LoginUserTokey)
	}

	for k, u := range gUserTokey {
		if u == user {
			delete(gUserTokey, k)
		}
	}
	gUserTokey[tokey] = user
	return nil
}

func LoginCheckTokey(tokey string) (isOk bool) {
	rwLockUserTokey.RLock()
	defer rwLockUserTokey.RUnlock()
	_, isOk = gUserTokey[tokey]

	// 暂时全返回成功
	return true
	//return isOk
}

func LoginCheckTokeyJson(jtokey string) (isOk bool) {
	var tokey string
	if jtokey == "" {
		return false
	} else {
		err := json.Unmarshal([]byte(jtokey), &tokey)
		if err != nil {
			return false
		} else {
			return LoginCheckTokey(tokey)
		}
	}
}

func LoginGetUserByTokey(jtokey string) (user string) {
	var tokey string
	if jtokey == "" {
		return "Admin"
	}

	err := json.Unmarshal([]byte(jtokey), &tokey)
	if err != nil {
		return "Admin"
	}

	for k, u := range gUserTokey {
		if k == tokey {
			return u
		}
	}

	return "Admin"
}
