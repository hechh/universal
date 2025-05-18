package mock

import (
	"fmt"
	"poker_server/common/config"
	"poker_server/common/dao"
	"poker_server/common/dao/repository/redis/login_session"
	"poker_server/common/pb"
	"poker_server/common/yaml"
	"poker_server/framework/domain"
	"poker_server/framework/library/async"
	"poker_server/framework/library/mlog"
	"poker_server/framework/network"
	"poker_server/server/gate/frame"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/spf13/cast"
	"golang.org/x/net/websocket"
)

type MockClient struct {
	node      *pb.Node
	uid       uint64
	sessionId string
	cfg       *yaml.Config
	conn      domain.INet
}

func NewMockClient(cfg string, nodeId int32, uid uint64) (*MockClient, error) {
	node := &pb.Node{
		Name: strings.ToLower(pb.NodeType_Gate.String()),
		Type: pb.NodeType_Gate,
		Id:   int32(nodeId),
	}
	// 加载配置
	yamlCfg, err := yaml.LoadConfig(cfg, node)
	if err != nil {
		return nil, err
	}
	// 初始化redis
	if err := dao.InitRedis(yamlCfg.Redis); err != nil {
		return nil, err
	}
	// 初始化游戏配置
	if err := config.InitLocal(yamlCfg.Configure); err != nil {
		return nil, err
	}
	ret := &MockClient{
		node: node,
		uid:  uid,
		cfg:  yamlCfg,
	}
	if err := ret.init(); err != nil {
		return nil, err
	}
	return ret, nil
}

// 发送消息
func (d *MockClient) Send(head *pb.Head, msg proto.Message) error {
	head.Id = d.uid
	head.RouteId = d.uid
	return d.conn.WriteMsg(head, msg)
}

func (d *MockClient) Close() error {
	return d.conn.Close()
}

func (d *MockClient) init() error {
	// 设置 session
	now := time.Now().Unix()
	d.sessionId = cast.ToString(now)
	ttl := int64(15)
	if err := login_session.Set(d.sessionId, &pb.LoginSession{
		SessionId:  d.sessionId,
		AccountId:  d.uid,
		CreateTime: now,
		ExpireTime: now + ttl,
	}, time.Duration(ttl)*time.Second); err != nil {
		return err
	}
	// 建立连接
	ws, err := websocket.Dial(fmt.Sprintf("ws://localhost%s/ws", d.node.Addr), "", "http://localhost/")
	if err != nil {
		return err
	}
	d.conn = network.NewSocket(ws, 1024*1024)
	d.conn.Register(&frame.Frame{})
	// 发送登录请求
	if err := d.login(); err != nil {
		return err
	}
	// 启动心跳
	async.SafeGo(mlog.Fatalf, d.keepAlive)
	return nil
}

func (d *MockClient) keepAlive() {
	head := &pb.Head{
		SendType:    pb.SendType_POINT,
		DstNodeType: pb.NodeType_Gate,
		DstNodeId:   d.node.Id,
		IdType:      pb.IdType_UID,
		Id:          d.uid,
		RouteId:     d.uid,
		Cmd:         uint32(pb.CMD_GATE_HEART_REQUEST),
	}
	// 循环发送心跳
	tt := time.NewTicker(3 * time.Second)
	for {
		<-tt.C
		if err := d.conn.WriteMsg(head, &pb.GateHeartRequest{}); err != nil {
			fmt.Println("发送心跳包失败:", err)
			break
		}
	}
}

func (d *MockClient) login() error {
	// 发送登录消息
	head := &pb.Head{
		DstNodeType: d.node.Type,
		DstNodeId:   d.node.Id,
		Id:          d.uid,
		RouteId:     d.uid,
		Cmd:         uint32(pb.CMD_GATE_LOGIN_REQUEST),
	}
	if err := d.conn.WriteMsg(head, &pb.GateLoginRequest{SessionId: d.sessionId}); err != nil {
		return err
	}
	// 接收登录返回消息
	pack := &pb.Packet{}
	if err := d.conn.Read(pack); err != nil {
		return err
	}
	loginRsp := &pb.GateLoginResponse{}
	if err := proto.Unmarshal(pack.Body, loginRsp); err != nil {
		return err
	}
	if loginRsp.Head != nil && loginRsp.Head.Code != 0 {
		return fmt.Errorf("登录失败: %v", loginRsp.Head)
	}
	return nil
}
