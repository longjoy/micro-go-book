package setup

import (
	"fmt"
	conf "github.com/longjoy/micro-go-book/ch13-seckill/pkg/config"
	"github.com/samuel/go-zookeeper/zk"
	"log"
	"testing"
	"time"
)

func ZkStateStringFormat(s *zk.Stat) string {
	return fmt.Sprintf("Czxid:%d\nMzxid: %d\nCtime: %d\nMtime: %d\nVersion: %d\nCversion: %d\nAversion: %d\nEphemeralOwner: %d\nDataLength: %d\nNumChildren: %d\nPzxid: %d\n",
		s.Czxid, s.Mzxid, s.Ctime, s.Mtime, s.Version, s.Cversion, s.Aversion, s.EphemeralOwner, s.DataLength, s.NumChildren, s.Pzxid)
}

func waitSecProductEvent(event zk.Event) {
	log.Print(">>>>>>>>>>>>>>>>>>>")
	log.Println("path:", event.Path)
	log.Println("type:", event.Type.String())
	log.Println("state:", event.State.String())
	log.Println("<<<<<<<<<<<<<<<<<<<")
	if event.Path == conf.Zk.SecProductKey {

	}
}

func TestInitZK(t *testing.T) {
	var hosts = []string{"39.98.179.73:2181"}
	option := zk.WithEventCallback(waitSecProductEvent)
	conn, _, err := zk.Connect(hosts, time.Second*5, option)
	if err != nil {
		fmt.Println(err)
		return
	}

	var path = "/zk_test_g1o"
	var data = []byte("hello")
	var flags int32 = 0
	// permission
	var acls = zk.WorldACL(zk.PermAll)

	// create
	p, err_create := conn.Create(path, data, flags, acls)
	if err_create != nil {
		fmt.Println(err_create)
		return
	}
	fmt.Println("created:", p)

	// get
	v, s, err := conn.Get(path)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("value of path[%s]=[%s].\n", path, v)
	fmt.Printf("state:\n")
	fmt.Printf("%s\n", ZkStateStringFormat(s))

	defer conn.Close()
}
