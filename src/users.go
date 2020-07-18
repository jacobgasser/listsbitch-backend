package main

func GetIDFromUsername(username string) string {
	user := User{}
	DB.First(&user, "username = ?", username)
	return user.ID
}
