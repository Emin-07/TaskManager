package domain

const (
	CreateOperation = "create"
	ChangeOperation = "put"
	PatchOperation  = "patch"
	DeleteOperation = "delete"
)

const (
	TopicTasks = "tasks"
	TopicUsers = "users"
)

type Message struct {
	Key   []byte
	Val   []byte
	Topic string
}
