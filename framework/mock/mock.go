package mock

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/golang/protobuf/proto"
	"github.com/spf13/cast"
)

var (
	ip   string = "localhost"
	port int    = 22345
)

func Init(ip string, port int) {
	ip = ip
	port = port
}

func Request(cmd uint32, msg proto.Message) error {
	buf, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	targetUrl := fmt.Sprintf("http://%s:%d/api", ip, port)
	u, _ := url.ParseRequestURI(targetUrl)

	params := url.Values{}
	params.Set("cmd", cast.ToString(cmd))
	params.Set("value", string(buf))

	u.RawQuery = params.Encode()
	_, err = http.Get(u.String())
	return err
}
