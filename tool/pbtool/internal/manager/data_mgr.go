package manager

import "universal/tool/pbtool/domain"

var (
	mgr  = make(map[string]domain.IClass)
	list = []domain.IClass{}
)

func Register(cls domain.IClass) {
	mgr[cls.GetName()] = cls
	list = append(list, cls)
}

func Get(name string) domain.IClass {
	return mgr[name]
}

func GetAll() []domain.IClass {
	return list
}
