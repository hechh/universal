package test

/*
import (
	"testing"
	"time"
	"hego/common/global"
	"hego/framework/cluster/internal/discovery"
)

var (
	cfg *global.Config
)

func TestMain(m *testing.M) {
	cfg, _ = global.Load("../../../env/config", "gate")
	m.Run()
}

func TestPut(t *testing.T) {
	cli, err := discovery.NewEtcdClient(cfg.Etcd.Endpoints...)
	if err != nil {
		t.Log(err)
		return
	}

	cli.Put("/root/key01", "value01")
	cli.Put("/root/key02", "value02")
	cli.Put("/root/key03", "value03")

	err = cli.Walk("/root", func(key, val string) {
		t.Log(key, "----------", val)
		cli.Delete(key)
	})
	if err != nil {
		t.Log(err)
	}
}

func TestKeepAlive(t *testing.T) {
	cli, err := discovery.NewEtcdClient(cfg.Etcd.Endpoints...)
	if err != nil {
		t.Log(err)
		return
	}
	cli.KeepAlive("/root/key04", "value04", 6)
	time.Sleep(8 * time.Second)
	err = cli.Walk("/root", func(key, val string) {
		t.Log(key, "----------", val)
	})
	if err != nil {
		t.Log(err)
	}
	cli.Close()
}

func TestWatch(t *testing.T) {
	cli, err := discovery.NewEtcdClient(cfg.Etcd.Endpoints...)
	if err != nil {
		t.Log(err)
		return
	}
	if err := cli.KeepAlive("/root/key05", "value05", 6); err != nil {
		t.Log(err)
		return
	}
	addF := func(key, val string) {
		t.Log(key, "-----add-----", val)
	}
	delF := func(key, val string) {
		t.Log(key, "-----del-----", val)
	}
	if err = cli.Watch("/root", addF, delF); err != nil {
		t.Log(err)
		return
	}
	cli.Close()
	time.Sleep(8 * time.Second)
}
*/
