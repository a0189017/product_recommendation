package database

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/jpillora/backoff"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBInterface interface {
	GetDB() *gorm.DB
	LoadSchemaFields()
}

type DBOptions struct {
	User          string
	Password      string
	Name          string
	Host          string
	Port          string
	SlowThreshold string
	Colorful      string
}

type Database struct {
	DB     *gorm.DB
	fields []string
}

func New(dbOptions *DBOptions) DBInterface {
	if len(dbOptions.User) == 0 {
		panic("db user is empty")
	}
	if len(dbOptions.Password) == 0 {
		panic("db password is empty")
	}

	// connect to mysql
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&time_zone=UTC",
		dbOptions.User, dbOptions.Password, dbOptions.Host, dbOptions.Port, dbOptions.Name,
	)
	dialector := mysql.Open(dsn)

	b := &backoff.Backoff{
		Factor: 1.5,
		Min:    1 * time.Second,
		Max:    32 * time.Second,
	}

	slowThreshold := dbOptions.SlowThreshold
	if slowThreshold == "" {
		slowThreshold = "1000"
	}
	t, err := strconv.ParseInt(slowThreshold, 10, 0)
	if err != nil {
		panic("db slowThreshold should be a valid number")
	}
	if t <= 0 {
		panic("db slowThreshold should be greater than 0")
	}
	colorful := dbOptions.Colorful
	if colorful == "" {
		colorful = "false"
	}
	bColorful, err := strconv.ParseBool(colorful)
	if err != nil {
		panic("db color should be true or false")
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Duration(t) * time.Millisecond,
			Colorful:      bColorful,
		},
	)

	for {
		db, err := gorm.Open(dialector, &gorm.Config{
			CreateBatchSize: 100,
			NowFunc: func() time.Time {
				return time.Now().UTC()
			},
			Logger: newLogger,
		})
		if err != nil {
			d := b.Duration()
			fmt.Printf("%s, reconnecting in %s", err, d)
			if d == b.Max {
				panic(err)
			}
			time.Sleep(d)

			continue
		}
		//connected
		b.Reset()

		d := db
		ddb, err := d.DB()
		if err != nil {
			panic(err)
			// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
		}
		ddb.SetMaxIdleConns(20)
		// SetMaxOpenConns sets the maximum number of open connections to the database.
		ddb.SetMaxOpenConns(150)
		// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
		ddb.SetConnMaxLifetime(time.Minute * 5)

		d = db.Debug()

		dbInstance := &Database{DB: d}
		dbInstance.LoadSchemaFields()
		return dbInstance
	}
}

func (d *Database) GetDB() *gorm.DB {
	return d.DB
}

func (d *Database) LoadSchemaFields() {
	var query string
	fields := []string{}
	query = `
				SELECT column_name
				FROM information_schema.columns
				WHERE table_schema = DATABASE()
				GROUP BY column_name
			`

	if err := d.DB.Raw(query).Scan(&fields).Error; err != nil {
		panic(err)
	}

	d.fields = fields
}
