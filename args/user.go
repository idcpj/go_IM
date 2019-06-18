package args

type User struct {
	UserId int64 `json:"user_id" form:"user_id"`
	Mobile string `json:"mobile" form:"mobile"`
}
