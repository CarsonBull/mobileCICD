package pipelines

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"os"
	"time"
)

var DB *gorm.DB

func initializeDB() {

	var err error

	fmt.Println("Waiting for DB")
	duration := time.Duration(1000 * 5 )
	time.Sleep(duration* time.Millisecond)
	fmt.Println("Trying to connect to DB")

	dbHost := os.Getenv("DB_URL")

	dbPort := os.Getenv("DB_PORT")

	ssl := os.Getenv("SSL")

	DB, err = gorm.Open("postgres", "host=" + dbHost + " port=" + dbPort + " user=test password=password sslmode=" + ssl)

	if err != nil {
		fmt.Println(err)
		panic("Unable to open database")
	}

	DB.AutoMigrate(&PipelineObject{})
}

func init() {
	initializeDB()
}
