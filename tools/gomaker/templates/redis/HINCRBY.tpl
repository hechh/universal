{{$total :=  .Key.Values.Join .Field.Values}}
func HIncrBy({{if $total}} {{$total.Arg}}, {{end}} incr int) (int, error) {
	return redis.HIncrBy(DBNAME, GetRedisKey({{.Key.Values.Val ""}}), GetRedisField({{.Field.Values.Val ""}}), incr)
}
