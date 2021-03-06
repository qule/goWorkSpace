package zk_util

import (
	"github.com/samuel/go-zookeeper/zk"
	"go_spider-master/core/common/config"
	"go_spider-master/core/common/util"
	"github.com/bitly/simplejson"
	"fmt"
	"time"
	"strings"
//	"reflect"
	"strconv"
)

type Zk_util struct {
	conn			*zk.Conn
	flag			int32
	ephemeralflag	int32
	acl				[]zk.ACL
}

func (this *Zk_util) Connect() *Zk_util {
	zkHosts, err := config.GetSpecConfig(config.CONF_FILE_PATH, "zookeeper", "hosts")
	util.Panicable(err)
	zks := strings.Split(zkHosts, ",")
	this.conn, _, err = zk.Connect(zks, time.Second)
	util.Panicable(err)

	this.flag = int32(0)
	this.ephemeralflag = int32(zk.FlagEphemeral)
	this.acl = zk.WorldACL(zk.PermAll)

	return this
}

func (this *Zk_util) createPathIfNotExist(dest string, ephemeral bool) *Zk_util {
	if dest == "/" {
		util.Printf("reach root node!")
		return nil
	}

	targetIndex := strings.LastIndex(dest, "/")
	parent := dest[0:targetIndex]
//	target := dest[targetIndex+1:len(dest)]

	parentExists, _, _ := this.conn.Exists(parent)
	if parentExists == false {
		util.Printf(fmt.Sprintf("parent: %v NOT exists, step into recursion", parent))
		this.createPathIfNotExist(parent[0:len(parent)], ephemeral)
	}
	destExists, _, _ := this.conn.Exists(dest)
	if destExists == false {
		util.Printf(fmt.Sprintf("cur: %v NOT exists, creating...", dest))

		var flag int32
		if ephemeral == false {flag = this.flag} else {flag = this.ephemeralflag}

		_, err := this.conn.Create(dest, []byte(""), flag, this.acl)
		if nil != err {
			panic(err)
		}
	}

	return this
}

func (this *Zk_util) CreatePermanentPathIfNotExist(dest string) *Zk_util {
	return this.createPathIfNotExist(dest, false)
}

func (this *Zk_util) CreateEphemeralPathIfNotExist(dest string) *Zk_util {
	return this.createPathIfNotExist(dest, true)
}

func (this *Zk_util) SetPathData(dest string, value string) *Zk_util {
	_, stat, err := this.conn.Get(dest)

//	destExists, _, _ := this.conn.Exists(dest)
	switch {
	case dest == "/":
		util.Printf("You cannot set value to root node!")
		return nil
	case err != nil:
		util.Printf(fmt.Sprintf("dest: %v MAY NOT exists", dest))
		return nil
	}

	_, err = this.conn.Set(dest, []byte(value), stat.Version)
	util.Panicable(err)

	return this
}

// data format: json
func (this *Zk_util) GetPathData(dest string) []byte {
	destExists, _, _ := this.conn.Exists(dest)
	if destExists == false {
		util.Printf(fmt.Sprintf("dest: %v NOT exists", dest))
		return nil
	}

	dataRes, _, err := this.conn.Get(dest)
	util.Panicable(err)

	return dataRes
}

func (this *Zk_util) DeletePath(dest string) *Zk_util {
	err := this.conn.Delete(dest, -1)
	if err != nil {
		util.Printf(fmt.Sprintf("dest: %v MAY NOT exists", dest))
	}

	return this
}


// create path of following nodes
// 1. master
func (this *Zk_util) CreateMaster() *Zk_util {
	return this.CreatePermanentPathIfNotExist(config.ZK_CRAWLER_MASTER_PATH)
}
func (this *Zk_util) CreateMasterHeartbeat() *Zk_util {
	return this.CreatePermanentPathIfNotExist(config.ZK_CRAWLER_MASTER_HEARTBEAT_PATH)
}
// 2. scheduler
func (this *Zk_util) CreateScheduler() *Zk_util {
	return this.CreatePermanentPathIfNotExist(config.ZK_CRAWLER_SCHEDULER_PATH)
}
func (this *Zk_util) CreateSchedulerHeartbeat() *Zk_util {
	return this.CreatePermanentPathIfNotExist(config.ZK_CRAWLER_SCHEDULER_HEARTBEAT_PATH)
}
// 3. dp
func (this *Zk_util) CreateDPStat() *Zk_util {
	return this.CreatePermanentPathIfNotExist(config.ZK_CRAWLER_DP_PATH)
}
func (this *Zk_util) CreateDP(id int) *Zk_util {
	path := config.ZK_CRAWLER_DP_PATH + "/" + strconv.Itoa(id)
	return this.CreatePermanentPathIfNotExist(path)
}
func (this *Zk_util) CreateDPHeartbeat(id int) *Zk_util {
	path := config.ZK_CRAWLER_DP_PATH + "/" + strconv.Itoa(id) + "/heartbeat"
	return this.CreatePermanentPathIfNotExist(path)
}
// 4. sink
func (this *Zk_util) CreateSink() *Zk_util {
	return this.CreatePermanentPathIfNotExist(config.ZK_CRAWLER_SINK_PATH)
}
func (this *Zk_util) CreateSinkHeartbeat() *Zk_util {
	return this.CreatePermanentPathIfNotExist(config.ZK_CRAWLER_SINK_HEARTBEAT_PATH)
}

