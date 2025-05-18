func HScan({{if .Key.Values}} {{.Key.Values.Arg}}, {{end}} cursor int64, match string, count int) (pos int64, result map[string]*pb.{{.Name}}, err error) {
	pos, rets, err := redis.HScan(DBNAME, GetRedisKey({{.Key.Values.Val ""}}), cursor, match, count)
	if err != nil {
		return 0, nil, err
	}
	result = make(map[string]*pb.{{.Name}})
	for field, val := range rets {
		elem := new{{.Name}}()
		if err := elem.Unmarshal([]byte(val)); err != nil {
			return 0, nil, err
		}
		result[field] = elem
	}
	return
}
