package main

import (
	"fmt"
	"github.com/go-ini/ini"
	"github.com/mkideal/log"
	"os"
	"runtime"
)

//声明成游戏
type Payload struct {
    name string
}

//任务
type Job struct {
    Payload Payload
}

type Mpool struct {
	name string
}

func (p *Mpool) initialize() chan Job {
	maxWorkers := runtime.NumCPU()
	maxQueue := maxWorkers
	JobQueue := make(chan Job,maxWorkers)

	log.Info("Mpool (%s) initialized.", p.name)
	return JobQueue
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
	pool := Mpool{"pool_default"}
	job_queue := pool.initialize()
	log.Info("Main exit normally.")
}
