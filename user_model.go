package main

type User struct {
	Id       string
	Password string
}

func findUser(id string) *User {
	var user User

	users := []User{{"user1", "passwd1"}, {"user2", "passwd2"}}
	for k, v := range users {
		if v.Id == id {
			user = users[k]
		}
	}

	return &user
}
