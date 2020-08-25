// Package persistance implements simple funcs for db access.
// It also implements a sctruct that utilizes short and long URLs
package persistance

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

const (
	dbname = "urls.db"
	dbms   = "sqlite3"
)

type shortLongURL struct {
	gorm.Model
	Short string
	Long  string
}

func init() {
	db := getConnection()
	defer db.Close()
	db.AutoMigrate(&shortLongURL{})
}

func getConnection() *gorm.DB {
	db, err := gorm.Open(dbms, dbname)
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

// Returns the ID of the most recent inserted URL
// This ID equals counter + 1 from shortener package 
func GetMostRecentUpdatedEntryID() int64 {
	db := getConnection()
	defer db.Close()

	m := shortLongURL{}
	db.Order("created_at desc").Last(&m)

	return int64(m.ID)
}

// Creates short and long URL in DB if the short URL is not present
// Updates long URL if short URL is present  
func Save(s string, l string) {
	db := getConnection()
	defer db.Close()

	u := shortLongURL{}
	if db.First(&u, "short = ?", s); u.Short != "" {
		u.Long = l
		db.Save(&u)
	}
	db.Create(&shortLongURL{Short: s, Long: l})
}

// Returns short URL if parIsLongUrl is true
// Returns long URL if parIsLongUrl is false
// Returns empty str if there is no such URL
func GetURLFromDB(l string, parIsLongUrl bool) string {
	db := getConnection()
	defer db.Close()

	var u shortLongURL
	if parIsLongUrl {
		db.First(&u, "long = ?", l)
		return u.Short
	} else {
		db.First(&u, "short = ?", l)
		return u.Long
	}
}
