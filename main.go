package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func handleChdir(c *gin.Context) {
	var pathChange dir

	// get json
	// there is an error here
	if err := c.BindJSON(&pathChange); err != nil {
		c.IndentedJSON(500, gin.H{"Error": err.Error()})
		return
	}

	// change current directory to given one
	err := change_dir(pathChange.Path)
	if err != nil {
		c.IndentedJSON(500, gin.H{"Error": err.Error()})
		return
	}

	var wd string
	wd, err = get_wd()
	if err != nil {
		c.IndentedJSON(500, gin.H{"state": "err"})
		return
	}

	c.IndentedJSON(200, gin.H{"cwd": wd})
}

func handleGetFiles(c *gin.Context) {
	data, err := getfiles()
	if err != nil {
		c.IndentedJSON(401, gin.H{"Error": err.Error()})
		return
	}

	c.IndentedJSON(200, gin.H{"files": data})
}

func main() {

	change_dir("../")
	data, _ := get_wd()
	fmt.Print(data)
	router := gin.Default()
	router.POST("/chdir", handleChdir)
	router.GET("/files", handleGetFiles)

	router.Run("localhost:6969")

}
