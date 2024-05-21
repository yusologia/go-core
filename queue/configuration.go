package queue

type JobConf struct {
	Context     interface{}
	JobFunc     interface{}
	QueueName   string
	JobName     string
	Priority    uint
	Concurrency uint
}
