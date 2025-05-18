func HKeys({{.Key.Values.Arg}}) ([]string, error) {
	return redis.HKeys(DBNAME, GetRedisKey({{.Key.Values.Val ""}}))
}
