package base

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	// mysql
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
)

type gormLogger struct{}

func (*gormLogger) Print(values ...interface{}) {
	if values[0] == "sql" {
		query := fmt.Sprintf("%.2fms", float64(values[2].(time.Duration).Nanoseconds()/1e4)/100.0)
		log.WithFields(log.Fields{"sql": values[3], "sql_bind": values[4], "used": query}).Info("[SQL]")
	}
}

// singleton
var _db *dBInstance
var once sync.Once

// DB database manager class
type dBInstance struct {
	Master *gorm.DB
	Slave  *gorm.DB
}

// DB returns a new Model without opening database connection
func DB() *dBInstance {
	var err error
	once.Do(func() {
		if _db == nil {
			_db = &dBInstance{}
			err = _db.init(DBConn)
		}
	})

	if err == nil {
		return _db
	}
	return nil
}

func (me *dBInstance) Close() {
	if me.Master != nil {
		me.Master.Close()
	}
	if me.Slave != nil {
		me.Slave.Close()
	}
}

func (me *dBInstance) init(str string) error {
	arr := strings.Split(str, ",")
	for _, v := range arr {
		if len(v) <= 0 {
			continue
		}
		if err := me.open(strings.TrimSpace(v)); err != nil {
			return err
		}
	}
	return nil
}

// Open opens database connection with the settings found in cfg
func (me *dBInstance) open(dbType string) error {
	switch dbType {
	case "master":
		me.Master = me.getConn(os.Getenv("DB_URL_MASTER"))
	case "slave":
		me.Slave = me.getConn(os.Getenv("DB_URL_SLAVE"))
	case "default":
		me.Master = me.getConn(os.Getenv("DB_URL"))
	default:
		return errors.New("db type err")
	}
	return nil
}

func (me *dBInstance) getConn(url string) *gorm.DB {
	conn, err := gorm.Open("mysql", url)
	log.Info(url)
	if err != nil {
		panic("fail open mysql connection")
	}

	_dbMaxIdleConns, _ := strconv.Atoi(os.Getenv("DB_MaxIdleConns"))
	_dbMaxOpenConns, _ := strconv.Atoi(os.Getenv("DB_MaxOpenConns"))

	conn.DB().SetMaxIdleConns(_dbMaxIdleConns)
	conn.DB().SetMaxOpenConns(_dbMaxOpenConns)

	if err = conn.DB().Ping(); err != nil {
		panic("ping mysql err")
	}

	if os.Getenv("DEBUG") == "true" {
		conn.LogMode(true)
		conn.SetLogger(&gormLogger{})
	}

	return conn
}
