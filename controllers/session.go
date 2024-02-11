package controllers

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/Akhilstaar/zero_trust_puppylove_2.0_Backend/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

func UserLogin(c *gin.Context) {
	info := new(models.UserLogin)
	if err := c.BindJSON(info); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input data format."})
		return
	}

	loginmodel := models.User{}
	verifyuser := Db.Model(&loginmodel).Where("id = ? AND pass = ?", info.Id, info.Pass).First(&loginmodel)
	if verifyuser.Error != nil {
		if errors.Is(verifyuser.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid Login Request."})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": verifyuser.Error.Error()})
			return
		}
	}

	token, err := generateJWTToken(info.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT token"})
		return
	}
	expirationTime := time.Now().Add(time.Hour * 24)
	cookie := &http.Cookie{
		Name:    "Authorization",
		Value:   token,
		Expires: expirationTime,
		Path:    "/",
		Domain:  os.Getenv("DOMAIN"),
		// For Http
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		// For Https
		// HttpOnly: false,
		// Secure:   true, // Set this to true if you're using HTTPS, false for HTTP
		// SameSite: http.SameSiteNoneMode,
	}

	http.SetCookie(c.Writer, cookie)
	c.JSON(http.StatusOK, gin.H{"message": "User logged in successfully !!",
		"pubKey":     loginmodel.PubK,
		"pvtKey_Enc": loginmodel.PrivK,
		"data":       loginmodel.Data,
		"s1submit":   loginmodel.S1submit,
		"s2submit":   loginmodel.S2submit,
		"certgiven":  loginmodel.Certgiven,
	})
}

func UserLogout(c *gin.Context) {
	cookie := &http.Cookie{
		Name:     "Authorization",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour), // Expire the cookie immediately
		Path:     "/",
		Domain:   os.Getenv("DOMAIN"),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteNoneMode,
		// For Htpps
		// HttpOnly: false,
		// Secure:   true, // Set this to true if you're using HTTPS, false for HTTP
		// SameSite: http.SameSiteNoneMode,
	}

	http.SetCookie(c.Writer, cookie)

	c.JSON(http.StatusOK, gin.H{
		"message": "User logged out successfully.",
	})
}

type AuthClaims struct {
	User_id string `json:"user_id"`
	jwt.StandardClaims
}

func generateJWTToken(userID string) (string, error) {
	var jwtSigningKey = os.Getenv("USER_JWT_SIGNING_KEY")
	claims := AuthClaims{
		userID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 12).Unix(),
			IssuedAt:  jwt.TimeFunc().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSigningKey))
	return tokenString, err
}