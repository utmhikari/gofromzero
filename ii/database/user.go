package database

type User struct {
	*DBModel
	Name string `json:"name"`
	Age uint `json:"age"`
}

type userDAO struct {}

var UserDAO userDAO

// Create create a user record
func (*userDAO) Create(user User) error {
	DBInstance.AutoMigrate(&User{})
	return DBInstance.Create(&user).Error
}

// First get the first record of user
func (*userDAO) First() (User, error) {
	var user User
	err := DBInstance.First(&user).Error
	return user, err
}

// Update update user record
func (*userDAO) Update(user User) error {
	return DBInstance.Model(&User{}).Updates(&user).Error
}

// Delete delete all
func (*userDAO) Delete() error {
	return DBInstance.Delete(&User{}).Error
}