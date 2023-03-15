package simpleworkerpool

type Context struct {
	Event  EventType
	Handle WorkerHandle
	Args   []interface{}
}
