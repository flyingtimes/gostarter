package Mpool

import (
	"fmt"
	"github.com/go-ini/ini"
	"github.com/mkideal/log"
	"os"
	"runtime"
	"strconv"
	"sync"
)


//任务
type Job struct {
		name string
		nextDispatcher *Dispatcher
}
func (p *Job) Run(pp *Dispatcher){

	fmt.Println("i am working")
  if pp!=nil{
		pp.JobQueue <- Job{
			fmt.Sprintf("工序2-[%s]",p.name),
			nil,
		}
	}


}
//  工人
type Worker struct {
    name string //工人的名字
    // WorkerPool chan JobQueue //对象池
    WorkerPool chan chan Job//对象池
    JobChannel chan Job //通道里面拿
    quit chan bool //
}



func (w *Worker) LoopWork(){
    //开一个新的协程
    go func(){

        for{
            //注册到对象池中,
						log.Info("woker[%s]返回任务池等待任务",w.name)
            w.WorkerPool <-w.JobChannel
            select {
            //接收到了新的任务
            case job :=<- w.JobChannel:
                log.Info("woker[%s]接收到了任务 [%s]",w.name,job.name)
								job.Run(job.nextDispatcher)
								log.Info("woker[%s]完成任务 [%s]",w.name,job.name)
								wg.Done()
            //接收到了任务
            case <-w.quit:
								log.Info("woker[%s]退出。",w.name)
								wg.Done()
                return
            }
        }
    }()
}

func (w Worker) Stop(){
    go func(){
        w.quit <- true
    }()
}

type Dispatcher struct {
                 //WorkerPool chan JobQueue
    name string //调度的名字
    maxWorkers int //获取 调试的大小
    WorkerPool chan chan Job //注册和工人一样的通道
		JobQueue chan Job
}

func (d *Dispatcher) Run(){
    // 开始运行
    for i :=0;i<d.maxWorkers;i++{
        worker := NewWorker(d.WorkerPool,fmt.Sprintf("%s-work-%s",d.name,strconv.Itoa(i)))
        //开始工作
        worker.LoopWork()
    }
    //监控
    go d.LoopGetTask()

}

func (d *Dispatcher) LoopGetTask()  {
    for {
        select {
        case job :=<-d.JobQueue:
            log.Info("调度者[%s][%d]接收到一个工作任务 %s ",d.name, len(d.WorkerPool),job.name)
            // 调度者接收到一个工作任务
            go func (job Job) {
                //从现有的对象池中拿出一个
                jobChannel := <-d.WorkerPool

                jobChannel <- job

            }(job)

        default:

            //fmt.Println("ok!!")
        }

    }
}

// 新建一个工人
func NewWorker(workerPool chan chan Job,name string) Worker{
    log.Info("创建了一个worker:%s \n",name);
    return Worker{
        name:name,//工人的名字
        WorkerPool: workerPool, //工人在哪个对象池里工作,可以理解成部门
        JobChannel:make(chan Job),//工人的任务
        quit:make(chan bool),
    }
}

func NewDispatcher(dname string,maxWorkers int) *Dispatcher {
	  jq := make(chan Job,maxWorkers)
    pool :=make(chan chan Job,maxWorkers)
		log.Info("调度者(%s) 初始化完毕.", dname)
    return &Dispatcher{
        WorkerPool:pool,// 将工人放到一个池中,可以理解成一个部门中
        name:dname,//调度者的名字
        maxWorkers:maxWorkers,//这个调度者有好多个工人
				JobQueue:jq,
    }
}


/*
func (p *Mpool) initialize(name string,workers int ) chan Job {
	dispatch := NewDispatcher(name,workers,JobQueue)
	dispatch.Run()
	log.Info("调度者(%s) 初始化完毕.", name)
	return JobQueue
}
*/
var wg sync.WaitGroup
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

	wg.Add(10000*2)
	dispacher1:= NewDispatcher("01",3)
	dispacher2 := NewDispatcher("02",3)
	dispacher1.Run()
	dispacher2.Run()
	for i:=0;i<10000;i++{

        dispacher1.JobQueue <- Job{
            fmt.Sprintf("工序1-[%s]",strconv.Itoa(i)),
						dispacher2,
        }


    }
   wg.Wait()
    //close(dispacher1.JobQueue)
	  //close(dispacher2.JobQueue)
	log.Info("Main exit normally.")
}
