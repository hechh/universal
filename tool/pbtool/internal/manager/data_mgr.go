package manager

import "universal/tool/pbtool/domain"

var (
	mgr = make(map[string]domain.IClass)
)

func Register(cls domain.IClass) {
	mgr[cls.GetName()] = cls
}

func Get(name string) domain.IClass {
	return mgr[name]
}
