package users

import (
	"bookstore_users-api/datasources/mysql/users_db"
	"bookstore_users-api/logger"
	"bookstore_users-api/utils/date_utils"
	"bookstore_users-api/utils/errors"
	"fmt"
	"strings"
)

const (
	indexUniqueEmail      = "email"
	errorNoRows           = "errorNoRows"
	queryInsertUser       = "INSERT INTO users (first_name, last_name, email, date_created, status, password) VALUES(?,?,?,?,?,?)"
	queryGetUser          = "SELECT id,first_name, last_name, email, date_created, status, password FROM users WHERE id=?  "
	queryUpdateUser       = "UPDATE users SET first_name=?, last_name=?, email=?, date_created=? where id = ?"
	queryDeleteUser       = "Delete from users where id = ?"
	queryFindUserByStatus = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE status = ?;"
	queryFindUsers        = "SELECT id, first_name, last_name, email, date_created, status, password FROM users;"
)

var (
	usersDB = make(map[int64]*User)
)

func (user *User) Get(userID int64) *errors.RestErr {

	stmt, err := users_db.Client.Prepare(queryGetUser)
	if err != nil {

		logger.Error("error when trying to prepare get user ", err)
		return errors.NewInternalServerError(err.Error())
	}

	defer stmt.Close()

	result := stmt.QueryRow(userID)
	if err := result.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status, &user.Password); err != nil {

		if strings.Contains(err.Error(), errorNoRows) {

			return errors.NewNotFoundError(fmt.Sprintf("user %d not found", user.ID))

		}

		return errors.NewInternalServerError(fmt.Sprintf("error when trying to get user %d: %s", user.ID, err.Error()))
	}

	// if err := users_db.Client.Ping(); err != nil {
	// 	panic(err)
	// }
	//
	// result := usersDB[userID]
	//
	// if result == nil {
	// 	return errors.NewNotFoundError(fmt.Sprintf("user %d not found", userID))
	// }
	// user.ID = result.ID
	// user.FirstName = result.FirstName
	// user.LastName = result.LastName
	// user.Email = result.Email
	// user.DateCreated = result.DateCreated

	// result, err := users_db.Client.Query("select * from users where id = ?", userID)
	// if err != nil {
	// 	// if strings.Contains(err.Error(), indexUniqueEmail) {
	// 	// 	return errors.NewNotInternalServerError(fmt.Sprintf("email %s already exists", user.Email))
	// 	// }
	//
	// 	return errors.NewNotFoundError("user not foun")
	// }
	//
	// fmt.Println()
	//
	// for result.Next() {
	// 	fmt.Println(result)
	// }

	return nil
}

func (user *User) Save() *errors.RestErr {

	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}

	defer stmt.Close()

	user.DateCreated = date_utils.GetNowDBFormat()

	insertResult, err := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated, user.Status, user.Password)
	if err != nil {

		// sqlErr, err := err

		if strings.Contains(err.Error(), indexUniqueEmail) {
			return errors.NewInternalServerError(fmt.Sprintf("email %s already exists", user.Email))
		}

		return errors.NewInternalServerError(fmt.Sprintf("error when trying to save user: %s", err.Error()))
	}

	userID, err := insertResult.LastInsertId()
	if err != nil {
		return errors.NewInternalServerError(fmt.Sprintf("error when trying to save user: %s", err.Error()))
	}

	user.ID = userID

	return nil
}

func (user *User) Update(userID int64) *errors.RestErr {

	stmt, err := users_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}

	defer stmt.Close()

	user.DateCreated = date_utils.GetNowString()

	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated, userID)
	if err != nil {

		// sqlErr, err := err

		if strings.Contains(err.Error(), indexUniqueEmail) {
			return errors.NewInternalServerError(fmt.Sprintf("email %s already exists", user.Email))
		}

		return errors.NewInternalServerError(fmt.Sprintf("error when trying to save user: %s", err.Error()))
	}

	// userID, err := insertResult.LastInsertId()
	// if err != nil {
	// 	return errors.NewNotInternalServerError(fmt.Sprintf("error when trying to save user: %s", err.Error()))
	// }
	//
	// user.ID = userID

	return nil
}

func (user *User) Delete(userID int64) (int64, *errors.RestErr) {

	stmt, err := users_db.Client.Prepare(queryDeleteUser)
	if err != nil {
		return 0, errors.NewInternalServerError(err.Error())
	}

	defer stmt.Close()

	result, err := stmt.Exec(userID)
	if err != nil {

		return 0, errors.NewInternalServerError(fmt.Sprintf("error when trying to save user: %s", err.Error()))
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {

		return 0, errors.NewInternalServerError(fmt.Sprintf("error when trying to save user: %s", err.Error()))
	}

	return rowsAffected, nil
}

// find data from status
func (user *User) FindByStatus(status string) ([]User, *errors.RestErr) {

	stmt, err := users_db.Client.Prepare(queryFindUserByStatus)
	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	rows, err := stmt.Query(status)
	defer rows.Close()
	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}
	defer rows.Close()

	results := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.DateCreated, &user.Email, &user.Status); err != nil {
			return nil, errors.NewInternalServerError(err.Error())
		}

		results = append(results, user)
	}

	if len(results) == 0 {
		return nil, errors.NewNotFoundError(fmt.Sprintf("no users matching status %s", status))
	}

	return results, nil

}

func (user *User) FindAllUsers() ([]User, *errors.RestErr) {

	stmt, err := users_db.Client.Prepare(queryFindUsers)
	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	defer rows.Close()
	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}
	defer rows.Close()

	results := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.DateCreated, &user.Email, &user.Status, &user.Password); err != nil {
			return nil, errors.NewInternalServerError(err.Error())
		}

		results = append(results, user)
	}

	if len(results) == 0 {
		return nil, errors.NewNotFoundError("not user found")
	}

	return results, nil

}
