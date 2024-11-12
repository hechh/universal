package typespec

type Stock struct {
	Code      string    // 股票代码
	Name      string    // 股票名字
	Reason    string    // 涨停原因
	Themes    []string  // 所属概念
	BeginTime string    // 首次涨停时间
	EndTime   string    // 最终涨停时间
	BeginUnix int64     // 首次涨停时间
	EndUnix   int64     // 最终涨停时间
	Continue  int       // 连续涨停天数
	Total     string    // 几天几板
	Time      *TimeInfo // 时间
}

type Hotspot struct {
	Name        string           // 热点名称
	MaxContinue int              // 最高连板
	Members     map[int][]*Stock // 连板数据
	List        []int            // 连板高度数列
}

type TimeInfo struct {
	Date  int64 // 时间
	Year  int32 // 年
	Month int32 // 月
	Day   int32 // 日
}
