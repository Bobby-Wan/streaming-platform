package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bobby-wan/streaming-platform/configure"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type ContentCategory uint32

const (
	Gaming ContentCategory = 1 << iota
	IRL
	Music
	Art
)

func (c ContentCategory) String() string {
	switch c {
	case Gaming:
		return "Gaming"
	case IRL:
		return "IRL"
	case Music:
		return "Music"
	case Art:
		return "Art"
	default:
		return fmt.Sprintf("%d", int(c))
	}
}

type Role uint16

const (
	Normal Role = 1 << iota
	Moderator
	Admin
)

type User struct {
	gorm.Model
	Email      string `gorm:"unique"`
	Username   string `gorm:"unique"`
	Password   []byte
	LastActive time.Time
	Role       Role
}

type ActiveStream struct {
	gorm.Model
	Username string `gorm:unique` //cannot have one user with two streams at the same time
	Title    string
	Viewers  uint
	Category uint32
	URL      string
	Active   bool
}

// type ActiveStreamCategories struct {
// 	Stream         ActiveStream
// 	ActiveStreamID uint
// 	Category       ContentCategory
// }

type UserControllerInterface interface {
	Get(username string) (*User, error)
	GetByEmail(email string) (*User, error)
	Create(user User) (*User, error)
	Delete(id uint) error
	Update(user User) error
}

type StreamControllerInterface interface {
	Create(stream ActiveStream) (*ActiveStream, error)
	// Get(id uint) (*ActiveStream, error)
	GetByUsername(username string) (*ActiveStream, error)
	GetByCategories(mask uint32) ([]ActiveStream, error)
	GetByViewCount(limit int) ([]ActiveStream, error)
	Delete(userID uint) error
	End(username string) (*ActiveStream, error)
	Start(username string) (*ActiveStream, error)
}

type LoginControllerInterface struct {
	PtrDB          *gorm.DB
	PtrCurrentUser *User
	LoggedIn       bool
}

//Initialize : db interface init function
func Initialize() (*gorm.DB, error) {
	db, err := connectToDB()
	if err != nil {
		return nil, err
	}

	if db.Migrator().HasTable(&ActiveStream{}) {
		//clear all active stream
		db.Migrator().DropTable(&ActiveStream{})
	}
	err = db.AutoMigrate(&User{}, &ActiveStream{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() (*gorm.DB, error) {
	//setup basic logger
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,   // Slow SQL threshold
			LogLevel:      logger.Silent, // Log level
			Colorful:      true,          // Enable color
		},
	)

	configure.Init()
	dbConf := configure.DbConfig{}

	configure.Config.UnmarshalKey("db", &dbConf)

	fmt.Println(dbConf)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/gocoursetest?charset=utf8mb4&parseTime=True&loc=Local", dbConf.User, dbConf.Password, dbConf.Host, dbConf.Port)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	sqlDb, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDb.SetMaxIdleConns(10)
	sqlDb.SetMaxOpenConns(100)

	return db, err
}
