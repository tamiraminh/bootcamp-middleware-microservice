package user

import (
	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/shared/failure"
)

type UserService interface {
	Create(requestFormat UserRequestFormat) (user User, err error)
	Login(requestFormat LoginRequestFormat) (login Login, err error)
	ResolveByUsername(username string) (user User, err error)
}

type UserServiceImpl struct {
	UserRepository UserRepository
	Config *configs.Config
}

func ProvideUserServiceImpl(userRepository UserRepository, config *configs.Config)  *UserServiceImpl {
	s := new(UserServiceImpl)
	s.UserRepository = userRepository
	s.Config = config

	return s
}

func (s *UserServiceImpl) Create(requestFormat UserRequestFormat) (user User, err error)  {
	user, err = user.NewFromRequestFormat(requestFormat)
	if err != nil {
		return user, failure.BadRequest(err)
	}

	err = s.UserRepository.Create(user)
	if err != nil {
		return
	}

	return
}


func (s *UserServiceImpl) Login(requestFormat LoginRequestFormat) (login Login, err error)  {
	login, err = login.NewFromRequestFormat(requestFormat)
	if err != nil {
		return login, failure.BadRequest(err)
	}
	var user User
	user, err = s.UserRepository.ResolveByUsername(login.Username)
	if err != nil {
		return
	}

	
	
	login.User = user
	
	if match := CheckPasswordHash(login.Password, login.User.Password); !match {
		return Login{}, failure.BadRequestFromString("Password False!")
	}

	return
}



func (s *UserServiceImpl) ResolveByUsername(username string) (user User, err error)  {
	user, err = s.UserRepository.ResolveByUsername(username)

	if user.IsDeleted() {
		return user, failure.NotFound("User")
	}


	return
}

