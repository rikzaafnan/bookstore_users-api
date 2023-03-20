package services

import (
	"bookstore-users-api/domain/users"
	"bookstore-users-api/utils/errors"
)

var (
	UsersService usersServiceInterface = &usersService{}
)

type usersService struct {
}

type usersServiceInterface interface {
	CreatedUser(user users.User) (*users.User, *errors.RestErr)
	GetUser(userID int64) (*users.User, *errors.RestErr)
	UpdateUser(isPartial bool, user users.User) (*users.User, *errors.RestErr)
	DeleteByUserID(userID int64) *errors.RestErr
	FindByStatus(status string) (users.Users, *errors.RestErr)
}

func (userServices *usersService) CreatedUser(user users.User) (*users.User, *errors.RestErr) {

	err := user.Validate()
	if err != nil {
		return nil, err
	}

	// err = user.GetUserByEmail()
	// if err == nil {

	// 	return nil, errors.NewBadRequestError(fmt.Sprintf("email %s already exist", user.Email))
	// }

	err = user.Save()
	if err != nil {
		return nil, err
	}

	return &user, nil

}

func (userServices *usersService) GetUser(userID int64) (*users.User, *errors.RestErr) {

	if userID <= 0 {
		return nil, errors.NewBadRequestError("invalid user id")
	}

	result := &users.User{ID: userID}

	err := result.FindByID()
	if err != nil {
		return nil, err
	}

	return result, nil

}

func (userServices *usersService) UpdateUser(isPartial bool, user users.User) (*users.User, *errors.RestErr) {

	// currentUser := users.User{ID: user.ID}

	// err := currentUser.FindByID()
	// if err != nil {
	// 	return nil, err
	// }

	currentUser, err := userServices.GetUser(user.ID)
	if err != nil {
		return nil, err
	}

	if isPartial {
		currentUser.Email = user.Email

		if user.LastName != "" {
			currentUser.FirstName = user.FirstName
		}
		if user.LastName != "" {
			currentUser.FirstName = user.FirstName
		}
		if user.LastName != "" {
			currentUser.FirstName = user.FirstName
		}
	} else {
		currentUser.FirstName = user.FirstName
		currentUser.LastName = user.LastName
		currentUser.Email = user.Email
	}

	err = currentUser.UpdateUser()
	if err != nil {
		return nil, err
	}

	return currentUser, nil

}

func (userServices *usersService) DeleteByUserID(userID int64) *errors.RestErr {

	currentUser, err := userServices.GetUser(userID)
	if err != nil {
		return err
	}

	err = currentUser.DeleteByID()
	if err != nil {
		return err
	}

	return nil

}

func (userServices *usersService) FindByStatus(status string) (users.Users, *errors.RestErr) {

	dao := &users.User{}
	users, err := dao.FindByStatus(status)
	if err != nil {
		return nil, err
	}

	return users, nil

}
