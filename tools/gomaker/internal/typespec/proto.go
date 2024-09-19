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
	Key     string
	Value   string
	Comment string
}

type Import struct {
	Docs    []string
	File    string
	Comment string
}

type Attribute struct {
	Type     string
	Name     string
	IsRepeat bool
	Index    int
}

type Message struct {
	Doc        string
	Name       string
	Attributes []*Attribute
	Comment    string
}

type MValue struct {
	Name    string
	Value   int
	Comment string
}

type MEnum struct {
	Doc     string
	Name    string
	Comment string
}
