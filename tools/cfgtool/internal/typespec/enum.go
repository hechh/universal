package typespec

type Value struct {
	Name  string
	Value int32
	Desc  string
}

type Enum struct {
	Name      string
	Values    map[string]*Value
	ValueList []*Value
	FileName  string
	Sheet     string
}

func (e *Enum) Set(file, sheet string) {
	e.FileName = file
	e.Sheet = sheet
}

func (e *Enum) AddValue(name string, val int32, desc string) {
	if e.Values == nil {
		e.Values = make(map[string]*Value)
	}
	value := &Value{Name: name, Value: val, Desc: desc}
	e.ValueList = append(e.ValueList, value)
	e.Values[name] = value
}
