package typespec

type Proto struct {
	FileName      string
	References    map[string]struct{}
	ReferenceList []string
	Enums         []*Enum
	Structs       []*Struct
	Configs       []*Config
}

func (p *Proto) AddReference(ref string) {
	if p.FileName == ref {
		return
	}
	if _, ok := p.References[ref]; ok {
		return
	}
	p.ReferenceList = append(p.ReferenceList, ref)
	p.References[ref] = struct{}{}
}

func (p *Proto) AddEnum(enum *Enum) {
	p.Enums = append(p.Enums, enum)
}

func (p *Proto) AddStruct(strct *Struct) {
	p.Structs = append(p.Structs, strct)
}

func (p *Proto) AddConfig(cfg *Config) {
	p.Configs = append(p.Configs, cfg)
}
