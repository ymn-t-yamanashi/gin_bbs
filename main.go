package main

import (
	"gin_bbs/controllers"
	"gin_bbs/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	r := gin.Default()
	store := cookie.NewStore([]byte("secret_key"))
	r.Use(sessions.Sessions("mysession", store))
	r.LoadHTMLGlob("templates/**/*")

	db, _ := gorm.Open(sqlite.Open("bbs.db"), &gorm.Config{})
	db.AutoMigrate(&models.BBS{})
	db.AutoMigrate(&models.Kind{})

	ctrl := &controllers.BBSController{DB: db}
	ctrl_k := &controllers.KindController{DB: db}

	g := r.Group("/bbs")
	{
		g.GET("/", ctrl.Index)
		g.GET("/new", ctrl.New)
		g.POST("", ctrl.Create)
		g.GET("/:id", ctrl.Show)
		g.GET("/:id/edit", ctrl.Edit)
		g.POST("/:id/update", ctrl.Update)
		g.POST("/:id/delete", ctrl.Delete)
	}

	gk := r.Group("/kind")
	{
		gk.GET("/", ctrl_k.Index)
		gk.GET("/new", ctrl_k.New)
		gk.POST("", ctrl_k.Create)
		gk.GET("/:id", ctrl_k.Show)
		gk.GET("/:id/edit", ctrl_k.Edit)
		gk.POST("/:id/update", ctrl_k.Update)
		gk.POST("/:id/delete", ctrl_k.Delete)
	}

	r.Run(":9999")
}
