package types

type Update struct {
	Text   string
	From   *User
	ChatID GroupID
}

type User struct {
	ID       UserID
	Username string
}
