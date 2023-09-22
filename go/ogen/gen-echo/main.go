package main

import (
	"github.com/GIT_USER_ID/GIT_REPO_ID/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

    //todo: handle the error!
	c, _ := handlers.NewContainer()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())


	// AddPet - Add a new pet to the store
	e.POST("/v3/pet", c.AddPet)

	// DeletePet - Deletes a pet
	e.DELETE("/v3/pet/:petId", c.DeletePet)

	// GetPetById - Find pet by ID
	e.GET("/v3/pet/:petId", c.GetPetById)

	// UpdatePet - Updates a pet in the store
	e.POST("/v3/pet/:petId", c.UpdatePet)


	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}