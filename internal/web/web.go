package web

import (
	"encoding/json"
	"fs/gofs_file_server/internal/database"
	"fs/gofs_file_server/internal/files"

	"github.com/gin-gonic/gin"
)

var chdir_logger *files.Logger
var getfiles_logger *files.Logger

func Init_web() {
	// start logger rutines
	chdir_logger = files.CreateLogger("CHDIR", "file.log")
	getfiles_logger = files.CreateLogger("GET_FILES", "file.log")
	go chdir_logger.KeepLogger()
	go getfiles_logger.KeepLogger()
}

func Dtor_web() {
	chdir_logger.DestroyLog()
	getfiles_logger.DestroyLog()
}

func HandleChdir(c *gin.Context) {
	var pathChange files.Dir

	// get json
	// there is an error here
	if err := c.BindJSON(&pathChange); err != nil {
		go getfiles_logger.LogStr(err.Error(), c.ClientIP())
		c.IndentedJSON(500, gin.H{"Error": err.Error()})
		return
	}

	// change current directory to given one
	err := files.Change_Dir(pathChange.Path)
	if err != nil {
		go getfiles_logger.LogStr(err.Error(), c.ClientIP())
		c.IndentedJSON(500, gin.H{"Error": err.Error()})
		return
	}

	var wd string
	wd, err = files.Get_wd()
	if err != nil {
		go getfiles_logger.LogStr(err.Error(), c.ClientIP())
		c.IndentedJSON(500, gin.H{"state": err.Error()})
		return
	}

	go getfiles_logger.LogStr("wd changed to: "+wd, c.ClientIP())
	c.IndentedJSON(200, gin.H{"cwd": wd})
}

func HandleGetFiles(c *gin.Context) {

	data, err := files.Getfiles()
	if err != nil {
		go getfiles_logger.LogStr(err.Error(), c.ClientIP())
		c.IndentedJSON(401, gin.H{"Error": err.Error()})
		return
	}

	go getfiles_logger.LogStr(data, c.ClientIP())
	c.IndentedJSON(200, gin.H{"files": json.RawMessage(data)})
}

func HandleLogin(c *gin.Context) {

}

// TODO : dont allow to see api gateway responses.
func HandleSignUp(c *gin.Context) {
	var user database.User
	if err := c.BindJSON(&user); err != nil {
		c.IndentedJSON(500, gin.H{"error": "internal error1"})
		return
	}

	db := database.GetInstanceDb()
	if res, err := db.DoesUserExists(user, "pending"); err != nil || res {
		c.IndentedJSON(500, gin.H{"error": "user already exists"})
		return
	}

	err := db.AddUserToPending(user)
	if err != nil {
		c.IndentedJSON(500, gin.H{"error": "problem inserting user to pending"})
		return
	}

}
