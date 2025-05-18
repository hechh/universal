func MGet(keys ...string) (result map[string]*pb.{{.Name}}, err error) {
	if len(keys) <= 0 {
		return
	}
	bufs, err := redis.MGet(DBNAME, keys...)
	if err != nil {
		return nil, err
	}
	result = make(map[string]*pb.{{.Name}})
	for i, buf := range bufs {
		elem := new{{.Name}}()
		if err := elem.Unmarshal(buf); err != nil {
			return nil, err
		}
		result[keys[i]] = elem
	}
	return
}
