package main

import (
	"flag"
	"log"

	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	ses "github.com/sourcegraph/go-ses"
)

const (
	errServer             = "Server error. Please contact your system administrator."
	requestInProcess      = "Request in process"
	friendRequestSent     = "Friend Request Sent"
	friendRequestApproved = "Friend Request Approved"
)

var (
	gormDatabase *gorm.DB

	dbname     = Config().GetString("dbname")
	dbpassword = Config().GetString("dbpassword")
	dbuser     = Config().GetString("dbuser")

	sessionStore = sessions.NewCookieStore([]byte(Config().GetString("sessionKey")))

	sesConfig = ses.Config{
		Endpoint:        Config().GetString("sesEndpoint"),
		AccessKeyID:     Config().GetString("sesAccessKeyID"),
		SecretAccessKey: Config().GetString("sesSecretAccessKey"),
	}
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
	gormDatabase.AutoMigrate(&Friendship{})
}
