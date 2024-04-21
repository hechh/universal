package basic

import "universal/common/pb"

type Context struct {
	*pb.PacketHead                        // rpc请求头
	readyOnlys     map[string]interface{} // 零时缓存
}

func NewContext(head *pb.PacketHead, datas map[string]interface{}) *Context {
	return &Context{
		PacketHead: head,
		readyOnlys: datas,
	}
}

func (d *Context) GetValue(key string) interface{} {
	return d.readyOnlys[key]
}
