package controllers

import (
	"github.com/Akmdvtm/database"
	"github.com/Akmdvtm/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"
)

const (
	secretKey = "secret"
)

// Register user | add user info in db (password hash)
func Register(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		panic(err)
	}

	//Crypt password
	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 1)
	user := &models.User{
		Name:     data["name"],
		Email:    data["email"],
		Password: string(password),
	}

	//Создаем user
	database.DB.Create(&user)

	return c.JSON(user)
}

//Login user | return jwt token with cookie
func Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		panic(err)
	}

	var user models.User

	//Получаем пользователя на основе полученных данных (email)
	database.DB.Where("email = ?", data["email"]).First(&user)

	//Если не валидные данные
	if user.Id == 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "user not found",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data["password"])); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "incorrect password",
		})
	}

	//Create JWT Token
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.Id)),
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days срок действия
	})

	token, err := claims.SignedString([]byte(secretKey))
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "could not login",
		})
	}

	//Create cookie
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24 * 7),
		HTTPOnly: true, //Frontend может запрашивать, но не получать доступ
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "success",
	})
}

// User Get authenticated user
func User(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}
	claims := token.Claims.(*jwt.StandardClaims)

	var user models.User

	database.DB.Where("id = ?", claims.Issuer).First(&user)

	return c.JSON(user)
}

//Logout user
func Logout(c *fiber.Ctx) error {
	// Объявляем что куки истекли час назад
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "success",
	})
}
