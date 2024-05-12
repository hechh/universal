package manager

import (
	"time"
	"universal/common/pb"
	"universal/framework/fbasic"

	"gopkg.in/mgo.v2"
)

var mongoPool = make(map[string]*mgo.Session)

func RegisterMongoPool(dbname, user, passwd string, addrs ...string) (err error) {
	if _, ok := mongoPool[dbname]; ok {
		return fbasic.NewUError(1, pb.ErrorCode_HasRegistered, dbname, addrs)
	}
	dialInfo := &mgo.DialInfo{
		Addrs:     addrs,
		Timeout:   time.Duration(60) * time.Second,
		Source:    dbname,
		Username:  user,
		Password:  passwd,
		PoolLimit: 200,
	}
	var sess *mgo.Session
	if sess, err = mgo.DialWithInfo(dialInfo); err != nil {
		return
	}
	// 判断是否可以ping
	if err = sess.Ping(); err != nil {
		return
	}
	mongoPool[dbname] = sess
	return
}

func GetMongo(dbName string) (conn *mgo.Session, err error) {
	sess, ok := mongoPool[dbName]
	if !ok {
		return nil, fbasic.NewUError(1, pb.ErrorCode_NotSupported, dbName)
	}
	conn = sess.Clone()
	return
}

func PutMongo(conn *mgo.Session) {
	if conn != nil {
		conn.Close()
	}
}
