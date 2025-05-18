func HMSet({{if .Key.Values}} {{.Key.Values.Arg}}, {{end}} values map[string]*pb.{{.Name}}) error {
	if len(values) <= 0 {
		return nil
	}
	args := []interface{}{}
	for field, val := range values {
		buf, err := val.Marshal()
		if err != nil {
			return err
		}
		args = append(args, field, buf)
	}
	return redis.HMSet(DBNAME, GetRedisKey({{.Key.Values.Val ""}}), args...)
}
