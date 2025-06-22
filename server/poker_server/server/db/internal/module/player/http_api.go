package player

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"poker_server/common/pb"
	"poker_server/common/yaml"
	"poker_server/library/uerror"
	"time"

	"github.com/spf13/cast"
)

var (
	phpCfg *yaml.PhpConfig
)

func Init(cfg *yaml.PhpConfig) {
	phpCfg = cfg
}

func PlayerInfoRequest(uid uint64, rsp *pb.HttpPlayerInfoRsp) error {
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	if err := writer.WriteField("uid", cast.ToString(uid)); err != nil {
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}

	// 1. 检查配置请求
	req, err := http.NewRequest("POST", phpCfg.UserInfoUrl, payload)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 3. 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body) // 尝试读取错误信息
		return uerror.New(1, pb.ErrorCode_REQUEST_FAIELD, "uid:%d 服务器返回错误状态码: %d, 响应: %s", uid, resp.StatusCode, string(body))
	}
	// 4. 解析JSON到protobuf结构体
	if err := json.NewDecoder(resp.Body).Decode(rsp); err != nil {
		return uerror.New(1, pb.ErrorCode_MARSHAL_FAILED, "uid:%d JSON解析失败: %v", uid, err)
	}
	if rsp.RespMsg.Code != 100 {
		return uerror.New(1, pb.ErrorCode_REQUEST_FAIELD, "uid:%d 请求失败: %s", uid, rsp.RespMsg.Message)
	}
	return nil
}
