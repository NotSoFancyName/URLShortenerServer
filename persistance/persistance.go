package persistance

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

const (
	dbname = "urls.db"
	dbms   = "sqlite3"
)

type ShortLongURL struct {
	gorm.Model
	Short string
	Long  string
}

func init() {
	db := getConnection()
	defer db.Close()
	db.AutoMigrate(&ShortLongURL{})
}

func getConnection() *gorm.DB {
	db, err := gorm.Open(dbms, dbname)
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

func GetCounter() int64 {
	db := getConnection()
	defer db.Close()

	m := ShortLongURL{}
	db.Order("created_at desc").Last(&m)

	return int64(m.ID)
}

func Save(s string, l string) {
	db := getConnection()
	defer db.Close()

	u := ShortLongURL{}
	if db.First(&u, "short = ?", s); u.Short != "" {
		u.Long = l
		db.Save(&u)
	}
	db.Create(&ShortLongURL{Short: s, Long: l})
}

func GetURLFromDB(l string, parIsLongUrl bool) string {
	db := getConnection()
	defer db.Close()

	var u ShortLongURL
	if parIsLongUrl {
		db.First(&u, "long = ?", l)
		return u.Short
	} else {
		db.First(&u, "short = ?", l)
		return u.Long
	}
}
