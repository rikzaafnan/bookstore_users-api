package users

import (
	"bookstore-users-api/datasources/mysql/users_db"
	"bookstore-users-api/utils/crypto_utils"
	"bookstore-users-api/utils/date_utils"
	"bookstore-users-api/utils/errors"
	"fmt"
	"log"
	"strings"

	"bookstore-users-api/logger"

	"github.com/go-sql-driver/mysql"
)

var (
	usersDB = make(map[int64]*User)
)

type Users []User

const (
	indexUniqueEmail      = "email"
	errorNoRows           = "no rows in result set"
	queryInsertUser       = "INSERT INTO users(first_name, last_name, email, date_created, password, status) VALUES (?, ?, ?, ?,?,?)"
	queryGetUser          = "SELECT id, first_name, last_name, email, date_created,status FROM users"
	queryUpdateUser       = "UPDATE users set first_name=?, last_name=?, email=? where id=?"
	queryDeleteUserByID   = "DELETE FROM users where id=?"
	queryFindUserByStatus = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE status=?"
)

func (user *User) Save() *errors.RestErr {

	stmt, err := users_db.ClientDB.Prepare(queryInsertUser)
	if err != nil {

		return errors.NewInternalServerError(err.Error())
	}

	defer stmt.Close()

	user.DateCreated = date_utils.GetNowDBFormatString()
	user.Status = StatusActive

	user.Password = crypto_utils.GetMd5(user.Password)

	insertResult, saveErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated, user.Password, user.Status)
	if saveErr != nil {
		log.Println(saveErr.(*mysql.MySQLError))

		sqlErr, ok := saveErr.(*mysql.MySQLError)
		if !ok {
			return errors.NewInternalServerError(fmt.Sprintf("error when trying to save user %s", err.Error()))
		}

		switch sqlErr.Number {
		case 1062:
			return errors.NewBadRequestError(fmt.Sprintf("email %s already exist", user.Email))
		}

		return errors.NewInternalServerError(fmt.Sprintf("error when trying to save user %s", saveErr.Error()))
	}

	userId, err := insertResult.LastInsertId()
	if err != nil {
		return errors.NewInternalServerError(fmt.Sprintf("error when trying to save user %s", err.Error()))
	}

	user.ID = userId

	return nil

}

func (user *User) Get() *errors.RestErr {

	if err := users_db.ClientDB.Ping(); err != nil {
		panic(err)
	}

	result := usersDB[user.ID]
	if result == nil {
		return errors.NewNotFoundError(fmt.Sprintf("user %d not found", user.ID))
	}

	user.ID = result.ID
	user.FirstName = result.FirstName
	user.LastName = result.LastName
	user.Email = result.Email
	user.DateCreated = result.DateCreated

	return nil

}

func (user *User) FindByID() *errors.RestErr {
	err := users_db.ClientDB.QueryRow(queryGetUser+" where id = ?", user.ID).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status)
	if err != nil {

		if strings.Contains(err.Error(), errorNoRows) {
			return errors.NewNotFoundError(fmt.Sprintf("user id %d not found", user.ID))
		}
		logger.Error("ini ada error : ", err)
		return errors.NewBadRequestError(fmt.Sprintf("error when trying to get user id %d: %s", user.ID, err.Error()))
	}

	return nil
}

func (user *User) GetUserByEmail() *errors.RestErr {

	err := users_db.ClientDB.QueryRow(queryGetUser+" where email = ?", user.Email).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status)
	if err != nil {
		return errors.NewInternalServerError(fmt.Sprintf("error when trying to get user email %s", err.Error()))
	}

	// rows, err := results
	// if err != nil {
	// 	log.Fatal(err)
	// 	return errors.NewInternalServerError(err.Error())
	// }
	// if rows != 1 {
	// 	log.Fatalf("expected to affect 1 row, affected %d", rows)
	// 	return errors.NewInternalServerError("failed get data")
	// }

	return nil

}

func (user *User) UpdateUser() *errors.RestErr {

	stmt, err := users_db.ClientDB.Prepare(queryUpdateUser)
	if err != nil {

		return errors.NewInternalServerError(err.Error())
	}

	defer stmt.Close()

	_, saveErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.ID)
	if saveErr != nil {
		log.Println(saveErr.(*mysql.MySQLError))

		sqlErr, ok := saveErr.(*mysql.MySQLError)
		if !ok {
			return errors.NewInternalServerError(fmt.Sprintf("error when trying to save user %s", err.Error()))
		}

		switch sqlErr.Number {
		case 1062:
			return errors.NewBadRequestError(fmt.Sprintf("email %s already exist", user.Email))
		}

		return errors.NewInternalServerError(fmt.Sprintf("error when trying to save user %s", saveErr.Error()))
	}

	return nil

}

func (user *User) DeleteByID() *errors.RestErr {

	stmt, err := users_db.ClientDB.Prepare(queryDeleteUserByID)
	if err != nil {

		return errors.NewInternalServerError(err.Error())
	}

	defer stmt.Close()

	_, deleteErr := stmt.Exec(user.ID)
	if deleteErr != nil {
		log.Println(deleteErr.(*mysql.MySQLError))

		sqlErr, ok := deleteErr.(*mysql.MySQLError)
		if !ok {
			return errors.NewInternalServerError(fmt.Sprintf("error when trying to save user %s", err.Error()))
		}

		switch sqlErr.Number {
		case 1062:
			return errors.NewBadRequestError(fmt.Sprintf("email %s already exist", user.Email))
		}

		return errors.NewInternalServerError(fmt.Sprintf("error when trying to save user %s", deleteErr.Error()))
	}

	return nil

}

func (user *User) FindByStatus(status string) ([]User, *errors.RestErr) {
	stmt, err := users_db.ClientDB.Prepare(queryFindUserByStatus)
	if err != nil {

		return nil, errors.NewInternalServerError(err.Error())
	}

	defer stmt.Close()

	rows, err := stmt.Query(status)
	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}
	defer rows.Close()

	results := make([]User, 0)

	for rows.Next() {
		var user User

		err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status)
		if err != nil {
			return nil, errors.NewInternalServerError(err.Error())
		}

		results = append(results, user)
	}

	if len(results) == 0 {
		return nil, errors.NewNotFoundError(fmt.Sprintf("no users matching status %s", status))
	}

	return results, nil

}
