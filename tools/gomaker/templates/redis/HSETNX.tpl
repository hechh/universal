{{$total :=  .Key.Values.Join .Field.Values}}
func HSetNx({{if $total}} {{$total.Arg}}, {{end}} data *pb.{{.Name}}) (ok bool, err error) {
	buf, err := data.Marshal()
	if err != nil {
		return false, err
	}
	if ok, err = redis.HSetNX(DBNAME, GetRedisKey({{.Key.Values.Val ""}}), GetRedisField({{.Field.Values.Val ""}}), buf); err == nil {
		manager.CallSets(TABLE_NAME, data {{if $total}}, {{$total.Val ""}} {{end}})
	}
	return
}
