package Mpool

import (
	"fmt"
	"github.com/flyingtimes/gostarter"
	"github.com/go-ini/ini"
	"github.com/mkideal/log"
	"os"
	"runtime"
	"strconv"
)

//任务
type Job struct {
	name           string
	nextDispatcher *Mpool.Dispatcher
}

func (j Job) Run(pp *Mpool.Dispatcher) {
	fmt.Println("i am working")
}
func (j Job) GetName() string {
	return j.name
}
func (j Job) GetNextDispatcher() *Mpool.Dispatcher {
	return j.nextDispatcher
}

func main() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		fmt.Println("找不到配置文件：", err)
		os.Exit(1)
	}
	section, err := cfg.GetSection("main")
	if err != nil {
		fmt.Println("找不到main的配置信息：", err)
		os.Exit(1)
	}
	key, err := section.GetKey("logfile")
	if err != nil {
		fmt.Println("找不到logfile的配置信息：", err)
		os.Exit(1)
	}
	log_file_name := key.String()
	defer log.Uninit(log.InitFile(log_file_name))
	log.Info("Main started.")

	runtime.GOMAXPROCS(4)

	dispacher1 := Mpool.NewDispatcher("01", 3)
	dispacher2 := Mpool.NewDispatcher("02", 3)
	dispacher1.Run()
	dispacher2.Run()
	for i := 0; i < 10000; i++ {
		dispacher1.JobQueue <- Job{
			fmt.Sprintf("工序1-[%s]", strconv.Itoa(i)),
			dispacher2,
		}

	}
	//close(dispacher1.JobQueue)
	//close(dispacher2.JobQueue)
	log.Info("Main exit normally.")
}
