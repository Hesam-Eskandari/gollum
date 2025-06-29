package database

import (
	"fmt"
	"sync"
	"time"

	"github.com/Hesam-Eskandari/gollum/library/lrucache"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type ConnectionPool struct {
	HostTag
	Conn *gorm.DB
}

var connectionsCache = lrucache.NewLRUCache[HostTag](5)

var mu sync.Mutex

func NewConnectionPool(config ConnectionConfig) (*ConnectionPool, error) {
	mu.Lock()
	defer mu.Unlock()
	var connPool *ConnectionPool
	poolGeneric, err := connectionsCache.Read(config.HostTag)
	if err == nil {
		connPool = poolGeneric.(*ConnectionPool)
		return connPool, nil
	}
	conn, err := gorm.Open(postgres.Open(buildPostgresConnectionString(config)), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		return nil, err
	}
	db, _ := conn.DB()
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(time.Hour * 24)
	connPool = &ConnectionPool{
		HostTag: config.HostTag,
		Conn:    conn,
	}
	_, _ = connectionsCache.Add(config.HostTag, connPool)
	return connPool, nil
}

func buildPostgresConnectionString(config ConnectionConfig) string {
	sslMode := "disable"
	if config.SSL {
		sslMode = "enable"
	}
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		config.User, config.Password, config.Host, config.Port, config.Database, sslMode)
}
