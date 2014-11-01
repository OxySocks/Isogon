package main

import (
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"errors"

	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
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

// Function that handles the login for Users. Compares the password available in the struct to the
// Hashed password in the database. Sets session to include the user if login was successful.
// If login went wrong, users will be redirected if they have provided invalid credentials,
// If any other error occurs a 500 error will be thrown.
func LoginUser(req *http.Request, s sessions.Session, res render.Render, db *gorm.DB, user User) {
	user, err := user.Login(db)

	if err != nil {
		log.Println(err)

		if err.Error() == "Wrong username or password." || err.Error() == "not found" {
			res.HTML(401, "users/login", "Wrong username or password.")
			return
		}
		log.Println("[WARNING] Unhandled error caught at user login")
		res.HTML(500, "users/login", "Internal server error. Please try again.")
		return
	}

	s.Set("user", user.ID)
	res.Redirect("/", 302)
	return
}

// GORM wrapper function that grabs a user from the database given the information available in the struct.
// Requires at least an ID to be set for this user so it can be grabbed from the database.
func (user User) Get(db *gorm.DB) (User, error) {
	query := db.Where(&User{ID: user.ID}).First(&user)

	if query.Error != nil {
		if query.Error == gorm.RecordNotFound {
			return user, errors.New("not found")
		}
		return user, query.Error
	}

	return user, nil
}

// Function to grab a user from the current sesion. Returns an error if no User is stored in the session
// Always returns a User struct, but it can/will be empty when no data was found.
func (user User) FromSession(db *gorm.DB, s sessions.Session) (User, error) {
	data := s.Get("user")
	id, exists := data.(int64)
	if exists {
		var user User
		user.ID = id
		user, err := user.Get(db)
		if err != nil {
			return user, err
		}
		return user, nil
	}
	return user, errors.New("unauthorized")
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
