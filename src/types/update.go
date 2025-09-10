package types

type Update struct {
	ID     int64
	Text   string
	From   *User
	ChatID GroupID
}

type User struct {
	ID       UserID
	Username string
}
