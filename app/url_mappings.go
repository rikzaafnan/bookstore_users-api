package app

import controller "bookstore-users-api/controllers"

func mapUrls() {

	router.GET("/ping", controller.Ping)

	router.GET("/users/:userID", controller.GetUser)
	router.GET("/users/search", controller.SearchUser)
	router.POST("/users", controller.CreateUser)
	router.PUT("/users/:userID", controller.UpdateUser)
	router.PATCH("/users/:userID", controller.UpdateUser)
	router.DELETE("/users/:userID", controller.DeleteUser)

	router.GET("/internal/users/search", controller.Search)
}
