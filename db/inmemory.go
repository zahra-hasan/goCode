package db

import (
	"slices"
)

// Create a struct for users
// notice the items in back ticks refer to the fields in the JSON data
// the struct uses the capitalised keys but when converted to json it will use the lower case ones
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// mocking ID
var count int = 3

// declare users
var users []User

func init() {
	// initialise in memory db with data (Create hardcoded slice using struct)
	users = []User{
		{ID: 1, Name: "Name 1"},
		{ID: 2, Name: "Name 2"},
		{ID: 3, Name: "Name 3"},
	}
}

// function that we use in our main.go "db.GetUsers()"
func GetUsers() []User {
	return users
}

func GetUser(id int) User {
	var user User

	for _, user := range users {
		if user.ID == id {
			return user
		}
	}

	return user
}

func AddUser(user User) int {
	count++
	user.ID = count

	users = append(users, user)

	return count
}

func DeleteUser(id int) []User {

	for index, user := range users {
		if user.ID == id {
			// slices.Delete will delete an item or items based on the passed in index
			start, end := index, index+1
			users = slices.Delete(users, start, end)
		}
	}
	return users
}

func UpdateUser(id int, userData User) User {
	var user User
	for index, user := range users {
		if user.ID == id {
			start, end := index, index+1
			userData.ID = id
			users = slices.Replace(users, start, end, userData)
			return users[index]
		}
	}

	return user
}
