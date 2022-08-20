package modeltests

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/jumkos/WartaGkjMedari/api/controllers"
	"github.com/jumkos/WartaGkjMedari/api/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var server = controllers.Server{}
var userInstance = models.User{}
var renunganInstance = models.Renungan{}

func TestMain(m *testing.M) {

	err := godotenv.Load(os.ExpandEnv("../../.env"))
	if err != nil {
		log.Fatalf("Error getting env %v\n", err)
	}
	Database()

	os.Exit(m.Run())
}

func Database() {

	var err error

	DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("TestDbUser"), os.Getenv("TestDbPassword"), os.Getenv("TestDbHost"), os.Getenv("TestDbPort"), os.Getenv("TestDbName"))
	server.DB, err = gorm.Open(mysql.Open(DBURL))
	if err != nil {
		fmt.Println("Cannot connect to database")
		log.Fatal("This is the error:", err)
	} else {
		fmt.Println("We are connected to the database")
	}
}

func refreshUserTable() error {
	server.DB.Migrator().DropTable(&models.User{})
	server.DB.AutoMigrate(&models.User{})
	log.Printf("Successfully refreshed table")
	return nil
}

func seedOneUser() (models.User, error) {

	refreshUserTable()

	user := models.User{
		Nickname: "Pet",
		Email:    "pet@gmail.com",
		Password: "password",
	}

	err := server.DB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		log.Fatalf("cannot seed users table: %v", err)
	}
	return user, nil
}

func seedUsers() error {

	users := []models.User{
		{
			Nickname: "Steven victor",
			Email:    "steven@gmail.com",
			Password: "password",
		},
		{
			Nickname: "Kenny Morris",
			Email:    "kenny@gmail.com",
			Password: "password",
		},
	}

	for i := range users {
		err := server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func refreshUserAndRenunganTable() error {

	server.DB.Migrator().DropTable(&models.User{}, &models.Renungan{})
	server.DB.AutoMigrate(&models.User{}, &models.Renungan{})
	log.Printf("Successfully refreshed tables")
	return nil
}

func seedOneUserAndOneRenungan() (models.Renungan, error) {

	err := refreshUserAndRenunganTable()
	if err != nil {
		return models.Renungan{}, err
	}
	user := models.User{
		Nickname: "Sam Phil",
		Email:    "sam@gmail.com",
		Password: "password",
	}
	err = server.DB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		return models.Renungan{}, err
	}
	post := models.Renungan{
		Title:    "This is the title sam",
		Content:  "This is the content sam",
		AuthorID: user.ID,
	}
	err = server.DB.Model(&models.Renungan{}).Create(&post).Error
	if err != nil {
		return models.Renungan{}, err
	}
	return post, nil
}

func seedUsersAndRenungan() ([]models.User, []models.Renungan, error) {

	var err error

	if err != nil {
		return []models.User{}, []models.Renungan{}, err
	}
	var users = []models.User{
		{
			Nickname: "Steven victor",
			Email:    "steven@gmail.com",
			Password: "password",
		},
		{
			Nickname: "Magu Frank",
			Email:    "magu@gmail.com",
			Password: "password",
		},
	}
	var posts = []models.Renungan{
		{
			Title:   "Title 1",
			Content: "Hello world 1",
		},
		{
			Title:   "Title 2",
			Content: "Hello world 2",
		},
	}

	for i := range users {
		err = server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		posts[i].AuthorID = users[i].ID

		err = server.DB.Model(&models.Renungan{}).Create(&posts[i]).Error
		if err != nil {
			log.Fatalf("cannot seed posts table: %v", err)
		}
	}
	return users, posts, nil
}
