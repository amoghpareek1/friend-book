package main

import (
	"flag"
	"log"

	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const (
	errServer   = "Server error. Please contact your system administrator."
	accessToken = "token"
)

var (
	gormDatabase *gorm.DB

	dbname     = Config().GetString("dbname")
	dbpassword = Config().GetString("dbpassword")
	dbuser     = Config().GetString("dbuser")

	jwtKey = []byte(Config().GetString("jwtkey"))

	sessionStore = sessions.NewCookieStore([]byte(Config().GetString("sessionKey")))
)

func init() {
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)

	flag.String("port", "", "port to listen on")

	flag.Parse()

	if flag.Lookup("port").Value.String() == "" {
		log.Fatal("-port is required.")
	}

	var err error
	gormDatabase, err = gorm.Open("mysql", dbuser+":"+dbpassword+"@/"+dbname+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal(err)
	}

	gormDatabase.AutoMigrate(&User{})
}
