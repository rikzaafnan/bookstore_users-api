package controller

import (
	"bookstore-users-api/domain/users"
	"bookstore-users-api/services"
	"bookstore-users-api/utils/errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {

	var user users.User

	err := c.ShouldBind(&user)
	if err != nil {
		restErr := errors.NewBadRequestError("invalid json body")

		c.JSON(restErr.Status, restErr)

		return
	}

	result, saveErr := services.UsersService.CreatedUser(user)
	if saveErr != nil {

		c.JSON(saveErr.Status, saveErr)
		return
	}

	c.JSON(http.StatusCreated, result.Marshall(c.GetHeader("X-Public") == "true"))
}

func GetUser(c *gin.Context) {

	userID, userErr := strconv.ParseInt(c.Param("userID"), 10, 64)
	if userErr != nil {

		err := errors.NewBadRequestError("user id should be a number")
		c.JSON(err.Status, err)
		return

	}

	user, saveErr := services.UsersService.GetUser(userID)
	if saveErr != nil {

		c.JSON(saveErr.Status, saveErr)
		return
	}

	c.JSON(http.StatusOK, user.Marshall(c.GetHeader("X-Public") == "true"))
}

func FindUser(c *gin.Context) {
	c.String(http.StatusNotImplemented, "implemented me")
}

func SearchUser(c *gin.Context) {
	c.String(http.StatusNotImplemented, "implemented me")
}

func UpdateUser(c *gin.Context) {

	var user users.User

	userID, userErr := strconv.ParseInt(c.Param("userID"), 10, 64)
	if userErr != nil {

		err := errors.NewBadRequestError("user id should be a number")
		c.JSON(err.Status, err)
		return

	}

	err := c.ShouldBind(&user)
	if err != nil {
		restErr := errors.NewBadRequestError("invalid json body")

		c.JSON(restErr.Status, restErr)

		return
	}

	user.ID = userID

	isPartial := c.Request.Method == http.MethodPatch

	result, updateErr := services.UsersService.UpdateUser(isPartial, user)
	if updateErr != nil {

		c.JSON(updateErr.Status, updateErr)
		return
	}

	c.JSON(http.StatusCreated, result.Marshall(c.GetHeader("X-Public") == "true"))
}

func DeleteUser(c *gin.Context) {

	var user users.User

	userID, userErr := strconv.ParseInt(c.Param("userID"), 10, 64)
	if userErr != nil {

		err := errors.NewBadRequestError("user id should be a number")
		c.JSON(err.Status, err)
		return

	}

	user.ID = userID

	deleteErr := services.UsersService.DeleteByUserID(userID)
	if deleteErr != nil {

		c.JSON(deleteErr.Status, deleteErr)
		return
	}

	c.JSON(http.StatusCreated, "OK")
}

func Search(c *gin.Context) {
	status := c.Query("status")

	users, err := services.UsersService.FindByStatus(status)
	if err != nil {

		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusCreated, users.Marshalls(c.GetHeader("X-Public") == "true"))
}
