package cfgtool

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"universal/framework/common/uerror"
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
			data[node.Key] = base.CastTo(node.Type.Name, val)
		} else {
			if vv, ok := data[node.Key]; ok {
				data[node.Key] = append(vv.([]interface{}), base.CastTo(node.Type.Name, val))
			} else {
				data[node.Key] = []interface{}{base.CastTo(node.Type.Name, val)}
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
	buffer := bytes.NewBuffer(nil)
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
			buffer.Write(buf)
			buffer.WriteRune('\n')
		}
		return true
	})
	if err != nil {
		return err
	}
	// 生成文档
	dstFile := cmdLine.Dst
	if !strings.HasSuffix(dstFile, ".json") {
		dstFile += ".json"
	}
	if err := os.MkdirAll(filepath.Dir(dstFile), os.FileMode(0777)); err != nil {
		return uerror.NewUError(1, -1, err)
	}
	if err := ioutil.WriteFile(dstFile, buffer.Bytes(), os.FileMode(0666)); err != nil {
		return uerror.NewUError(1, -1, err)
	}
	return nil
}

func Init() {
	manager.Register("json", maker.NewBaseMaker(Gen, "-action=json -param=test.xlsx", "xlsx表转json"))
}
