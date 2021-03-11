package setting

func newTaskService() {
	taskSec := Cfg.Section("task")
	queueTaskSec := Cfg.Section("queue.task")
	switch taskSec.Key("QUEUE_TYPE").MustString(ChannelQueueType) {
	case ChannelQueueType:
		queueTaskSec.Key("TYPE").MustString("persistable-channel")
		queueTaskSec.Key("CONN_STR").MustString(taskSec.Key("QUEUE_CONN_STR").MustString(""))
	case RedisQueueType:
		queueTaskSec.Key("TYPE").MustString("redis")
		queueTaskSec.Key("CONN_STR").MustString(taskSec.Key("QUEUE_CONN_STR").MustString("addrs=127.0.0.1:6379 db=0"))
	}
	queueTaskSec.Key("LENGTH").MustInt(taskSec.Key("QUEUE_LENGTH").MustInt(1000))
}
