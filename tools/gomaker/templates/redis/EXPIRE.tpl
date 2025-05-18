func Expire({{if .Key.Values}} {{.Key.Values.Arg}}, {{end}} ttl int) error {
	return redis.Expire(DBNAME, GetRedisKey({{.Key.Values.Val ""}}), ttl)
}
