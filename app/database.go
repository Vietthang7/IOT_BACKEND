package app

import (
	"backend/internal/model"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const CTimeOut = 10 * time.Second

type DatabaseConfig struct {
	*gorm.DB
	Driver      string `yaml:"driver" env:"DB_DRIVER"`
	Host        string `yaml:"host" env:"DB_HOST"`
	Username    string `yaml:"username" env:"DB_USER"`
	Password    string `yaml:"password" env:"DB_PASSWORD"`
	DBName      string `yaml:"db_name" env:"DB_NAME"`
	Port        string `yaml:"port" env:"DB_PORT"`
	Connections int    `yaml:"connections" env:"DB_CONNECTIONS"`
	Debug       bool   `yaml:"debug"`
	MaxIdleConn int    `env:"MAX_IDLE_CONN"`
	MaxOpenConn int    `env:"MAX_OPEN_CONN"`
	MaxLifetime int64  `env:"MAX_LIFE_TIME_PER_CONN"`
	sqlDB       *sql.DB
}

func (cg *DatabaseConfig) Setup() {
	logrus.SetLevel(logrus.DebugLevel)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,         // Disable color
		},
	)
	mainDbDNS := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", cg.Username, cg.Password, cg.Host, cg.Port, cg.DBName)
	DB, err := gorm.Open(
		mysql.New(mysql.Config{
			DSN:               mainDbDNS,
			DefaultStringSize: 256, // default size for string fields
			// DisableDatetimePrecision:  true,  // disable datetime precision, which not supported before MySQL 5.6
			// DontSupportRenameIndex:    true,  // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
			// DontSupportRenameColumn:   true,  // `change` when rename column, rename column not supported before MySQL 8, MariaDB
			// SkipInitializeWithVersion: false, // auto configure based on currently MySQL version
		}),
		&gorm.Config{
			PrepareStmt: true,
			Logger:      newLogger,
			//DisableForeignKeyConstraintWhenMigrating: true,
		})
	if err != nil {
		logrus.Panic("Failed to connect database: "+mainDbDNS, err)
	}
	sqlDB, _ := DB.DB()
	if cg.MaxIdleConn > 0 {
		sqlDB.SetMaxIdleConns(cg.MaxIdleConn) // MAX_IDLE_CONN
	}
	if cg.MaxOpenConn > 0 {
		sqlDB.SetMaxOpenConns(cg.MaxOpenConn) // MAX_OPEN_CONN
	}
	if cg.MaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(time.Duration(cg.MaxLifetime) * time.Minute)
	}
	if cg.Debug {
		DB = DB
	}
	cg.DB = DB
	cg.sqlDB = sqlDB
	if err = MigrateDatabase(DB); err != nil {
		logrus.Fatal(err)
	}
	fmt.Println("*************** DB AUTO MIGRATE FINISHED  ***************")
}

func (cg *DatabaseConfig) CloseConnection() error {
	return cg.sqlDB.Close()
}

func MigrateDatabase(DB *gorm.DB) error {
	if err := DB.AutoMigrate(&model.DataSensor{}); err != nil {
		logrus.Debugf("Err AutoMigrate DataSensor: %v", err)
	}
	if err := DB.AutoMigrate(&model.DeviceHistory{}); err != nil {
		logrus.Debugf("Err AutoMigrate DeviceHistory: %v", err)
	}

	logrus.Info("Migration finish")
	return nil
}
