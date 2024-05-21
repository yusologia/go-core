package queue

import (
	"fmt"
	"github.com/gocraft/work"
	"gorm.io/gorm/utils"
	"log"
	"os"
	"os/signal"
	"strings"
)

type Queue struct {
	Names string
}

func (q Queue) Work(workers []JobConf) {
	var names []string

	if len(q.Names) > 0 {
		q.Names = strings.ReplaceAll(q.Names, " ", "")
		names = strings.Split(q.Names, ",")
	}

	RegisterRedis()

	var pools []*work.WorkerPool

	for _, worker := range workers {
		if len(names) > 0 {
			if !utils.Contains(names, worker.QueueName) {
				continue
			}
		}

		pool := work.NewWorkerPool(worker.Context, worker.Concurrency, worker.QueueName, RedisPool)
		pool.JobWithOptions(worker.JobName, work.JobOptions{
			Priority: worker.Priority,
		}, worker.JobFunc)

		pool.Start()

		pools = append(pools, pool)

		log.Println(fmt.Sprintf("%s:%s", worker.QueueName, worker.JobName))
	}

	// Wait for a signal to quit:
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan

	for _, pool := range pools {
		pool.Stop()
	}

	fmt.Println("All worker is done!")
}
