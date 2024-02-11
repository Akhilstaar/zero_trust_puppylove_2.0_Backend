package main

import (
	"os"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/Akhilstaar/zero_trust_puppylove_2.0_Backend/db"
	"github.com/Akhilstaar/zero_trust_puppylove_2.0_Backend/router"
	"github.com/Akhilstaar/zero_trust_puppylove_2.0_Backend/utils"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}

	var CfgAdminPass = os.Getenv("CFG_ADMIN_PASS")
	Db := db.InitDB()

	utils.Randinit()
	store := cookie.NewStore([]byte(CfgAdminPass))
	r := gin.Default()

	// Local Testing Frontend Running on localhost:3000
	// r.Use(cors.New(cors.Config{AllowCredentials: true, AllowOrigins: []string{"http://localhost:3000"}, AllowHeaders: []string{"content-type", "g-recaptcha-response"}}))

	// Allow all origins
	r.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowHeaders:     []string{"content-type", "g-recaptcha-response"},
	}))

	r.Use(sessions.Sessions("adminsession", store))
	router.PuppyRoute(r, *Db)

	r.Run(":8080")

	// if err := r.Run(config.CfgAddr); err != nil {
	// 	fmt.Println("[Error] " + err.Error())
	// }
}
