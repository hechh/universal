func Exists({{.Key.Values.Arg}}) (bool, error) {
	return redis.Exists(DBNAME, GetRedisKey({{.Key.Values.Val ""}}))
}
