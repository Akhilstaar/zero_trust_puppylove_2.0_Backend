package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Akhilstaar/zero_trust_puppylove_2.0_Backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserFirstLogin(c *gin.Context) {
	// User already authenticated in router.go by gin.HandlerFunc

	// Validate the input format
	var info models.TypeUserFirst
	if err := c.BindJSON(&info); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input data format."})
		return
	}

	// Check if user exists and is not dirty
	var tempU models.MailData
	if err := Db.Model(&models.User{}).Where("id = ?", info.Id).First(&tempU).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusForbidden, gin.H{"error": "User not found !!"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong, please try again."})
		}
		return
	}
	if tempU.Dirty {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "User already registered"})
		return
	}

	// Check if public key is unique
	var user models.User
	if err := Db.Model(&user).Where("pub_k = ?", info.PubKey).First(&user).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please enter another public key !!"})
		return
	}

	// Check if the AuthCode is correct
	if err := Db.Model(&user).Where("id = ? AND auth_c = ?", info.Id, info.AuthCode).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Incorrect AuthCode entered !!"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong, please try again."})
		}
		return
	}

	// Update user information
	if err := Db.Model(&user).Updates(models.User{
		Id:     info.Id,
		Pass:   info.PassHash,
		PubK:   info.PubKey,
		PrivK:  info.PrivKey,
		AuthC:  " ",
		S1Data: "FIRST_LOGIN",
		S2Data: "Not_Sent_YET",
		Dirty:  true,
	}).Error; err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong in records, please try again."})
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

	// Use a separate database session for each operation
	// to avoid conflicts with table name "users" specified more than once
	tx := Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var user models.User
	record := tx.Model(&user).Where("id = ?", userID).First(&user)
	if record.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong, Please try again."})
		return
	}
	if user.S1submit {
		tx.Rollback()
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Hearts already sent."})
		return
	}

	// Check if the entry with userID exists
	existingHeart := models.Stage1{}
	if err := tx.Where("id = ?", userID.(string)).First(&existingHeart).Error; err != nil {
		// Entry not found, create a new one
		if gorm.ErrRecordNotFound == err { // Should've written opposite ?  -- this one has better readability ;}
			newHeart := models.Stage1{
				Id:     userID.(string),
				M1:     info.M1,
				M2:     info.M2,
				M3:     info.M3,
				M4:     info.M4,
				Submit: true,
			}

			if err := tx.Create(&newHeart).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{"error": err})
				return
			}
		} else {
			// Other error occurred
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong, Please try again."})
			return
		}
	} else {
		// Entry found, update the existing one
		existingHeart.M1 = info.M1
		existingHeart.M2 = info.M2
		existingHeart.M3 = info.M3
		existingHeart.M4 = info.M4
		existingHeart.Submit = true
		if err := tx.Save(&existingHeart).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong, Please try again."})
			return
		}
	}

	if err := tx.Model(&user).Where("id = ?", userID).Updates(models.User{
		S1Data:   info.S1Data,
		S1submit: true,
	}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong, Please try again."})
		return
	}

	// Commit the transaction if everything is successful
	tx.Commit()

	c.JSON(http.StatusAccepted, gin.H{"message": "Hearts Sent Successfully !!"})
}

