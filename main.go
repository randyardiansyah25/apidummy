package main

import (
	"github.com/joho/godotenv"
	"github.com/kpango/glg"
	"os"
)

func main() {
	_ = glg.Log(os.Getenv("application.name"))
	_ = glg.Log(os.Getenv("application.desc"))
	_ = glg.Log(os.Getenv("application.ver"))
	if err := StartServer(); err != nil{
		_=glg.Log(err.Error())
	}
}

func init() {
	log := glg.FileWriter("log/xroute.log", 0666)
	glg.Get().
		SetMode(glg.BOTH).
		AddLevelWriter(glg.LOG, log).
		AddLevelWriter(glg.ERR, log).
		AddLevelWriter(glg.WARN, log).
		AddLevelWriter(glg.DEBG, log).
		AddLevelWriter(glg.INFO, log)

	if err:=godotenv.Load("config.env"); err != nil {
		glg.Fatalln(err.Error())
	}

	DbConn.Connect(
		os.Getenv("database.addr"),
		os.Getenv("database.port"),
		os.Getenv("database.user"),
		os.Getenv("database.pass"),
		os.Getenv("database.name"),
	)

}
