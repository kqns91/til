package handlers
import (
    "github.com/GIT_USER_ID/GIT_REPO_ID/models"
    "github.com/labstack/echo/v4"
    "net/http"
)

// AddPet - Add a new pet to the store
func (c *Container) AddPet(ctx echo.Context) error {
    return ctx.JSON(http.StatusOK, models.HelloWorld {
        Message: "Hello World",
    })
}


// DeletePet - Deletes a pet
func (c *Container) DeletePet(ctx echo.Context) error {
    return ctx.JSON(http.StatusOK, models.HelloWorld {
        Message: "Hello World",
    })
}


// GetPetById - Find pet by ID
func (c *Container) GetPetById(ctx echo.Context) error {
    return ctx.JSON(http.StatusOK, models.HelloWorld {
        Message: "Hello World",
    })
}


// UpdatePet - Updates a pet in the store
func (c *Container) UpdatePet(ctx echo.Context) error {
    return ctx.JSON(http.StatusOK, models.HelloWorld {
        Message: "Hello World",
    })
}
