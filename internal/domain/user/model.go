package user

import (
	"encoding/json"
	"log"
	"time"

	"github.com/evermos/boilerplate-go/shared/jwtmodel"
	"github.com/evermos/boilerplate-go/shared/nuuid"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt"
	"github.com/guregu/null"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id uuid.UUID `db:"id"`
	Username 	string `db:"username"`
	Name 	 	string `db:"name"`
	Password 	string `db:"password"`
	Role 	 	string `db:"role"`
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


	return
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

	return
}


func GenerateJWT(user User) (string, error)  {
	secret := viper.GetString("JWT_SECRET")

	claims := jwtmodel.Claims{
		UserId: user.Id,
		Username: user.Username,
		Role: user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
			Issuer: "evermos",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err 
	}

	return tokenString, nil
}

func (u User) ToResponseFormat() UserResponseFormat {
	accessToken, err := GenerateJWT(u)
	if err != nil {
		log.Println(err.Error())
	} 
	resp := UserResponseFormat{
		Id: u.Id,
		Username: u.Username,
		Name: u.Name,
		Role: u.Role,
		AccessToken: accessToken,
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
	Id uuid.UUID `db:"id"`
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
	accessToken, err := GenerateJWT(l.User)

	if err != nil {
		log.Println(err.Error())
	} 
	resp := LoginResponseFormat{
		AccessToken: accessToken,

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


