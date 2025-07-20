package http

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"poker_server/common/config/repository/php_config"
	"poker_server/common/pb"
	"poker_server/framework"
	"poker_server/library/mlog"
	"poker_server/library/uerror"
	"time"

	"github.com/spf13/cast"
)

func ChargeTransInRequest(param *pb.TransParam, rsp *pb.HttpTransferInRsp) error {
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	if err := writer.WriteField("uid", cast.ToString(param.Uid)); err != nil {
		return err
	}
	if err := writer.WriteField("game_sn", cast.ToString(param.GameSn)); err != nil {
		return err
	}
	if err := writer.WriteField("incr", cast.ToString(param.Incr)); err != nil {
		return err
	}
	if err := writer.WriteField("game_type", cast.ToString(int32(param.GameType))); err != nil {
		return err
	}
	if err := writer.WriteField("coin_type", cast.ToString(int32(param.CoinType))); err != nil {
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}

	// 1. 检查配置请求
	req, err := http.NewRequest("POST", php_config.MGetEnvTypeName(framework.GetEnvType(), "TransferOutUrl").Url, payload)
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
		return uerror.New(1, pb.ErrorCode_REQUEST_FAIELD, "uid:%d 服务器返回错误状态码: %d, 响应: %s", param.Uid, resp.StatusCode, string(body))
	}

	body, _ := io.ReadAll(resp.Body)

	// 4. 解析JSON到protobuf结构体
	if err = json.Unmarshal(body, rsp); err != nil {
		mlog.Infof("chargeTransInRequest resp:%v", string(body))
		mlog.Infof("chargeTransInRequest url:%v", php_config.MGetEnvTypeName(framework.GetEnvType(), "TransferOutUrl").Url)
		return uerror.New(1, pb.ErrorCode_MARSHAL_FAILED, "uid:%d JSON解析失败: %v", param.Uid, err)
	}
	if rsp.RespMsg.Code != 100 {
		return uerror.New(1, pb.ErrorCode_REQUEST_FAIELD, "uid:%d 请求失败ERROR : %v", param.Uid, rsp.RespMsg)
	}
	return nil
}

func ChargeTransOutRequest(param *pb.TransParam, rsp *pb.HttpTransferOutRsp) error {
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	if err := writer.WriteField("uid", cast.ToString(param.Uid)); err != nil {
		return err
	}
	if err := writer.WriteField("game_sn", cast.ToString(param.GameSn)); err != nil {
		return err
	}
	if err := writer.WriteField("incr", cast.ToString(param.Incr)); err != nil {
		return err
	}
	if err := writer.WriteField("game_type", cast.ToString(int32(param.GameType))); err != nil {
		return err
	}
	if err := writer.WriteField("coin_type", cast.ToString(int32(param.CoinType))); err != nil {
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}

	// 1. 检查配置请求 todo
	req, err := http.NewRequest("POST", php_config.MGetEnvTypeName(framework.GetEnvType(), "TransferInUrl").Url, payload)
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
		return uerror.New(1, pb.ErrorCode_REQUEST_FAIELD, "uid:%d 服务器返回错误状态码: %d, 响应: %s", param.Uid, resp.StatusCode, string(body))
	}

	body, _ := io.ReadAll(resp.Body)
	// 4. 解析JSON到protobuf结构体
	if err = json.Unmarshal(body, rsp); err != nil { // 尝试读取错误信息
		mlog.Infof("ChargeTransOutRequest resp:%v", string(body))
		mlog.Infof("ChargeTransOutRequest url:%v", php_config.MGetEnvTypeName(framework.GetEnvType(), "TransferOutUrl").Url)
		return uerror.New(1, pb.ErrorCode_MARSHAL_FAILED, "uid:%d JSON解析失败: %v", param.Uid, err)
	}
	if rsp.RespMsg.Code != 100 {
		return uerror.New(1, pb.ErrorCode_REQUEST_FAIELD, "uid:%d 请求失败: %s", param.Uid, rsp.RespMsg.Message)
	}
	return nil
}
