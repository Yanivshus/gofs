package main

import (
	"encoding/json"
	"os"
	"fs/gofs_file_server/internal/files"
	"fs/gofs_file_server/internal/auth"
	"github.com/gin-gonic/gin"
)

var chdir_logger *files.Logger
var getfiles_logger *files.Logger

func handleChdir(c *gin.Context) {
	var pathChange files.Dir

	// get json
	// there is an error here
	if err := c.BindJSON(&pathChange); err != nil {
		go getfiles_logger.Log_str(err.Error(), c.ClientIP())
		c.IndentedJSON(500, gin.H{"Error": err.Error()})
		return
	}

	// change current directory to given one
	err := files.Change_Dir(pathChange.Path)
	if err != nil {
		go getfiles_logger.Log_str(err.Error(), c.ClientIP())
		c.IndentedJSON(500, gin.H{"Error": err.Error()})
		return
	}

	var wd string
	wd, err = files.Get_wd()
	if err != nil {
		go getfiles_logger.Log_str(err.Error(), c.ClientIP())
		c.IndentedJSON(500, gin.H{"state": err.Error()})
		return
	}

	go getfiles_logger.Log_str("wd changed to: "+wd, c.ClientIP())
	c.IndentedJSON(200, gin.H{"cwd": wd})
}

func handleGetFiles(c *gin.Context) {

	data, err := files.Getfiles()
	if err != nil {
		go getfiles_logger.Log_str(err.Error(), c.ClientIP())
		c.IndentedJSON(401, gin.H{"Error": err.Error()})
		return
	}

	go getfiles_logger.Log_str(data, c.ClientIP())
	c.IndentedJSON(200, gin.H{"files": json.RawMessage(data)})
}

func main() {
	// start logger rutines
	chdir_logger = files.Create_logger("CHDIR", "file.log")
	getfiles_logger = files.Create_logger("GET_FILES", "file.log")
	go chdir_logger.Keep_logger()
	go getfiles_logger.Keep_logger()

	err := os.Chdir("fs") // change to the directory
	if err != nil {
		panic(err)
	}
	files.Proj_dir, err = files.Get_wd() // need to have project folder saved.
	if err != nil {
		panic(err)
	}

	router := gin.Default()
	router.POST("/chdir", handleChdir)
	router.GET("/files", handleGetFiles)

	router.Run("localhost:6969")

	chdir_logger.Destroy_log()
	getfiles_logger.Destroy_log()
}
