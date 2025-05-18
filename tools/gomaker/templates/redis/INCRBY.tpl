func IncrBy({{if .Key.Values}} {{.Key.Values.Arg}}, {{end}} incr int) (int, error) {
	return redis.IncrBy(DBNAME, GetRedisKey({{.Key.Values.Val ""}}), incr)
}
