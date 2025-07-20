package player

import (
	"poker_server/common/config/repository/php_config"
	"poker_server/common/pb"
	"poker_server/framework"
	"poker_server/library/httplib"
	"poker_server/library/uerror"

	"github.com/spf13/cast"
)

func PlayerInfoRequest(uid uint64, rsp *pb.HttpPlayerInfoRsp) error {
	cfg := php_config.MGetEnvTypeName(framework.GetEnvType(), "UserInfoUrl")
	if cfg == nil {
		return uerror.New(1, pb.ErrorCode_CONFIG_NOT_FOUND, "Php配置不存在")
	}
	params := map[string]interface{}{"uid": cast.ToString(uid)}
	if err := httplib.POST(cfg.Url, params, rsp); err != nil {
		return err
	}
	if rsp.RespMsg.Code != 100 {
		return uerror.New(1, pb.ErrorCode_REQUEST_FAIELD, "uid:%d 请求失败: %s", uid, rsp.RespMsg.Message)
	}
	return nil
}
