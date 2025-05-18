func Del({{.Key.Values.Arg}}) error {
	return redis.Del(DBNAME, GetRedisKey({{.Key.Values.Val ""}}))
}
