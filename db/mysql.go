package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"mock-api/model"
	"time"
)

type User struct {
	AccessToken string `gorm:"size:255;primary_key"`
	UserID      string `gorm:"size:255;primary_key"`
	Email       string `gorm:"size:255"`
	Password    string `gorm:"size:255"`
	CreatedAt   time.Time
	LoggedInAt	time.Time
}

type ImplDB struct {
	DB *gorm.DB
}

func (impl *ImplDB) Start() {
	impl.connectDB()
	impl.initMigration()
}

func (impl *ImplDB) Find(identifier model.Identifier) (User, error) {
	user := User{}
	error := impl.DB.Find(&user, "email = ? AND password = ?", identifier.Email, identifier.Password).Error
	return user, error
}

func (impl *ImplDB) FindByToken(token string) (User, error) {
	user := User{}
	error := impl.DB.Find(&user, "access_token = ?", token).Error
	return user, error
}

func (impl *ImplDB) FindByUserId(userId string) (User, error) {
	user := User{}
	error := impl.DB.Find(&user, "user_id = ?", userId).Error
	return user, error
}

func (impl *ImplDB) UpdateLoggedInAt(user User) User {
	impl.DB.Model(&user).Update("logged_in_at", time.Now())
	return user
}

func (impl *ImplDB) Create(user User) {
	impl.DB.Create(&user)
}

func (impl *ImplDB) DeleteByToken(token string) error {
	user, error := impl.FindByToken(token)
	if error != nil {
		return error
	}
	return impl.DB.Delete(&user).Error
}

func (impl *ImplDB) Close() {
	impl.DB.Close()
}

func (impl *ImplDB) connectDB() {
	dbms := "mysql"
	user := "root"
	password := "password"
	protocol := "tcp(dummy-mysql:3306)"
	dbname := "mysql"

	connect := user + ":" + password + "@" + protocol + "/" + dbname + "?parseTime=true"
	var error error
	impl.DB, error = gorm.Open(dbms, connect)
	if error != nil {
		fmt.Println("Retry connect DB...")
		time.Sleep(time.Second * 30)
		impl.DB, error = gorm.Open(dbms, connect)
		if error != nil {
			panic(error.Error())
		}
	}
}

func (impl *ImplDB) initMigration() {
	impl.DB.AutoMigrate(&User{})
}
