package httpkit

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	_ "net/http/pprof"
	"path/filepath"
	"reflect"
	"runtime/debug"
	"sort"
	"strings"
	"sync"
	"time"
	"universal/framework/async"
	"universal/framework/basic"
	"universal/framework/handler"
	"universal/library/mlog"
	"universal/tools/client/domain"
	"universal/tools/client/internal/player"

	"github.com/golang/protobuf/proto"
	"github.com/spf13/cast"
)

const (
	editHtml = `<h1>客户端请求JSON数据</h1>
        <form action="%s" method="post">
            <textarea name="json" rows="30" cols="150">{{.}}</textarea><br>
            <button type="submit">发送请求</button>
        </form>`
)

var (
	fatalNotify = func(err interface{}) {
		mlog.Fatal("error: %v, stack: %s", err, string(debug.Stack()))
	}
	jsons = make(map[string]string)
)

func registerJson(name string, reqJson interface{}) {
	if val := reflect.ValueOf(reqJson).Elem().Field(0); val.IsNil() {
		val.Set(reflect.ValueOf(&IPacket{}))
	}
	buf, _ := json.Marshal(reqJson)
	jsons[name] = string(buf)
}

func Init() {
	http.HandleFunc("/api", api)
	http.HandleFunc("/login", login)
	http.HandleFunc("/online", online)
	http.HandleFunc("/logout", logout)

	basic.SafeGo(fatalNotify, func() {
		if err := http.ListenAndServe("localhost:22345", nil); err != nil {
			mlog.Info("start http server is failed, error: %v", err)
		}
		mlog.Info("start http server is success")
	})
}

func logout(w http.ResponseWriter, r *http.Request) {
	player.Close()
	fmt.Fprintf(w, "在线人数: %d\n", player.GetOnline())
}

func online(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "在线人数: %d\n", player.GetOnline())
}

func api(w http.ResponseWriter, r *http.Request) {
	reqs := []string{}
	handler.Walk(func(api *handler.ApiInfo) bool {
		if strings.HasSuffix(api.GetReqName(), "Request") {
			reqs = append(reqs, api.GetReqName())
		}
		return true
	})
	sort.Slice(reqs, func(i, j int) bool {
		return strings.Compare(reqs[i], reqs[j]) <= 0
	})
	for i, name := range reqs {
		fmt.Fprintf(w, `<a href="%s/%s"> %3d. %s</a><br>`, r.URL.Path, name, i+1, name)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	wg := sync.WaitGroup{}
	cache := async.NewQueue()
	result := func(api *domain.Result) {
		defer basic.SafeRecover(fatalNotify, wg.Done)
		cache.Push(api)
	}
	// 解析请求参数
	r.ParseForm()
	if suid := r.Form.Get("uid"); len(suid) > 0 {
		uid := cast.ToUint64(suid)
		if pl := player.GetPlayer(uid); pl == nil {
			wg.Add(1)
			player.Login(uid, result)
		}
	} else {
		begin := r.Form.Get("begin")
		end := r.Form.Get("end")
		for uid := cast.ToUint64(begin); uid <= cast.ToUint64(end); uid++ {
			if pl := player.GetPlayer(uid); pl == nil {
				wg.Add(1)
				player.Login(uid, result)
			}
		}
	}
	// 超时关闭
	deadline(30*time.Second, &wg)
	now := basic.GetNowUnixSecond()
	basic.SafeRecover(fatalNotify, wg.Wait)
	response(w, cache, "LoginRequest", basic.GetNowUnixSecond()-now)
}

func handle(w http.ResponseWriter, r *http.Request) {
	opr := filepath.Base(r.URL.Path)
	r.ParseForm()
	jsonData := r.Form.Get("json")
	// 动态网页生成
	if len(jsonData) <= 0 {
		tmpl := template.Must(template.New("index").Parse(fmt.Sprintf(editHtml, r.URL.Path)))
		if err := tmpl.Execute(w, jsons[opr]); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	// 设置回调接口
	cache := async.NewQueue()
	wg := sync.WaitGroup{}
	result := func(api *domain.Result) {
		defer basic.SafeRecover(fatalNotify, wg.Done)
		cache.Push(api)
	}
	// 发送请求
	if err := player.Send(&wg, opr, jsonData, result); err != nil {
		cache.Push(fmt.Sprintf("error: %v\n", err))
	}
	// 超时关闭
	deadline(15*time.Second, &wg)
	now := basic.GetNowUnixSecond()
	basic.SafeRecover(fatalNotify, wg.Wait)
	response(w, cache, opr, basic.GetNowUnixMilli()-now)
}

func deadline(dur time.Duration, wg *sync.WaitGroup) {
	basic.SafeGo(nil, func() {
		<-time.NewTimer(dur).C
		for i := 0; i < 10000; i++ {
			wg.Done()
		}
	})
}

type Result struct {
	Name     string
	Cost     int64
	Total    int64
	Success  int64
	Failure  int64
	Response []proto.Message
}

func response(w http.ResponseWriter, cache *async.Queue, name string, cost int64) {
	ret := &Result{
		Name:  name,
		Cost:  cost,
		Total: cache.GetCount(),
	}
	for str := cache.Pop(); str != nil; str = cache.Pop() {
		vv := str.(*domain.Result)
		head := vv.Response.(handler.IHead).GetPacketHead()
		if vv.Error != nil || head.Code != 0 {
			ret.Failure++
		} else {
			ret.Success++
		}
		head.Id = vv.UID
		ret.Response = append(ret.Response, vv.Response)
	}
	buf, _ := json.Marshal(ret)
	w.Write(buf)
}
