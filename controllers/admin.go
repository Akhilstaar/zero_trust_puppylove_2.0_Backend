package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/Akhilstaar/zero_trust_puppylove_2.0_Backend/db"
	"github.com/Akhilstaar/zero_trust_puppylove_2.0_Backend/models"
	"github.com/Akhilstaar/zero_trust_puppylove_2.0_Backend/utils"
	"github.com/gin-gonic/gin"
)

var Db db.PuppyDb
var permit bool = true

func AdminLogin(c *gin.Context) {
	info := new(models.AdminLogin)
	if err := c.BindJSON(info); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input data format."})
		return
	}

	if info.Id != os.Getenv("ADMIN_ID") {
		c.JSON(http.StatusForbidden, gin.H{"error": "This action will be reported."})
		return
	}

	if info.Pass != os.Getenv("ADMIN_PASS") {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid Password."})
		return
	}

	token, err := generateJWTToken(info.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT token"})
		return
	}
	expirationTime := time.Now().Add(time.Hour * 24)
	cookie := &http.Cookie{
		Name:     "Authorization",
		Value:    token,
		Expires:  expirationTime,
		Path:     "/",
		Domain:   os.Getenv("DOMAIN"),
		HttpOnly: true,
		Secure:   false, // Set this to true if you're using HTTPS, false for HTTP
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(c.Writer, cookie)
	c.JSON(http.StatusOK, gin.H{"message": "Admin logged in successfully !!"})
}

func AddNewUser(c *gin.Context) {
	// TODO: Modify this function to handle multiple concatenated json inputs -- Done

	// TODO: Implement admin authentication logic -- Done

	// Validate the input format
	info := new(models.AddNewUser)
	if err := c.BindJSON(info); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input data format."})
		return
	}

	// Create user
	for _, user := range info.TypeUserNew {

		newUser := models.User{
			Id:        user.Id,
			Name:      user.Name,
			Email:     user.Email,
			Pass:      "",
			PubK:      "",
			PrivK:     "",
			AuthC:     utils.RandStringRunes(15),
			Data:      "",
			S1submit:  false,
			S2submit:  false,
			Dirty:     false,
			Certgiven: false,
		}

		// Insert the user into the database

		if err := Db.Create(&newUser).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully."})
}

func DeleteUser(c *gin.Context) {
	// TODO: Implement admin authentication logic -- Done

	// Validate the input format
	info := new(models.TypeUserNew)
	if err := c.BindJSON(info); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input data format."})
		return
	}

	newUser := models.User{
		Id:     info.Id,
		Name:   info.Name,
		Email:  info.Email,
	}

	if err := Db.Unscoped().Delete(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User Deleted successfully."})
}

func DeleteAllUsers(c *gin.Context) {
	// TODO: Implement admin authentication logic -- Done

	newUser := models.User{}
	if err := Db.Unscoped().Where("1 = 1").Delete(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "All Users Deleted successfully."})
}
