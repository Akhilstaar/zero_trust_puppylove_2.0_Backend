package controllers

import (
	"errors"
	"net/http"

	"github.com/Akhilstaar/zero_trust_puppylove_2.0_Backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserFirstLogin(c *gin.Context) {
	// User already authenticated in router.go by gin.HandlerFunc

	// Validate the input format
	info := new(models.TypeUserFirst)
	if err := c.BindJSON(info); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input data format."})
		return
	}

	tempU := models.MailData{}
	tempUser := models.User{}
	tempRecord := Db.Model(&tempUser).Where("id = ?", info.Id).First(&tempU)
	if tempRecord.Error != nil {
		if errors.Is(tempRecord.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusForbidden, gin.H{"error": "User not found !!"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong, Please try again."})
			return
		}
	}
	if tempU.Dirty {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "User already registered"})
		return
	}

	// See U later ;) ... I mean, see U in the zero trust world.
	user := models.User{}
	publicK := Db.Model(&user).Where("pub_k = ?", info.PubKey).First(&user)
	if publicK.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please enter another public key !!"})
		return
	}

	record := Db.Model(&user).Where("id = ? AND auth_c = ?", info.Id, info.AuthCode).First(&user)
	if record.Error != nil {
		if errors.Is(record.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Incorrect AuthCode entered !!"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong, Please try again."})
			return
		}
	}

	// var newuser models.User
	if err := record.Updates(models.User{
		Id:    info.Id,
		Pass:  info.PassHash,
		PubK:  info.PubKey,
		PrivK: info.PrivKey,
		AuthC: " ",
		Data:  info.Data,
		Dirty: true,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong, Please try again."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User Created Successfully."})
}

func SendHeart_Stage1(c *gin.Context) {
	// User already authenticated in router.go by gin.HandlerFunc
	// TODO: ADD a permit variable to control the api call output ie. accept or reject
	info := new(models.Stage1_Hearts)
	if err := c.BindJSON(info); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input data format."})
		return
	}

	userID, _ := c.Get("user_id")
	var user models.User
	record := Db.Model(&user).Where("id = ?", userID).First(&user)
	if record.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong, Please try again."})
		return
	}
	if user.S1submit {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Hearts already sent."})
		return
	}

	newHeart := models.Stage1{
		Id: userID.(string),
		M1: info.M1,
		M2: info.M2,
		M3: info.M3,
		M4: info.M4,
	}

	if err := Db.Create(&newHeart).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	if err := record.Updates(models.User{
		S1submit: true,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong, Please try again."})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "Hearts Sent Successfully !!"})
}

// TODO: Will implement virtualization later ;)
// func Stage1_Virtual(c *gin.Context) {
// 	info := new(models.Stage1_Hearts)
// 	if err := c.BindJSON(info); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong Format"})
// 		return
// 	}

// 	userID, _ := c.Get("user_id")
// 	var user models.User
// 	record := Db.Model(&user).Where("id = ?", userID).First(&user)
// 	if record.Error != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "User does not exist."})
// 		return
// 	}

// 	if user.S1submit {
// 		c.JSON(http.StatusOK, gin.H{"error": "Hearts already sent."})
// 		return
// 	}

// 	jsonData, err := json.Marshal(info.Hearts)
// 	if err != nil {
// 		return
// 	}

// 	if err := record.Updates(models.User{
// 		Data: string(jsonData),
// 	}).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update Data field of User."})
// 		return
// 	}

// 	c.JSON(http.StatusAccepted, gin.H{"message": "Virtual Hearts Sent Successfully !!"})
// }

func getActiveUsers(c *gin.Context, condition string, columnName string) {
	var users []models.User
	fetchUsers := Db.Model(&models.User{}).Where(condition).Pluck(columnName, &users)
	if fetchUsers.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Some Error Occurred"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": users})
}

func GetStage1ActiveUsers(c *gin.Context) {
	getActiveUsers(c, "dirty = true", "id")
}

func GetStage2ActiveUsers(c *gin.Context) {
	getActiveUsers(c, "s1submit = true", "id")
}

func GetStage3ActiveUsers(c *gin.Context) {
	getActiveUsers(c, "s2submit = true", "id")
}

func FetchPublicKeys(c *gin.Context) {
	var publicKeys []models.UserPublicKey
	fetchPublicKey := Db.Model(&models.User{}).Select("id, pub_k").Find(&publicKeys)
	if fetchPublicKey.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Some Error Occurred"})
		return
	}
	c.JSON(http.StatusOK, publicKeys)
}