// TODO: Will implement virtualization later ;) -- Done
// func Stage1_Virtual(c *gin.Context) {
func SendHeartVirtual_Stage1(c *gin.Context) {
	// User already authenticated in router.go by gin.HandlerFunc
	// TODO: ADD a permit variable to control the api call output ie. accept or reject
	info := new(models.Stage1_Hearts)
	if err := c.BindJSON(info); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input data format."})
		return
	}

	userID, _ := c.Get("user_id")

	// Use a separate database session for each operation
	// to avoid conflicts with table name "users" specified more than once
	tx := Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var user models.User
	record := tx.Model(&user).Where("id = ?", userID).First(&user)
	if record.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong, Please try again."})
		return
	}
	if user.S1submit {
		tx.Rollback()
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Hearts already sent."})
		return
	}

	// Check if the entry with userID exists
	existingHeart := models.Stage1{}
	if err := tx.Where("id = ?", userID.(string)).First(&existingHeart).Error; err != nil {
		// Entry not found, create a new one
		if gorm.ErrRecordNotFound == err { // Should've written opposite ?  -- this one has better readability ;}
			newHeart := models.Stage1{
				Id:     userID.(string),
				M1:     info.M1,
				M2:     info.M2,
				M3:     info.M3,
				M4:     info.M4,
				Submit: false,
			}

			if err := tx.Create(&newHeart).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{"error": err})
				return
			}
		} else {
			// Other error occurred
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong, Please try again."})
			return
		}
	} else {
		// Entry found, update the existing one
		existingHeart.M1 = info.M1
		existingHeart.M2 = info.M2
		existingHeart.M3 = info.M3
		existingHeart.M4 = info.M4

		if err := tx.Save(&existingHeart).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong, Please try again."})
			return
		}
	}

	if err := tx.Model(&user).Where("id = ?", userID).Updates(models.User{
		S1Data: info.S1Data,
	}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong, Please try again."})
		return
	}

	// Commit the transaction if everything is successful
	tx.Commit()

	c.JSON(http.StatusAccepted, gin.H{"message": "Hearts Sent Successfully !!"})
}

func SendHeart_Stage2(c *gin.Context) {
	// User already authenticated in router.go by gin.HandlerFunc
	// TODO: ADD a permit variable to control the api call output ie. accept or reject
	info := new(models.Stage2_Hearts)
	if err := c.BindJSON(info); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input data format."})
		return
	}

	userID, _ := c.Get("user_id")

	// Use a separate database session for each operation
	// to avoid conflicts with table name "users" specified more than once
	tx := Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var user models.User
	record := tx.Model(&user).Where("id = ?", userID).First(&user)
	if record.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong, Please try again."})
		return
	}
	if user.S2submit {
		tx.Rollback()
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Hearts already sent."})
		return
	}

	// Check if the entry with userID exists
	existingHeart := models.Stage2{}
	if err := tx.Where("id = ?", userID.(string)).First(&existingHeart).Error; err != nil {
		// Entry not found, create a new one
		if gorm.ErrRecordNotFound == err { // Should've written opposite ?  -- this one has better readability ;}
			newHeart := models.Stage2{
				Id:     userID.(string),
				Ka1:    info.Ka1,
				Ka2:    info.Ka2,
				Ka3:    info.Ka3,
				Ka4:    info.Ka4,
				Kb1:    info.Kb1,
				Kb2:    info.Kb2,
				Kb3:    info.Kb3,
				Kb4:    info.Kb4,
				Kc1:    info.Kc1,
				Kc2:    info.Kc2,
				Kc3:    info.Kc3,
				Kc4:    info.Kc4,
				Kd1:    info.Kd1,
				Kd2:    info.Kd2,
				Kd3:    info.Kd3,
				Kd4:    info.Kd4,
				Submit: true,
			}

			if err := tx.Create(&newHeart).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{"error": err})
				return
			}
		} else {
			// Other error occurred - maybe due to db
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong, Please try again."})
			return
		}
	} else {
		// Entry found, update the existing one
		existingHeart.Ka1 = info.Ka1
		existingHeart.Ka2 = info.Ka2
		existingHeart.Ka3 = info.Ka3
		existingHeart.Ka4 = info.Ka4
		existingHeart.Kb1 = info.Kb1
		existingHeart.Kb2 = info.Kb2
		existingHeart.Kb3 = info.Kb3
		existingHeart.Kb4 = info.Kb4
		existingHeart.Kc1 = info.Kc1
		existingHeart.Kc2 = info.Kc2
		existingHeart.Kc3 = info.Kc3
		existingHeart.Kc4 = info.Kc4
		existingHeart.Kd1 = info.Kd1
		existingHeart.Kd2 = info.Kd2
		existingHeart.Kd3 = info.Kd3
		existingHeart.Kd4 = info.Kd4
		existingHeart.Submit = true

		if err := tx.Save(&existingHeart).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong, Please try again."})
			return
		}
	}

	if err := tx.Model(&user).Where("id = ?", userID).Updates(models.User{
		S2submit: true,
		S2Data:   info.S2Data,
	}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong, Please try again."})
		return
	}

	tx.Commit()

	c.JSON(http.StatusAccepted, gin.H{"message": "Hearts Sent Successfully !!"})
}

