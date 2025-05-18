func HGetAll({{.Key.Values.Arg}}) (result map[string]*pb.{{.Name}}, err error) {
	vals, err := redis.HGetAll(DBNAME, GetRedisKey({{.Key.Values.Val ""}}))
	if err != nil {
		return nil, err
	}

	result = make(map[string]*pb.{{.Name}})
	for field, val := range vals {
		elem := new{{.Name}}()
		if err := elem.Unmarshal([]byte(val)); err != nil {
			return nil, err
		}
		result[field] = elem
	}
	return
}
