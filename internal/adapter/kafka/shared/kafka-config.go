package shared

const (
	CreateOperation = "create"
	ChangeOperation = "put"
	PatchOperation  = "patch"
	DeleteOperation = "delete"
)

type KafkaConfig struct {
	TaskTopic string
	UserTopic string
	Addr      []string
}

func NewKafkaConfig() *KafkaConfig {
	return &KafkaConfig{
		TaskTopic: "tasks",
		UserTopic: "users",
		Addr:      []string{"localhost:9092"},
	}
}

type MsgOperation struct {
	Operation string
}

//type TaskResponse struct {
//	Id       int       `json:"id"`
//	Title    string    `json:"title"`
//	Text     string    `json:"text"`
//	Priority int       `json:"priority"`
//	Expires  time.Time `json:"expires"`
//}

type TaskBasic struct {
	Title      string `json:"title"`
	Text       string `json:"text"`
	Priority   int    `json:"priority"`
	ExpireDays int    `json:"expire_days"`
}

type TaskCreate struct {
	TaskBasic
	UserId string `json:"user_id"`
}

type TaskDelete struct {
	ID     string `json:"id"`
	Role   string `json:"role"`
	UserId string `json:"user_id"`
}

type TaskPatch struct {
	TaskBasic
	TaskDelete
}

//type UserResponse struct {
//	Id        int       `json:"id"`
//	Username  string    `json:"username"`
//	Role      string    `json:"role"`
//	Email     string    `json:"email"`
//	CreatedAt time.Time `json:"created_at"`
//}

type UserBasic struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserDelete struct {
	ID string `json:"id"`
}

type UserPatch struct {
	UserBasic
	ID string `json:"id"`
}