func SendHeartVirtual_Stage2(c *gin.Context) {
	// User already authenticated in router.go by gin.HandlerFunc
	// TODO: ADD a permit variable to control the api call output ie. accept or reject
	info := new(models.Stage2_Hearts)
	if err := c.BindJSON(info); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input data format."})
		return
	}

	userID, _ := c.Get("user_id")

	// Use a separate database session for each operation
	// to avoid conflicts with table name "users" specified more than once
	tx := Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var user models.User
	record := tx.Model(&user).Where("id = ?", userID).First(&user)
	if record.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong, Please try again."})
		return
	}
	if user.S2submit {
		tx.Rollback()
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Hearts already sent."})
		return
	}

	// Check if the entry with userID exists
	existingHeart := models.Stage2{}
	if err := tx.Where("id = ?", userID.(string)).First(&existingHeart).Error; err != nil {
		// Entry not found, create a new one
		if gorm.ErrRecordNotFound == err { // Should've written opposite ?  -- this one has better readability ;}
			newHeart := models.Stage2{
				Id:     userID.(string),
				Ka1:    info.Ka1,
				Ka2:    info.Ka2,
				Ka3:    info.Ka3,
				Ka4:    info.Ka4,
				Kb1:    info.Kb1,
				Kb2:    info.Kb2,
				Kb3:    info.Kb3,
				Kb4:    info.Kb4,
				Kc1:    info.Kc1,
				Kc2:    info.Kc2,
				Kc3:    info.Kc3,
				Kc4:    info.Kc4,
				Kd1:    info.Kd1,
				Kd2:    info.Kd2,
				Kd3:    info.Kd3,
				Kd4:    info.Kd4,
				Submit: false,
			}

			if err := tx.Create(&newHeart).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{"error": err})
				return
			}
		} else {
			// Other error occurred - maybe due to db
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong, Please try again."})
			return
		}
	} else {
		// Entry found, update the existing one
		existingHeart.Ka1 = info.Ka1
		existingHeart.Ka2 = info.Ka2
		existingHeart.Ka3 = info.Ka3
		existingHeart.Ka4 = info.Ka4
		existingHeart.Kb1 = info.Kb1
		existingHeart.Kb2 = info.Kb2
		existingHeart.Kb3 = info.Kb3
		existingHeart.Kb4 = info.Kb4
		existingHeart.Kc1 = info.Kc1
		existingHeart.Kc2 = info.Kc2
		existingHeart.Kc3 = info.Kc3
		existingHeart.Kc4 = info.Kc4
		existingHeart.Kd1 = info.Kd1
		existingHeart.Kd2 = info.Kd2
		existingHeart.Kd3 = info.Kd3
		existingHeart.Kd4 = info.Kd4

		if err := tx.Save(&existingHeart).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong, Please try again."})
			return
		}
	}

	if err := tx.Model(&user).Where("id = ?", userID).Updates(models.User{
		S2Data: info.S2Data,
	}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong, Please try again."})
		return
	}

	tx.Commit()

	c.JSON(http.StatusAccepted, gin.H{"message": "Hearts Sent Successfully !!"})
}

