{{$total :=  .Key.Values.Join .Field.Values}}
func HExists({{$total.Arg}}) (int, error) {
	return redis.HExists(DBNAME, GetRedisKey({{.Key.Values.Val ""}}), GetRedisField({{.Field.Values.Val ""}}))
}
