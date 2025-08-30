package manager

import "universal/tool/dbtool/internal/base"

var (
	stringMgr = make(map[string]*base.String)
	hashMgr   = make(map[string]*base.Hash)
)

func AddString(aa *base.String) {
	stringMgr[aa.Name] = aa
}

func WalkString(f func(*base.String) bool) {
	for _, v := range stringMgr {
		if !f(v) {
			break
		}
	}
}

func AddHash(aa *base.Hash) {
	hashMgr[aa.Name] = aa
}

func WalkHash(f func(*base.Hash) bool) {
	for _, v := range hashMgr {
		if !f(v) {
			break
		}
	}
}
