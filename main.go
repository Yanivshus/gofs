package main

import (
	"github.com/gin-gonic/gin"
)

func handleChdir(c *gin.Context) {
	var pathChange dir

	// get json
	if err := c.BindJSON(&pathChange); err != nil {
		c.IndentedJSON(401, gin.H{"Error": err.Error()})
		return
	}

	// change current directory to given one
	err := change_dir(pathChange.path)
	if err != nil {
		c.IndentedJSON(401, gin.H{"Error": err.Error()})
		return
	}

	var wd string
	wd, err = get_wd()
	if err != nil {
		c.IndentedJSON(200, gin.H{"state": "err"})
		return
	}

	c.IndentedJSON(200, gin.H{"cwd": wd})
}

func main() {
	router := gin.Default()
	router.POST("/chdir", handleChdir)

	router.Run("localhost:6969")

}