func Send_Cert(c *gin.Context) {
	// User already authenticated in router.go by gin.HandlerFunc
	// TODO: ADD a permit variable to control the api call output ie. accept or reject
	info := new(models.Share_Cert)
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
	if user.Certgiven {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Certificates already sent."})
		return
	}

	newCert := models.Stage3{
		Id:  userID.(string),
		Sk1: info.CertAs,
		Sk2: info.CertBs,
		Sk3: info.CertCs,
		Sk4: info.CertDs,
	}

	if err := Db.Create(&newCert).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	if err := Db.Model(&user).Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong, Please try again."})
		return
	}

	if err := Db.Model(&user).Updates(models.User{
		Certgiven: true,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong, Please try again."})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "Certificates Sent Successfully !!"})
}

// TODO: Do not delete this function
// func getActiveUsers(c *gin.Context, condition string, columnName string) {
// 	// var users []models.User
// 	var res []models.UserPublicKey
// 	fetchUsers := Db.Model(&models.User{}).Where(condition).Pluck(columnName, &res)
// 	if fetchUsers.Error != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Some Error Occurred"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"users": res})
// }

func GetStage1ActiveUsers(c *gin.Context) {
	// getActiveUsers(c, "dirty = true", "id")  // TODO: See if a common function serves the purpose & with efficiency.
	var users []models.User
	var userDB models.User
	Db.Model(userDB).Find(&users)
	var results []string
	for _, user := range users {
		if user.Dirty {
			results = append(results, user.Id)
		}
	}
	c.JSON(http.StatusOK, gin.H{"users": results})
}

func GetStage2ActiveUsers(c *gin.Context) {
	// getActiveUsers(c, "s1submit = true", "id") // TODO: See if a common function serves the purpose & with efficiency.
	var users []models.User
	var userDB models.User
	Db.Model(userDB).Find(&users)
	var results []string
	for _, user := range users {
		if user.S1submit {
			results = append(results, user.Id)
		}
	}
	c.JSON(http.StatusOK, gin.H{"users": results})
}

func GetStage3ActiveUsers(c *gin.Context) {
	// getActiveUsers(c, "s2submit = true", "id") // TODO: See if a common function serves the purpose & with efficiency.
	var users []models.User
	var userDB models.User
	Db.Model(userDB).Find(&users)
	var results []string
	for _, user := range users {
		if user.S2submit {
			results = append(results, user.Id)
		}
	}
	c.JSON(http.StatusOK, gin.H{"users": results})
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

func FetchStage1Keys(c *gin.Context) {
	var stage1Keys []models.Fetch_Stage1
	fetchStage1Keys := Db.Model(&models.Stage1{}).Select("id, m1, m2, m3, m4").Where("submit = ?", true).Find(&stage1Keys)
	if fetchStage1Keys.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Some Error Occurred"})
		return
	}
	c.JSON(http.StatusOK, stage1Keys)
}

func FetchStage2Keys(c *gin.Context) {
	var stage2Keys []models.Fetch_Stage2
	fetchStage2Keys := Db.Model(&models.Stage2{}).Select("id, Ka1, Ka2, Ka3, Ka4, Kb1, Kb2, Kb3, Kb4, Kc1, Kc2, Kc3, Kc4, Kd1, Kd2, Kd3, Kd4").Where("submit = ?", true).Find(&stage2Keys)
	if fetchStage2Keys.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Some Error Occurred"})
		return
	}
	c.JSON(http.StatusOK, stage2Keys)
}

func FetchCerts(c *gin.Context) {
	var stage2Keys []models.Fetch_Cert
	fetchStage2Keys := Db.Model(&models.Stage3{}).Select("id, sk1, sk2, sk3, sk4").Where("submit = ?", true).Find(&stage2Keys)
	if fetchStage2Keys.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Some Error Occurred"})
		return
	}
	c.JSON(http.StatusOK, stage2Keys)
}