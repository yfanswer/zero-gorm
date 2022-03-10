package db

import (
	"database/sql"
	perrors "github.com/pkg/errors"
	"time"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/syncx"
	"gorm.io/gorm"
)

const cacheSafeGapBetweenIndexAndPrimary = time.Second * 5

var (
	// ErrNotFound is an alias of sqlx.ErrNotFound.
	ErrNotFound = sqlx.ErrNotFound

	// can't use one SingleFlight per conn, because multiple conns may share the same cache key.
	exclusiveCalls = syncx.NewSingleFlight()
	stats          = cache.NewStat("yorm")
)

type (
	// ExecFn defines the exec method.
	ExecFn func(conn *gorm.DB) (interface{}, error)
	// IndexQueryFn defines the query method that based on unique indexes.
	IndexQueryFn func(conn *gorm.DB, v interface{}) (interface{}, error)
	// PrimaryQueryFn defines the query method that based on primary keys.
	PrimaryQueryFn func(conn *gorm.DB, v, primary interface{}) error
	// QueryFn defines the query method.
	QueryFn func(conn *gorm.DB, v interface{}) error
	// FindFn defines the find method.
	FindFn func(conn *gorm.DB, v interface{}) error
	// InsertFn defines the insert method.
	InsertFn func(conn *gorm.DB) error
	// DeleteFn defines the delete method.
	DeleteFn func(conn *gorm.DB) error
	// UpdateFn defines the update method.
	UpdateFn func(conn *gorm.DB) error
	// RawsFn defines the raw method.
	RawsFn func(*sql.Rows) error
	// TXFn defines the transaction method.
	TXFn func(conn DBConn) error

	DBConn struct {
		db    *gorm.DB
		cache cache.Cache
	}
)

func NewDBConn(db *gorm.DB, c cache.CacheConf, opts ...cache.Option) DBConn {
	if len(c) == 0 || cache.TotalWeights(c) <= 0 {
		return DBConn{
			db: db,
		}
	}
	return DBConn{
		db:    db,
		cache: cache.New(c, exclusiveCalls, stats, sql.ErrNoRows, opts...),
	}
}

// DB *gorm.DB
func (cc DBConn) DB() *gorm.DB {
	return cc.db
}

// Cache entity
func (cc DBConn) Cache() cache.Cache {
	return cc.cache
}

// DelCache deletes cache with keys.
func (cc DBConn) DelCache(keys ...string) error {
	if cc.cache != nil {
		return cc.cache.Del(keys...)
	}
	return nil
}

func (cc DBConn) SetCache(key string, v interface{}) error {
	if cc.cache != nil {
		return cc.cache.Set(key, v)
	}
	return nil
}

// GetCache unmarshal cache with given key into v.
func (cc DBConn) GetCache(key string, v interface{}) error {
	if cc.cache != nil {
		return cc.cache.Get(key, v)
	}
	return perrors.New("cache nil")
}

// ExecFn Exec runs given exec on given keys, and returns execution result.
func (cc DBConn) ExecFn(exec ExecFn) (interface{}, error) {
	return exec(cc.db)
}

// ExecNoCache runs exec with given sql statement, without affecting cache.
func (cc DBConn) ExecNoCache(q string, args ...interface{}) error {
	return cc.db.Exec(q, args...).Error
}

// QueryRow unmarshals into v with given key and query func.
func (cc DBConn) QueryRow(v interface{}, key string, query QueryFn) error {
	if cc.cache == nil {
		return query(cc.db, v)
	}
	return cc.cache.Take(v, key, func(v interface{}) error {
		return query(cc.db, v)
	})
}

//// QueryRowIndex unmarshals into v with given key.
//func (cc DBConn) QueryRowIndex(v interface{}, key string, keyer func(primary interface{}) string,
//	indexQuery IndexQueryFn, primaryQuery PrimaryQueryFn) error {
//	var primaryKey interface{}
//	var found bool
//
//	if err := cc.cache.TakeWithExpire(&primaryKey, key, func(val interface{}, expire time.Duration) (err error) {
//		primaryKey, err = indexQuery(cc.db, v)
//		if err != nil {
//			return
//		}
//
//		found = true
//		return cc.cache.SetWithExpire(keyer(primaryKey), v, expire+cacheSafeGapBetweenIndexAndPrimary)
//	}); err != nil {
//		return err
//	}
//
//	if found {
//		return nil
//	}
//
//	return cc.cache.Take(v, keyer(primaryKey), func(v interface{}) error {
//		return primaryQuery(cc.db, v, primaryKey)
//	})
//}

// // QueryRowNoCache unmarshals into v with given statement.
// func (cc DBConn) QueryRowNoCache(v interface{}, q string, args ...interface{}) error {
// 	cc.db.Raw(q, args...).Scan(&v).ro
// 	cc.db.
// 	cc.db.Exec(q, args...).Scan()
// 	return cc.db.QueryRow(v, q, args...)
// }

func (cc DBConn) QueryRowsNoCache(rawsFunc RawsFn, q string, values ...interface{}) error {
	rows, err := cc.db.Raw(q, values...).Rows()
	if err != nil {
		return err
	}
	defer rows.Close()

	if err = rawsFunc(rows); err != nil {
		return err
	}

	return nil
}

// SetCache sets v into cache with given key.

// Transaction runs given fn in transaction mode.
func (cc DBConn) Transaction(fn TXFn) error {
	tx := cc.db.Begin()
	if err := fn(cc); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

func (cc DBConn) InsertIndex(insert InsertFn, keys ...string) error {
	if err := insert(cc.db); err != nil {
		return err
	}

	if err := cc.DelCache(keys...); err != nil {
		return err
	}

	return nil
}

func (cc DBConn) DelIndex(del DeleteFn, keys ...string) error {
	if err := del(cc.db); err != nil {
		return err
	}

	if err := cc.DelCache(keys...); err != nil {
		return err
	}

	return nil
}

func (cc DBConn) UpdateIndex(update UpdateFn, keys ...string) error {
	if err := update(cc.db); err != nil {
		return err
	}

	if err := cc.DelCache(keys...); err != nil {
		return err
	}

	return nil
}

func (cc DBConn) FindIndex(v interface{}, key string, find FindFn) error {
	if cc.cache == nil {
		return find(cc.db, v)
	}
	return cc.cache.Take(v, key, func(v interface{}) error {
		return find(cc.db, v)
	})
}

func (cc DBConn) FindNoCache(v interface{}, find FindFn) error {
	return find(cc.db, v)
}
