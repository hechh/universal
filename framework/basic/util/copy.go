package util

import (
	"regexp"
	"universal/common/pb"
	"universal/framework/basic/uerror"
)

func CopyHead(head *pb.Head) *pb.Head {
	new := *head
	newRouteInfo := *head.Table
	new.Table = &newRouteInfo
	return &new
}

func Filter(pattern string, vals ...string) (rets []string, err error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, uerror.NewUError(1, -1, "%s, %v, %v", pattern, vals, err)
	}
	for _, val := range vals {
		if re.MatchString(val) {
			continue
		}
		rets = append(rets, val)
	}
	return
}
