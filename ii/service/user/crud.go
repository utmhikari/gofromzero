package user

import "github.com/gofromzero/ii/database"

type Form struct {
	Name string `json:"name"`
	Age  uint   `json:"age"`
}

// Create create user on form
func Create(form Form) error {
	return database.UserDAO.Create(database.User{
		Name: form.Name,
		Age:  form.Age,
	})
}

// First get first user from database
func First() (database.User, error) {
	return database.UserDAO.First()
}

// Update update user records
func Update(form Form) error {
	return database.UserDAO.Update(database.User{
		Name: form.Name,
		Age:  form.Age,
	})
}

// Delete delete user records
func Delete() error {
	return database.UserDAO.Delete()
}
