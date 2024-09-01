package controllers

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/samuel-fonseca/task-manager-api/database"
	"github.com/samuel-fonseca/task-manager-api/model"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(c *fiber.Ctx) error {
	var payload = new(model.UserRegistrationData)
	err := c.BodyParser(&payload)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if payload.Password != payload.PasswordConfirm {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Your password confirmation does match the password.",
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	user := model.User{
		Username: payload.Username,
		Email:    payload.Email,
		Password: string(hashedPassword),
	}

	errors := ValidateStruct(user)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": errors,
		})
	}

	result := database.DB.Create(&user)

	if result.Error != nil {
		errMsg := result.Error.Error()

		if strings.Contains(errMsg, "duplicate key value violates unique") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"message": "Your username and email must be unique.",
			})
		}

		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"message": errMsg,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created.",
		"user":    user,
	})
}

// Login to the application
func LoginUser(c *fiber.Ctx) error {
	var payload = new(model.UserLoginData)
	err := c.BodyParser(&payload)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	errors := ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": errors,
		})
	}

	var user model.User
	username := strings.ToLower(payload.Username)
	result := database.DB.First(&user, "email = ? or username = ?", username, username)
	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "User account not found.",
		})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid password.",
		})
	}

	token, err := generateJwtToken(user)

	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"message": fmt.Sprintf("Could not generate JWT Token: %v", err),
		})
	}

	max_age, err := strconv.Atoi(os.Getenv("JWT_MAX_AGE"))
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"message": fmt.Sprintf("Could not generate JWT Token: %v", err),
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		MaxAge:   max_age * 60,
		Secure:   false,
		HTTPOnly: true,
		Domain:   os.Getenv("CLIENT_DOMAIN"),
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": token,
	})
}

func GetUserDetails(c *fiber.Ctx) error {
	user := c.Locals("user").(model.UserDetailsData)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{"user": user},
	})
}

func generateJwtToken(user model.User) (token string, err error) {
	secret, expire_in := jwtConfig()
	tokenByte := jwt.New(jwt.SigningMethodHS256)
	now := time.Now().UTC()
	claims := tokenByte.Claims.(jwt.MapClaims)

	claims["sub"] = user.ID
	claims["exp"] = now.Add(expire_in).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	return tokenByte.SignedString([]byte(secret))
}

func jwtConfig() (string, time.Duration) {
	secret := os.Getenv("JWT_SECRET")
	expire_in := os.Getenv("JWT_EXPIRE_IN")

	duration, err := time.ParseDuration(expire_in)
	if err != nil {
		log.Fatal("Error parsing the JWT_EXPIRE_IN duration.")
		duration = time.Hour * 24
	}

	return secret, duration
}
