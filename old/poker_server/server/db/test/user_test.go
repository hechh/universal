package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"poker_server/common/pb"
	"poker_server/framework/mock"
	"testing"
	"time"

	"github.com/spf13/cast"
)

func TestUserInfo(t *testing.T) {
	uid := uint64(144)
	routeId := uint64(pb.DataType_DataTypeUserInfo)
	t.Run("UserInfoReq", func(t *testing.T) {
		mock.Request(uid, routeId, pb.CMD_GET_USER_INFO_REQ, &pb.GetUserInfoReq{
			UidList: []uint64{144, 145},
		})
	})
}

func TestApi(t *testing.T) {
	url := "http://192.168.50.250:8000/mg/callback/ACE/getUserInfo"
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	if err := writer.WriteField("uid", cast.ToString(1000157)); err != nil {
		t.Log(err)
		return
	}
	if err := writer.Close(); err != nil {
		t.Log(err)
		return
	}

	// 1. 检查配置请求
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		t.Log(err)
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		t.Log(err)
		return
	}
	defer resp.Body.Close()

	// 3. 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body) // 尝试读取错误信息
		t.Log(string(body))
		return
	}
	// 4. 解析JSON到protobuf结构体
	rsp := &pb.HttpPlayerInfoRsp{}
	if err := json.NewDecoder(resp.Body).Decode(rsp); err != nil {
		t.Log(err)
		return
	}
	t.Log(rsp)
}

func TestPost(t *testing.T) {
	url := "http://192.168.50.250:8000/mg/callback/ACE/getUserInfo"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("uid", "144")
	err := writer.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	rsp := &pb.HttpPlayerInfoRsp{}
	err = json.NewDecoder(res.Body).Decode(rsp)
	fmt.Println(err, "---->", rsp)
	/*
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(body))
	*/
}
