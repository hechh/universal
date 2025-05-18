func MSet(vals map[string]*pb.{{.Name}}) (err error) {
	args := []interface{}{}
	for key, val := range vals {
		if buf, err := val.Marshal(); err != nil {
			return err
		} else {
			args = append(args, key, buf)
		}
	}
	return redis.MSet(DBNAME, args...)
}
