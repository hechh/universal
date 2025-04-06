package main

import (
	"encoding/binary"
	"flag"
	"hego/Library/ulog"
	"hego/common/dao"
	"hego/common/global"
	"hego/common/pb"
	"hego/framework/basic"
	"hego/framework/cluster"
	"hego/framework/socket"
	"net/http"
	"os"

	"golang.org/x/net/websocket"
	"google.golang.org/protobuf/proto"
)

type Frame struct{}

func (d *Frame) GetHeadSize() int {
	return 4
}

func (d *Frame) GetBodySize(head []byte) int {
	return int(binary.BigEndian.Uint32(head))
}

func (d *Frame) Check(head []byte, body []byte) bool {
	return true
}

func (d *Frame) Build(frame []byte, body []byte) []byte {
	// 设置包头
	binary.BigEndian.PutUint32(frame, uint32(len(body)))
	// 拷贝
	headSize := d.GetHeadSize()
	copy(frame[headSize:], body)
	return frame
}

func main() {
	var id int64
	var path string
	flag.Int64Var(&id, "id", 1, "服务节点ID")
	flag.StringVar(&path, "cfg", "./", "yaml配置文件路径")
	flag.Parse()

	// 加载配置
	if err := global.Init(path, pb.SERVER_Gate, uint32(id)); err != nil {
		panic(err)
	}

	// 初始化plog
	if logCfg := global.GetLogCfg(); logCfg != nil {
		ulog.Init(logCfg.Level, logCfg.Path, global.GetServerName())
	}

	// 初始化redis
	if err := dao.InitRedis(global.GetRedisCfg()); err != nil {
		panic(err)
	}

	// 初始化集群
	if err := cluster.Init(global.GetCfg(), global.GetServerType(), uint32(id), 600); err != nil {
		panic(err)
	}

	// 初始化websocket
	http.Handle("/ws", websocket.Handler(wsHandle))
	go func() {
		srvCfg := global.GetServerCfg()
		if srvCfg == nil {
			panic("yaml配置加载错误")
		}
		if err := http.ListenAndServe(srvCfg.Host, nil); err != nil {
			panic(err)
		}
	}()

	// 等待结束
	basic.SignalHandle(func(s os.Signal) {
		ulog.Info("gate服务退出")
		ulog.Close()
	}, os.Interrupt, os.Kill)
}

// websocket包处理
func wsHandle(ws *websocket.Conn) {
	ulog.Info("客户端(%s)连接成功！！！", ws.RemoteAddr().String())
	soc := socket.NewSocket(&Frame{}, ws)

	// 循环接受处理消息
	basic.SafeGo(ulog.Catch, func() {
		for {
			// 接受请求
			buf, err := soc.Read()
			if err != nil {
				ulog.Error("websocket接受数据失败: %v", err)
				break
			}

			// 解析packet
			pac := new(pb.Packet)
			if err := proto.Unmarshal(buf, pac); err != nil {
				ulog.Error("协议错误: %v", err)
				break
			}

			// 对请求路由
			if err := cluster.Dispatcher(pac.Head); err != nil {
				ulog.Error("路由错误: %v, error: %v", pac.Head, err)
				break
			}
		}
	})
}
