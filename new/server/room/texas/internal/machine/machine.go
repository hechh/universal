package machine

import (
	"fmt"
	"poker_server/common/pb"
	"poker_server/server/room/texas/domain"
	"poker_server/server/room/texas/internal/manager"
	"time"
)

type TexasMachine struct {
	startTime int64             // 状态机开始时间
	curState  pb.TexasGameState // 当前状态
}

func NewTexasMachine(nowMs int64, state pb.TexasGameState, room domain.IRoom) *TexasMachine {
	// 初始化状态
	handle := manager.GetState(state)
	if handle == nil {
		panic(fmt.Sprintf("TexasMachine状态机未注册状态: %d", state))
	}
	handle.OnEnter(nowMs, state, room)

	// 创建状态机
	return &TexasMachine{curState: state, startTime: time.Now().UnixMilli()}
}

func (m *TexasMachine) SetStateStartTime(nowMs int64) {
	m.startTime = nowMs
}

func (m *TexasMachine) GetStateStartTime() int64 {
	return m.startTime
}

func (m *TexasMachine) GetCurState() pb.TexasGameState {
	return m.curState
}

// 状态转移
func (m *TexasMachine) moveTo(newState pb.TexasGameState, nowMs int64, room domain.IRoom) {
	// 执行旧状态退出逻辑
	state := manager.GetState(m.curState)
	if state == nil {
		panic(fmt.Sprintf("状态机未注册状态: %d", newState))
	}
	state.OnExit(nowMs, m.curState, room)

	// 更新状态时间
	m.curState = newState
	m.startTime = nowMs

	// 执行进入逻辑
	handle := manager.GetState(m.curState)
	if handle == nil {
		panic(fmt.Sprintf("TexasMachine状态机未注册状态: %d", newState))
	}
	handle.OnEnter(nowMs, m.curState, room)
}

func (m *TexasMachine) Handle(nowMs int64, room domain.IRoom) {
	// 获取状态机
	handle := manager.GetState(m.curState)
	if handle == nil {
		panic(fmt.Sprintf("TexasMachine状态机未注册状态: %d", m.curState))
	}

	// 状态处理
	nextState := handle.OnTick(nowMs, m.curState, room)
	if nextState != m.curState {
		m.moveTo(nextState, nowMs, room)
	}
}
