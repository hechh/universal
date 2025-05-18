func SetEx({{if .Key.Values}} {{.Key.Values.Arg}}, {{end}} ttl int, data *pb.{{.Name}}) (err error) {
	buf, err := data.Marshal()
	if err != nil {
		return err
	}
	if err = redis.SetEX(DBNAME, GetRedisKey({{.Key.Values.Val ""}}), ttl, buf); err == nil {
		manager.CallSets(TABLE_NAME, data {{if .Key}}, {{.Key.Values.Val ""}} {{end}})
	}
	return
}
