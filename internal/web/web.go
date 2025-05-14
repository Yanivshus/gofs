package web

import (
	"encoding/json"
	"fs/gofs_file_server/internal/database"
	"fs/gofs_file_server/internal/files"

	"fs/gofs_file_server/internal/crypto"

	"github.com/gin-gonic/gin"
)

var chdir_logger *files.Logger
var getfiles_logger *files.Logger

func Init_web() {
	// start logger rutines
	chdir_logger = files.Create_logger("CHDIR", "file.log")
	getfiles_logger = files.Create_logger("GET_FILES", "file.log")
	go chdir_logger.Keep_logger()
	go getfiles_logger.Keep_logger()
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

func HandleGetFiles(c *gin.Context) {

	data, err := files.Getfiles()
	if err != nil {
		go getfiles_logger.Log_str(err.Error(), c.ClientIP())
		c.IndentedJSON(401, gin.H{"Error": err.Error()})
		return
	}

	go getfiles_logger.Log_str(data, c.ClientIP())
	c.IndentedJSON(200, gin.H{"files": json.RawMessage(data)})
}

func HandleLogin(c *gin.Context) {

}

func HandleSignUp(c *gin.Context) {
	var user database.User
	if err := c.BindJSON(&user); err != nil {
		c.IndentedJSON(500, gin.H{"error": "internal error"})
		return
	}

	db := database.GetInstanceDb()
	if err := db.DoesUserExists(user); err != nil {
		c.IndentedJSON(500, gin.H{"error": "user already exists"})
		return
	}

	salt, err := crypto.GenSalt()
	if err != nil {
		c.IndentedJSON(500, gin.H{"error": "internal error"})
		return
	}
	
	
}
