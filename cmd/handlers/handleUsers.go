package handlers

import (
	"errors"
	"fitness-api/cmd/models"
	"fitness-api/cmd/repositories"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type TokenRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

func CreateUser(c echo.Context) error {
	var err error
	user := models.User{}
	c.Bind(&user)

	err = hashPassword(&user.Password);

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

func hashPassword(password *string) (error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(*password), 14)
	*password = string(bytes)

	if err != nil {
		return err
	}
	
	return nil
}

func Auth(c echo.Context) error {
	var request TokenRequest

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	user, err := repositories.GetByEmail(request.Email)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	credentionalError := CheckPassword(user, request.Password)

	if credentionalError != nil {
		authError := errors.New("invalid creadentials")
		return c.JSON(http.StatusUnauthorized, authError)
	}

	tokenString, err := GenerateJWT(user.Email, user.Name)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, tokenString)
}

func CheckPassword(user models.User, providerPassword string) error {
	fmt.Println(user.Password)
	fmt.Println(providerPassword)
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providerPassword))

	if err != nil {
		return err
	}

	return nil
}