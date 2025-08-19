package typespec

type Type struct {
	Name    string
	TypeOf  int
	ValueOf int
}

type Field struct {
	Name     string
	Type     *Type
	Desc     string
	Position int
	ConvFunc func(string) interface{}
}

func (d *Field) Convert(vals ...string) (rets []interface{}) {
	for _, val := range vals {
		rets = append(rets, d.ConvFunc(val))
	}
	return
}

type Struct struct {
	Name      string
	Fields    map[string]*Field
	FieldList []*Field
	Converts  map[string][]*Field
	FileName  string
	Sheet     string
	Rows      [][]string
}

func (d *Struct) Set(file, sheet string, datas [][]string) {
	d.FileName = file
	d.Sheet = sheet
	d.Rows = datas
}

func (d *Struct) AddField(field *Field) {
	d.FieldList = append(d.FieldList, field)
	d.Fields[field.Name] = field
}
