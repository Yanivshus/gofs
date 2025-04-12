package main

import (
	"strings"

	"github.com/gin-gonic/gin"
)

var chdir_logger *logger
var getfiles_logger *logger

func handleChdir(c *gin.Context) {
	var pathChange dir

	// get json
	// there is an error here
	if err := c.BindJSON(&pathChange); err != nil {
		c.IndentedJSON(500, gin.H{"Error": err.Error()})
		go getfiles_logger.log_str(err.Error(), c.ClientIP())
		return
	}

	// change current directory to given one
	err := change_dir(pathChange.Path)
	if err != nil {
		c.IndentedJSON(500, gin.H{"Error": err.Error()})
		go getfiles_logger.log_str(err.Error(), c.ClientIP())
		return
	}

	var wd string
	wd, err = get_wd()
	if err != nil {
		c.IndentedJSON(500, gin.H{"state": err.Error()})
		go getfiles_logger.log_str(err.Error(), c.ClientIP())
		return
	}

	c.IndentedJSON(200, gin.H{"cwd": wd})
	go getfiles_logger.log_str("wd changed to: "+wd, c.ClientIP())
}

func handleGetFiles(c *gin.Context) {

	data, err := getfiles()
	if err != nil {
		c.IndentedJSON(401, gin.H{"Error": err.Error()})
		go getfiles_logger.log_str(err.Error(), c.ClientIP())
		return
	}

	go getfiles_logger.log_str(strings.Join(data, ","), c.ClientIP())
	c.IndentedJSON(200, gin.H{"files": data})
}

func main() {
	chdir_logger = create_logger("CHDIR", "file.log")
	getfiles_logger = create_logger("GET_FILES", "file.log")

	router := gin.Default()
	router.POST("/chdir", handleChdir)
	router.GET("/files", handleGetFiles)

	router.Run("localhost:6969")

	chdir_logger.destroy_log()
	getfiles_logger.destroy_log()
}
