package user

import (
	"encoding/json"
	"log"
	"time"

	"github.com/evermos/boilerplate-go/shared"
	"github.com/evermos/boilerplate-go/shared/nuuid"

	"github.com/gofrs/uuid"
	"github.com/guregu/null"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id uuid.UUID `db:"id"`
	Username 	string `db:"username"`
	Name 	 	string `db:"name"`
	Password 	string `db:"password"`
	Role 	 	string `db:"role"`
	AccessToken string `db:"-"`
	CreatedAt   time.Time   `db:"createdAt"`
	CreatedBy   uuid.UUID   `db:"createdBy"`
	UpdatedAt   null.Time   `db:"updatedAt"`
	UpdatedBy   nuuid.NUUID `db:"updatedBy"`
	DeletedAt   null.Time   `db:"deletedAt"`
	DeletedBy   nuuid.NUUID `db:"deletedBy"`
}

func (u *User) IsDeleted() (deleted bool) {
	return u.DeletedAt.Valid && u.DeletedBy.Valid
}

func (u User) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.ToResponseFormat())
}

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

func (u *User) Update(req UserRequestFormat, user User) (err error) {
	hashPassword, err := HashPassword(req.Password)
	if err != nil {
		log.Println(err.Error())
		return 
	}

	u.Username = req.Username
	u.Name = req.Name
	u.Password = hashPassword
	u.Role = req.Role
	u.UpdatedAt = null.TimeFrom(time.Now())
	u.UpdatedBy = nuuid.From(user.Id)

	err = u.Validate()
	return
}

func (u *User) Validate() (err error) {
	validator := shared.GetValidator()
	return validator.Struct(u)
}

func (u User) NewFromRequestFormat(req UserRequestFormat) (newUser User, err error) {
	userID, _ := uuid.NewV4()
	passwordHashed, err := HashPassword(req.Password)
	if err != nil {
		log.Println(err.Error())
		return
	} 
	newUser = User{
		Id: userID,
		Username: req.Username,
		Name: req.Name,
		Password: passwordHashed,
		Role: req.Role,
		CreatedAt:   time.Now(),
		CreatedBy:   userID,
	}

	err = newUser.Validate()

	return
}




func (u User) ToResponseFormat() UserResponseFormat {
	resp := UserResponseFormat{
		Id: u.Id,
		Username: u.Username,
		Name: u.Name,
		Role: u.Role,
		AccessToken: u.AccessToken,
		CreatedAt: u.CreatedAt,
		CreatedBy: u.CreatedBy,
		UpdatedAt: u.UpdatedAt,
		UpdatedBy: u.UpdatedBy.Ptr(),
		DeletedAt: u.DeletedAt,
		DeletedBy: u.DeletedBy.Ptr(),

	}
	return resp
}



type UserRequestFormat struct {
	Username 	string  `json:"username"`
	Name 	    string  `json:"name"`
	Password 	string  `json:"password"`
	Role 	    string  `json:"role"`
}

type UserResponseFormat struct {
	Id uuid.UUID 			`json:"id"`
	Username 	string 		`json:"username"`
	Name 	 	string 		`json:"name"`
	Role 	 	string 		`json:"role"`
	AccessToken string 		`json:"accessToken"`
	CreatedAt   time.Time   `json:"createdAt"`
	CreatedBy   uuid.UUID   `json:"createdBy"`
	UpdatedAt   null.Time   `json:"updatedAt,omitempty"`
	UpdatedBy   *uuid.UUID 	`json:"updatedBy,omitempty"`
	DeletedAt   null.Time   `json:"deletedAt,omitempty"`
	DeletedBy   *uuid.UUID 	`json:"deletedBy,omitempty"`	
}


type Login struct {
	Username 	string 		
	Password 	string  	
	User		User
	AccessToken string
	
}

func (l Login) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.ToResponseFormat())
}

func (l Login) NewFromRequestFormat(req LoginRequestFormat) (newLogin Login, err error) {
	if err != nil {
		log.Println(err.Error())
		return
	} 
	newLogin = Login{
		Username: req.Username,
		Password: req.Password,
	}

	return
}

func (l Login) ToResponseFormat() LoginResponseFormat {
	
	resp := LoginResponseFormat{
		AccessToken: l.AccessToken,

	}
	return resp
}

type LoginRequestFormat struct {
	Username 	string 		`json:"username" validate:"required"`
	Password 	string  	`json:"password" validate:"required"`
	
}

type LoginResponseFormat struct {
	AccessToken string 		`json:"accessToken"`
}


