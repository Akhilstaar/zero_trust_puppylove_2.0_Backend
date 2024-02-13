package router

import (
	"net/http"

	"github.com/Akhilstaar/zero_trust_puppylove_2.0_Backend/controllers"
	"github.com/Akhilstaar/zero_trust_puppylove_2.0_Backend/db"
	"github.com/gin-gonic/gin"
)

func PuppyRoute(r *gin.Engine, db db.PuppyDb) {

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Hello from the other side!")
	})

	// assigning here cuz I'm only importing controllers here, & considering their size better import them here.
	controllers.Db = db

	// Captcha
	captcha := r.Group("/captcha/user")
	{
		captcha.Use(controllers.Captchacheck())
		captcha.POST("/mail/:id", controllers.UserMail)
		captcha.POST("/login", controllers.UserLogin)
	}
	// User administration
	users := r.Group("/users")
	{
		users.POST("/login/first", controllers.UserFirstLogin)
		users.Use(controllers.AuthenticateUser())
		users.GET("/activeusers/stage1", controllers.GetStage1ActiveUsers)
		users.GET("/activeusers/stage2", controllers.GetStage2ActiveUsers)
		users.GET("/activeusers/stage3", controllers.GetStage3ActiveUsers)
		users.GET("/fetchPublicKeys", controllers.FetchPublicKeys)
		users.GET("/fetchStage1Keys", controllers.FetchStage1Keys)
		users.GET("/fetchStage2Keys", controllers.FetchStage2Keys)
		// users.GET("/fetchCerts", controllers.FetchCerts)

		users.POST("/sendheart/stage1", controllers.SendHeart_Stage1)
		users.POST("/sendheart/stage2", controllers.SendHeart_Stage2)
		users.POST("/sendheartvirtual/stage1", controllers.SendHeartVirtual_Stage1)
		users.POST("/sendheartvirtual/stage2", controllers.SendHeartVirtual_Stage2)
		// users.POST("/sendcert", controllers.Send_Cert)

		// users.GET("/fetchall", controllers.FetchHearts)
	}

	// Session administration
	session := r.Group("/session")
	{
		session.POST("/admin/login", controllers.AdminLogin)
		session.GET("/logout", controllers.UserLogout)
	}

	// admin
	admin := r.Group("/admin")
	{
		admin.Use(controllers.AuthenticateAdmin())
		admin.GET("/user/deleteallusers", controllers.DeleteAllUsers)
		admin.POST("/user/new", controllers.AddNewUser)
		admin.POST("/user/delete", controllers.DeleteUser)
		// admin.GET("/publish", controllers.PublishResults)
		// admin.GET("/TogglePermit", controllers.TogglePermit)
	}
	// r.GET("/stats", controllers.GetStats)

}
