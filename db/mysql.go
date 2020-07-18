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

func (impl *ImplDB) Find(identifier model.Identifier) User {
	user := User{}
	impl.DB.Find(&user, "email = ? AND password = ?", identifier.Email, identifier.Password)
	return user
}

func (impl *ImplDB) FindByToken(token string) User {
	user := User{}
	impl.DB.Find(&user, "access_token = ?", token)
	fmt.Println("取得したuserの値は: ", user)
	return user
}

func (impl *ImplDB) UpdateLoggedInAt(user User) User {
	impl.DB.Model(&user).Update("logged_in_at", time.Now())
	return user
}

func (impl *ImplDB) Create(user User) {
	impl.DB.Create(&user)
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

//func (impl *ImplDB) dummy() {
//	insertUser := User{}
//	user := User{}
//	user2 := User{}
//
//	insertUser.UserID = 1
//	insertUser.UserName = "hoge"
//
//	// 作成
//	// INSERT INTO users(user_id,user_name) VALUES(1,'hoge');
//	impl.DB.Create(&insertUser)
//
//	// 取得
//	// SELECT * FROM users WHERE user_id = 1;
//	impl.DB.Find(&user, "user_id = ?", 1)
//
//	fmt.Println("取得したuserの値は", user)
//
//	// 更新
//	// UPDATE users SET user_name = 'fuga' WHERE user_id = 1 and user_name = 'hoge';
//	impl.DB.Model(&user).Update("user_name", "fuga")
//
//	fmt.Println("更新後のuser:", user)
//
//	// 削除
//	// DELETE FROM users WHERE user_id = 1 and user_name = 'hoge';
//	impl.DB.Delete(&user)
//
//	if err := impl.DB.Find(&user2, "user_id = ?", 1).Error; err != nil {
//		// エラーハンドリング
//		fmt.Println("存在しませんでした...")
//	} else {
//		fmt.Println("取得したuserの値は", user2)
//	}
//}