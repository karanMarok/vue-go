package controllers

import (
	"backend/db"
	"backend/helpers"
	"backend/middlewares"
	"backend/models"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User

		//Binding the instantiated struct with request body
		err := c.ShouldBindJSON(&user)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//Get values from database and Check if email already exists
		var existingUser models.User

		if err := db.Db.Where("email = ?", user.Email).First(&existingUser).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Hash the password before saving
				hashedPwd := helpers.HashPassword([]byte(user.Password))
				user.Password = hashedPwd

				// Create the user
				if err := db.Db.Create(&user).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				c.JSON(http.StatusOK, gin.H{"message": "User registered successfully", "user": user})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
		} else {
			// User with the same email already exists
			c.JSON(http.StatusConflict, gin.H{"error": "User with this email already exists"})
		}
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User

		//Get data from body
		err := c.ShouldBindJSON(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//Check if user exists and get all data -- saved in existingUser
		var existingUser models.User

		if err := db.Db.Where("email=?", user.Email).Find(&existingUser).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//Compare passwords
		isPasswordSame := helpers.ComparePassword([]byte(existingUser.Password), []byte(user.Password))

		//If passwords match
		if !isPasswordSame {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password does not match!"})
		} else {
			//Generate both tokens and add them to user table
			accessToken, refreshToken, err := middlewares.GenerateToken(existingUser.Email)

			if err != nil {
				// Log the error or handle it appropriately
				log.Printf("Error generating tokens: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
				return
			}

			//Update tokens to database
			existingUser.AccessToken = accessToken
			existingUser.RefreshToken = refreshToken

			//Save changes to database
			if err := db.Db.Save(&existingUser).Error; err != nil {
				fmt.Println(err.Error(), "sbdsgdjsvjdvsj")
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save tokens"})
				return
			}

			//Set cookies
			http.SetCookie(c.Writer, &http.Cookie{
				Name:     "access-token",
				Value:    accessToken,
				Expires:  time.Now().Add(50 * time.Minute),
				HttpOnly: true,
				Secure:   true,
			})

			http.SetCookie(c.Writer, &http.Cookie{
				Name:     "refresh-token",
				Value:    refreshToken,
				Expires:  time.Now().Add(50 * time.Minute),
				HttpOnly: true,
				Secure:   true,
			})

			//Send message of successful login
			c.JSON(http.StatusOK, gin.H{"message": "Successfully Logged In!"})
		}
	}
}

func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.SetCookie("access-token", "", -1, "/", "", false, true)
		c.SetCookie("refresh-token", "", -1, "/", "", false, true)
		c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
	}
}
