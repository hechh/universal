package cfgtool

import (
	"encoding/json"
	"fmt"
	"strings"
	"universal/tools/gomaker/domain"
	"universal/tools/gomaker/internal/common/base"
	"universal/tools/gomaker/internal/common/types"
	"universal/tools/gomaker/internal/maker"
	"universal/tools/gomaker/internal/manager"

	"github.com/spf13/cast"
)

type Node struct {
	Key   string
	Index int
	Type  *types.Type
	next  *Node
}

func ToCast(t *types.Type, str string) interface{} {
	switch t.Name {
	case "int":
		return cast.ToInt(str)
	case "int32":
		return cast.ToInt32(str)
	case "int64":
		return cast.ToInt64(str)
	case "uint":
		return cast.ToUint(str)
	case "uint32":
		return cast.ToUint32(str)
	case "uint64":
		return cast.ToUint64(str)
	case "string":
		return str
	case "float32":
		return cast.ToFloat32(str)
	case "float64":
		return cast.ToFloat64(str)
	}
	return nil
}

func ParseType(stname string, node *Node) {
	if val := manager.GetStruct(stname); val != nil {
		vv, ok := val.Fields[node.Key]
		if !ok {
			return
		}
		if node.next == nil {
			node.Type = vv.Type
			return
		}
		ParseType(vv.Type.Name, node.next)
	}
}

func ParseNode(key string) *Node {
	pos := strings.Index(key, ".")
	if pos == -1 {
		index := 0
		if pos = strings.Index(key, "#"); pos != -1 {
			index = cast.ToInt(key[pos+1:])
			key = key[:pos]
		}
		return &Node{Key: key, Index: index}
	}
	result := ParseNode(key[:pos])
	result.next = ParseNode(key[pos+1:])
	return result
}

func parse(data, record map[string]interface{}, node *Node, val string) {
	if node.next == nil {
		if node.Index <= 0 {
			data[node.Key] = ToCast(node.Type, val)
		} else {
			if vv, ok := data[node.Key]; ok {
				data[node.Key] = append(vv.([]interface{}), ToCast(node.Type, val))
			} else {
				data[node.Key] = []interface{}{ToCast(node.Type, val)}
			}
		}
		return
	}
	if node.Index <= 0 {
		if _, ok := data[node.Key]; !ok {
			data[node.Key] = map[string]interface{}{}
		}
		parse(data[node.Key].(map[string]interface{}), record, node.next, val)
		return
	}
	rkey := fmt.Sprintf("%s#%d", node.Key, node.Index)
	if _, ok := record[rkey]; !ok {
		record[rkey] = map[string]interface{}{}
	}
	if vv, ok := data[node.Key]; !ok {
		data[node.Key] = []interface{}{record[rkey]}
	} else {
		data[node.Key] = append(vv.([]interface{}), record[rkey])
	}
	parse(record[rkey].(map[string]interface{}), record, node.next, val)
}

func Gen(cmdLine *domain.CmdLine, tpls *base.Templates) error {
	nodes := map[int]*Node{}
	jsons := map[int]string{}
	err := base.ParseXlsx(cmdLine.Param, func(sheet string, row int, cols []string) bool {
		switch row {
		case 0:
		case 1:
			for index, val := range cols {
				if len(val) > 0 {
					nodes[index] = ParseNode(val)
					ParseType(sheet, nodes[index])
				}
			}
		default:
			data := map[string]interface{}{}
			record := map[string]interface{}{}
			for index, val := range cols {
				if node, ok := nodes[index]; ok {
					parse(data, record, node, val)
				}
			}
			buf, _ := json.Marshal(data)
			jsons[row] = string(buf)
		}
		return true
	})
	if err != nil {
		return err
	}
	fmt.Println("------>", jsons)
	return nil
}

func Init() {
	manager.Register("xlsx", maker.NewBaseMaker(Gen, "-action=xlsx -param=test.xlsx", "xlsx表转json"))
}
