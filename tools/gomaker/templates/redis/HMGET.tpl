func HMGet({{if .Key.Values}} {{.Key.Values.Arg}}, {{end}} fields ...interface{}) (result map[string]*pb.{{.Name}}, err error) {
	if len(fields) <= 0 {
		return
	}
	bufs, err := redis.HMGet(DBNAME, GetRedisKey({{.Key.Values.Val ""}}), fields...)
	if err != nil {
		return nil, err
	}
	result = make(map[string]*pb.{{.Name}})
	for i, buf := range bufs {
		elem := new{{.Name}}()
		if err := elem.Unmarshal(buf); err != nil {
			return nil, err
		}
		result[fields[i].(string)] = elem
	}
	return
}
