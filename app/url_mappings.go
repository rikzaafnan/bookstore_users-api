package app

import controller "bookstore-users-api/controllers"

func mapUrls() {

	router.GET("/ping", controller.Ping)

	router.GET("/users/:userID", controller.GetUser)
	router.GET("/users/search", controller.SearchUser)
	router.POST("/users", controller.CreateUser)
}
