package main

import (
	"github.com/brianvoe/gofakeit"
	"github.com/gin-gonic/gin"
	"github.com/kpango/glg"
	"io/ioutil"
	"os"
	"time"
)

func StartServer() error {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard
	router := gin.Default()
	router.Use(gin.Recovery())
	RegisterHandler(router)
	port := ":" + os.Getenv("application.port")
	_ = glg.Log("Listening at,", port)
	return router.Run(port)
}

func RegisterHandler(router *gin.Engine) {
	router.POST("/", home)
	router.GET("/", home)
	router.GET("/init", initTable)
	router.GET("/addnew", addNew)
	router.GET("/list", list)
}

func home(c *gin.Context) {

	c.JSON(200, gin.H{
		"app.name":        os.Getenv("application.name"),
		"app.desc":        os.Getenv("application.desc"),
		"app.ver":         os.Getenv("application.ver"),
		"port.listener":   os.Getenv("application.port"),
		"database.status": DbConn.Status,
	})
}

func responseError(c *gin.Context, err error) {
	c.JSON(200, gin.H{
		"status": "error",
		"desc":   err.Error(),
	})
}

func responseSuccess(c *gin.Context, data interface{}) {
	c.JSON(200, gin.H{
		"status": "success",
		"data":   data,
	})
}

func addNew(c *gin.Context) {
	gofakeit.Seed(time.Now().UnixNano())
	person := gofakeit.Person()
	employee := Employee{
		EmployeeId: gofakeit.Int32(),
		FirstName:  person.FirstName,
		LastName:   person.LastName,
		Address:    person.Address.Address,
		Phone:      gofakeit.Phone(),
	}

	conn := DbConn.GetConnection()
	_, err := conn.Exec("INSERT INTO employee "+
		"values(?,?,?,?,?)", employee.EmployeeId, employee.FirstName, employee.LastName, employee.Address, employee.Phone)
	if err != nil {
		responseError(c, err)
		return
	}
	responseSuccess(c, employee)
}

func list(c *gin.Context) {
	conn := DbConn.GetConnection()
	rows, err := conn.Query("SELECT " +
		"employee_id, first_name, last_name, address, phone from employee")
	if err != nil {
		responseError(c, err)
		return
	}
	defer func() {
		_ = rows.Close()
	}()

	employess := make([]Employee, 0)
	for rows.Next() {
		employee := Employee{}
		if err = rows.Scan(&employee.EmployeeId, &employee.FirstName, &employee.LastName, &employee.Address, &employee.Phone); err != nil {
			responseError(c, err)
			return
		}
		employess = append(employess, employee)
	}
	responseSuccess(c, gin.H{
		"count":   len(employess),
		"records": employess,
	})
}
func initTable(c *gin.Context) {
	db := Database{}
	db.Connect(
		os.Getenv("database.addr"),
		os.Getenv("database.port"),
		os.Getenv("database.user"),
		os.Getenv("database.pass"),
		"",
	)
	conn := db.GetConnection()
	defer func() {
		_ = conn.Close()
	}()

	if _, err := conn.Exec("CREATE DATABASE " + os.Getenv("database.name")); err != nil {
		responseError(c, err)
		return
	}

	if _, err := conn.Exec("USE " + os.Getenv("database.name")); err != nil {
		responseError(c, err)
		return
	}

	stmt, err := conn.Prepare("CREATE TABLE employee (" +
		"employee_id int NOT NULL," +
		"first_name varchar(120)," +
		"last_name varchar(120)," +
		"address varchar(300)," +
		"phone varchar(25)," +
		"PRIMARY KEY(employee_id)" +
		");")
	if err != nil {
		responseError(c, err)
		return
	}

	_, err = stmt.Exec()
	if err != nil {
		responseError(c, err)
		return
	}

	DbConn.Connect(
		os.Getenv("database.addr"),
		os.Getenv("database.port"),
		os.Getenv("database.user"),
		os.Getenv("database.pass"),
		os.Getenv("database.name"),
	)

	c.JSON(200, gin.H{
		"status": "success",
		"desc":   "create table success!",
	})
}
