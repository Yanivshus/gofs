package main

import (
	"os"

	"github.com/Yanivshus/gofs/internal/database"
	"github.com/Yanivshus/gofs/internal/files"
	"github.com/Yanivshus/gofs/internal/web"

	"github.com/gin-gonic/gin"
)

func main() {
	web.Init_web()
	database.InitDb()
	database.GetInstanceDb()

	err := os.Chdir("fs") // change to the directory
	if err != nil {
		panic(err)
	}
	files.Proj_dir, err = files.Get_wd() // need to have project folder saved.
	if err != nil {
		panic(err)
	}

	router := gin.Default()
	router.POST("/chdir", web.HandleChdir)
	router.GET("/files", web.HandleGetFiles)
	router.POST("/signup", web.HandleSignUp)
	router.Run("localhost:6969")
	web.Dtor_web()
}
