package main

import (
	"fmt"
	"fs/gofs_file_server/internal/crypto"
	"fs/gofs_file_server/internal/database"
)

func main() {
	/*web.Init_web()

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

	router.Run("localhost:6969")
	web.Dtor_web()*/

	db := database.ConnectPostgres()
	db.ClosePostgres()

	str, _ := crypto.GenSalt()
	fmt.Println(str)

}
