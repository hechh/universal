package typespec

type Syntax struct {
	Docs    []string
	Type    string
	Name    string
	Comment string
}

type Package struct {
	Docs    []string
	Type    string
	Name    string
	Comment string
}

type Option struct {
	Docs    []string
	Type    string
	Key     string
	Value   string
	Comment string
}

type Import struct {
	Docs    []string
	Type    string
	File    string
	Comment string
}

type Attribute struct {
	Docs     []string
	Type     string
	Name     string
	Comment  string
	IsRepeat bool
	Index    int
}

type Message struct {
	Docs       []string
	Type       string
	Name       string
	Attributes []*Attribute
}

type MValue struct {
	Docs    []string
	Name    string
	Value   int
	Comment string
}

type MEnum struct {
	Docs   []string
	Type   string
	Name   string
	Values []*MValue
}
