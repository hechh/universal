package player

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"universal/common/global"
	"universal/common/pb"
	"universal/framework/handler"
	"universal/library/mlog"
	"universal/tools/client/domain"

	"github.com/spf13/cast"
)

var (
	platform uint64
	url      string = "ws://172.16.126.208:801/ws" // 内网域名
	gates    []uint64
	gatecfg  []*global.ServerConfig
	rwMutex  = new(sync.RWMutex)
	uids     = make(map[uint64]*Player)
	closes   = make(chan uint64, 100)
)

func RandomGate(uid uint64) string {
	switch platform {
	case 1:
		return url
	}
	ll := uint64(len(gates))
	gateCfg := gatecfg[uint32(gates[uid%ll])]
	return fmt.Sprintf("ws://%s/ws", gateCfg.Host)
}

func Init(cfg map[uint32]*global.ServerConfig, plat uint64, srvs string) {
	platform = plat
	for _, val := range cfg {
		gatecfg = append(gatecfg, val)
	}
	for _, val := range strings.Split(srvs, ",") {
		gates = append(gates, cast.ToUint64(val))
	}
	go func() {
		for uid := range closes {
			DelPlayer(uid)
			mlog.Info("客户端关闭, uid: %d", uid)
		}
	}()
}

func Close() {
	rwMutex.Lock()
	defer rwMutex.Unlock()
	for uid, pl := range uids {
		pl.Close()
		delete(uids, uid)
	}
}

func GetOnline() int {
	return len(uids)
}

func GetPlayer(uid uint64) *Player {
	rwMutex.RLock()
	defer rwMutex.RUnlock()
	return uids[uid]
}

func AddPlayer(pl *Player) {
	rwMutex.Lock()
	defer rwMutex.Unlock()
	uids[pl.GetUID()] = pl
}

func DelPlayer(uid uint64) {
	rwMutex.Lock()
	defer rwMutex.Unlock()
	if val, ok := uids[uid]; ok {
		val.Close()
		delete(uids, uid)
	}
}

func Send(wg *sync.WaitGroup, name, strJson string, cb domain.ResultCB) error {
	pac := handler.GetByName(name)
	if pac == nil {
		return fmt.Errorf("协议(%s)不支持", name)
	}

	// 解析请求
	req := pac.NewRequest()
	if err := json.Unmarshal([]byte(strJson), req); err != nil {
		return fmt.Errorf("json字符串解析错误, name: %s, json: %s", name, strJson)
	}

	// 发送请求
	rwMutex.RLock()
	defer rwMutex.RUnlock()
	for _, pl := range uids {
		wg.Add(1)
		pl.Send(&pb.RpcHead{FuncName: name}, req, cb)
	}
	return nil
}

func Login(uid uint64, cb domain.ResultCB) {
	// 创建玩家
	if pl, err := NewPlayer(nil, uid, closes); err != nil {
		cb(&domain.Result{UID: uid, Error: err})
	} else {
		// 添加管理
		AddPlayer(pl)
		// 开始登录
		pl.Login(cb)
	}
}
