package packet

type Header struct {
	SrcNodeType uint32
	SrcNodeId   uint32
	DstNodeType uint32
	DstNodeId   uint32
	Uid         uint64
	RouteId     uint64
	Cmd         uint32
	ActorName   string
	FuncName    string
}

func (h *Header) GetSrcNodeType() uint32 {
	return h.SrcNodeType
}

func (h *Header) GetSrcNodeId() uint32 {
	return h.SrcNodeId
}

func (h *Header) GetDstNodeType() uint32 {
	return h.DstNodeType
}

func (h *Header) GetDstNodeId() uint32 {
	return h.DstNodeId
}

func (h *Header) GetCmd() uint32 {
	return h.Cmd
}

func (h *Header) GetUid() uint64 {
	return h.Uid
}

func (h *Header) GetRouteId() uint64 {
	return h.RouteId
}

func (h *Header) GetActorName() string {
	return h.ActorName
}

func (h *Header) GetFuncName() string {
	return h.FuncName
}
