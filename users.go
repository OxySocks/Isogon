package main

import (
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"errors"

	"github.com/codegangsta/martini-contrib/render"
)

// Function that handles the registration of users. Currently only handles GET requests.
// Returns 403 if registrations are not allowed at this time, as placed in the global Settings variable.
func RegisterUser(req *http.Request, r render.Render, db *gorm.DB, user User) {
	user, err := user.Save(db)

	if err != nil {
		log.Println(err)
		if err.Error() == "UNIQUE constraint failed: users.email" {
			r.Error(422)
			return
		}
		r.Error(500)
		return
	}

	user, err = user.Login(db)

	if err != nil {
		log.Println(err)
		r.Error(500)
		return
	}

	r.Redirect("/", 302)
}

// Method to get a user by  e-mail address. Used in forms as a User is returned from Martini bindings.
// Returns an error if the User could not be found. Always returns the user.
func (user User) GetByEmail(db *gorm.DB) (User, error) {
	query := db.Where(&User{Email: user.Email}).First(&user)

	if query.Error != nil {
		if query.Error == gorm.RecordNotFound {
			return user, errors.New("User not found.")
		}

		return user, query.Error
	}

	return user, nil
}

// Method to log a user in using the database. Uses the GetByEmail function to see if the User's e-mail address
// is present in the database. Returns an error if the login was not successful. Always returns the user.
func (user User) Login(db *gorm.DB) (User, error) {
	user, err := user.GetByEmail(db)

	if err != nil {
		return user, err
	}

	if !CompareHash(user.PasswordHash, user.Password) {
		return user, errors.New("Wrong username or password.")
	}

	return user, nil
}

// Method to save a user to the database. Generates a password hash (bcrypt) using the user's password field.
// The Password field is NOT stored in the database, only hashes are actually stored. The Password field
// get returned from the Martini bindings and is only used for temporary storage.
func (user User) Save(db *gorm.DB) (User, error) {
	passwordHash, err := GenerateHash(user.Password)

	if err != nil {
		return user, err
	}

	user.PasswordHash = passwordHash
	query := db.Create(&user)

	if query.Error != nil {
		return user, query.Error
	}

	return user, nil
}
