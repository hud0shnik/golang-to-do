package todo

type User struct {
	Id       int    `json:"-"`
	Name     string `json:"name"`
	UserName string `json:"username"`
	Password string `json:"password"`
}
