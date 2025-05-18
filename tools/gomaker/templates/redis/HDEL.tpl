{{$total :=  .Key.Values.Join .Field.Values}}
func HDel({{$total.Arg}}) (int, error) {
	return redis.HDel(DBNAME, GetRedisKey({{.Key.Values.Val ""}}), GetRedisField({{.Field.Values.Val ""}}))
}