// get data of following nodes
// 1. master
func (this *Zk_util) GetMaster() []byte {
	return this.GetPathData(config.ZK_CRAWLER_MASTER_PATH)
}
func (this *Zk_util) GetMasterHeartbeat() []byte {
	return this.GetPathData(config.ZK_CRAWLER_MASTER_HEARTBEAT_PATH)
}
// 2. scheduler
func (this *Zk_util) GetScheduler() []byte {
	return this.GetPathData(config.ZK_CRAWLER_SCHEDULER_PATH)
}
func (this *Zk_util) GetSchedulerHeartbeat() []byte {
	return this.GetPathData(config.ZK_CRAWLER_SCHEDULER_HEARTBEAT_PATH)
}
// 3. dp
func (this *Zk_util) GetDPStat() []byte {
	return this.GetPathData(config.ZK_CRAWLER_DP_PATH)
}
func (this *Zk_util) GetDP(id int) []byte {
	path := config.ZK_CRAWLER_DP_PATH + "/" + strconv.Itoa(id)
	return this.GetPathData(path)
}
func (this *Zk_util) GetDPHeartbeat(id int) []byte {
	path := config.ZK_CRAWLER_DP_PATH + "/" + strconv.Itoa(id) + "/heartbeat"
	return this.GetPathData(path)
}
// 4. sink
func (this *Zk_util) GetSink() []byte {
	return this.GetPathData(config.ZK_CRAWLER_SINK_PATH)
}
func (this *Zk_util) GetSinkHeartbeat() []byte {
	return this.GetPathData(config.ZK_CRAWLER_SINK_HEARTBEAT_PATH)
}

// set data of following nodes
// 1. master
func (this *Zk_util) SetMaster(data []byte) *Zk_util {
	return this.SetPathData(config.ZK_CRAWLER_MASTER_PATH, string(data))
}
func (this *Zk_util) UpMasterHeartBeat(ts int64) *Zk_util {
	return this.SetPathData(config.ZK_CRAWLER_MASTER_HEARTBEAT_PATH, strconv.Itoa(int(ts)))
}
// 2. scheduler
func (this *Zk_util) SetScheduler(data []byte) *Zk_util {
	return this.SetPathData(config.ZK_CRAWLER_SCHEDULER_PATH, string(data))
}
func (this *Zk_util) UpSchedulerHeartBeat(ts int64) *Zk_util {
	return this.SetPathData(config.ZK_CRAWLER_SCHEDULER_HEARTBEAT_PATH, strconv.Itoa(int(ts)))
}
// 3. dp
func (this *Zk_util) SetDPStat(data []byte) *Zk_util {
	return this.SetPathData(config.ZK_CRAWLER_DP_PATH, string(data))
}
func (this *Zk_util) SetDP(id int, data []byte) *Zk_util {
	path := config.ZK_CRAWLER_DP_PATH + "/" + strconv.Itoa(id)
	return this.SetPathData(path, string(data))
}
func (this *Zk_util) UpDPHeartBeat(id int, ts int64) *Zk_util {
	path := config.ZK_CRAWLER_DP_PATH + "/" + strconv.Itoa(id) + "/heartbeat"
	return this.SetPathData(path, strconv.Itoa(int(ts)))
}
// 4. sink
func (this *Zk_util) SetSink(data []byte) *Zk_util {
	return this.SetPathData(config.ZK_CRAWLER_SINK_PATH, string(data))
}
func (this *Zk_util) UpSinkHeartBeat(ts int64) *Zk_util {
	return this.SetPathData(config.ZK_CRAWLER_SINK_HEARTBEAT_PATH, strconv.Itoa(int(ts)))
}




func main() {
	zkUtil := Zk_util{}
	zkUtil.Connect()

	zkUtil.CreatePermanentPathIfNotExist(config.ZK_CRAWLER_DP_PATH)

	//json example string for testing
	js, err := simplejson.NewJson([]byte("{}"))
	util.Panicable(err)
	js.Set("rootpath1", "rootpathvalue1")
	js.SetPath([]string{"path1", "path2"}, "pathvalue")
	js.SetPath([]string{"path1", "path3", "path4"}, "path3value")
	testJson, err := js.Map()

	config.ShowParseJsonMap(testJson)

//	testJsonStr := fmt.Sprintf("%v", testJson)
	testJsonStr, _ := js.MarshalJSON()

	zkUtil.SetPathData(config.ZK_CRAWLER_DP_PATH, string(testJsonStr))

	fmt.Println(string(zkUtil.GetPathData(config.ZK_CRAWLER_DP_PATH)))

	zkUtil.DeletePath(config.ZK_CRAWLER_DP_PATH)
}
