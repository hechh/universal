func TTL({{.Key.Values.Arg}}) (int, error) {
	return redis.TTL(DBNAME, GetRedisKey({{.Key.Values.Val ""}}))
}
