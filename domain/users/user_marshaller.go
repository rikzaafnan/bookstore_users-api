package users

import "github.com/goccy/go-json"

type PublicUser struct {
	ID int64 `json:"user_id"`
	// FirstName   string `json:"first_name"`
	// LastName    string `json:"last_name"`
	// Email       string `json:"email"`
	DateCreated string `json:"date_created"`
	Status      string `json:"status"`
	// Password    string `json:"password"`
}

type PrivateUser struct {
	ID          int64  `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	DateCreated string `json:"date_created"`
	Status      string `json:"status"`
}

func (user *User) Marshall(isPublic bool) interface{} {
	if isPublic {
		return PublicUser{
			ID:          user.ID,
			DateCreated: user.DateCreated,
			Status:      user.Status,
		}
	}

	userJson, _ := json.Marshal(user)

	var privateUser PrivateUser
	json.Unmarshal(userJson, &privateUser)

	return privateUser
}

func (users Users) Marshalls(isPublic bool) []interface{} {
	results := make([]interface{}, len(users))

	for index, user := range users {
		results[index] = user.Marshall(isPublic)
	}

	return results

}
