package web

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	b64 "encoding/base64"
	"github.com/Yanivshus/gofs/internal/database"
	"github.com/Yanivshus/gofs/internal/files"

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

	err := files.CreateDirIfNeeded("assets")
	if err != nil {
		panic(err)
	}

	err = files.CreateDirIfNeeded("fs")
	if err != nil {
		panic(err)
	}
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
	file, _, err := c.Request.FormFile("upload")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	// copy file bytes to buffer , later will be saved as blob.
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		c.IndentedJSON(500, gin.H{"error": "internal error1"})

	}

	userData := c.Request.FormValue("data")
	var user database.User
	err = json.Unmarshal([]byte(userData), &user)
	if err != nil {
		c.IndentedJSON(500, gin.H{"error": "json conversion err"}) // for now
		return
	}

	db := database.GetInstanceDb()
	if res, err := db.DoesUserExists(user, "pending"); err != nil || res {
		c.IndentedJSON(500, gin.H{"error": "user waiting for admin approval"})
		return
	}

	stringPhoto := b64.StdEncoding.EncodeToString(buf.Bytes())
	err = db.InsertImageByUsername(stringPhoto, user.Username)
	if err != nil {
		c.IndentedJSON(500, gin.H{"error": "rpoblem inserting image"})
		
	}

	err = db.AddUserToPending(user)
	if err != nil {
		c.IndentedJSON(500, gin.H{"error": "problem inserting user to pending"})
		return
	}

}
