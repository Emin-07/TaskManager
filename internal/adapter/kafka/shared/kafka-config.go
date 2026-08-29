package shared

//TODO Delete kafkaconfig, and make receiving addr thru env

type KafkaConfig struct {
	Addr []string
}

func NewKafkaConfig() *KafkaConfig {
	return &KafkaConfig{
		Addr: []string{"localhost:9092"},
	}
}

type MsgOperation struct {
	Operation string `json:"operation"`
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
