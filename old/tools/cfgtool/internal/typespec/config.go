package typespec

type Index struct {
	Name      string
	Type      *Type
	FieldList []*Field
}

type Config struct {
	Name      string
	Fields    map[string]*Field
	FieldList []*Field
	Indexs    map[string]*Index
	IndexList []*Index
	FileName  string
	Sheet     string
	Rows      [][]string
	Rules     []string
}

func (c *Config) Set(file, sheet string, datas [][]string, rules []string) {
	c.FileName = file
	c.Sheet = sheet
	c.Rows = datas
	c.Rules = rules
}

func (d *Config) AddField(field *Field) {
	d.FieldList = append(d.FieldList, field)
	d.Fields[field.Name] = field
}

func (d *Config) AddIndex(ind *Index) {
	d.Indexs[ind.Name] = ind
	d.IndexList = append(d.IndexList, ind)
}
