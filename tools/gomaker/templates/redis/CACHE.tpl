{{if .IsString}}
{{if .Key.Values.HasUID}}
func GetCacheKey({{.Key.Values.Arg}}) string {
	return fmt.Sprintf("{{.Name}}:{{.Key.Values.Format}}", {{.Key.Values.Val ""}})
}

func GetByCache({{.Key.Values.Arg}}) (ret *pb.{{.Name}}, exist bool, err error) {
	if sess, ok := framework.GetUserSession(uid); ok && sess != nil {
		if ret, exist = sess.GetValue(GetCacheKey({{.Key.Values.Val ""}})).(*pb.{{.Name}}); exist && ret != nil {
			return
		}
	}
	return Get({{.Key.Values.Val ""}})
}
{{end}}
{{else if .IsHash}}
{{if .Key.Values.HasUID}} {{$total :=  .Key.Values.Join .Field.Values}}
func GetCacheKey({{$total.Arg}}) string {
	return fmt.Sprintf("{{.Name}}:{{$total.Format}}", {{$total.Val ""}})
}

func GetByCache({{$total.Arg}}) (result *pb.{{.Name}}, exist bool, err error) {
	if sess, ok := framework.GetUserSession(uid); ok && sess != nil {
		if result, exist = sess.GetValue(GetCacheKey({{$total.Val ""}})).(*pb.{{.Name}}); exist && result != nil {
			return
		}
	}
	result, exist, err = HGet({{$total.Val ""}})
	return
}
{{end}}
{{end}}
