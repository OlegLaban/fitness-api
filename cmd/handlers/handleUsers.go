package handlers

import (
	"fitness-api/cmd/models"
	"fitness-api/cmd/repositories"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(c echo.Context) error {
	var err error
	user := models.User{}
	c.Bind(&user)
	user.Password, err = hashPassword(user.Password);

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	newUser, err := repositories.CreateUser(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, newUser)
}

func HandleUpdateUser(c echo.Context) error {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	user := models.User{}
	c.Bind(&user)
	updatedUser, err := repositories.UpdateUser(user, idInt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, updatedUser)
}

func HandleGetUsers(c echo.Context) error {
	users, err := repositories.UserList()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, users)
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	
	return string(bytes), nil
}

func CheckPassword(user models.User, providerPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providerPassword))

	if err != nil {
		return err
	}

	return nil
}