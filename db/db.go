package db

import (
	"context"
	"database/sql"
	"errors"
	perrors "github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/syncx"
	"gorm.io/gorm"
	"time"
)

const cacheSafeGapBetweenIndexAndPrimary = time.Second * 5

var (
	ErrCacheNil = errors.New("cache nil")

	// can't use one SingleFlight per conn, because multiple conns may share the same cache key.
	exclusiveCalls = syncx.NewSingleFlight()
	stats          = cache.NewStat("yorm")
)

type (
	// ExecFn defines the exec method.
	ExecFn func(conn *gorm.DB) (interface{}, error)
	// ExecCtxFn defines the exec method with ctx.
	ExecCtxFn func(ctx context.Context, conn *gorm.DB) (interface{}, error)
	// IndexQueryFn defines the query method that based on unique indexes.
	IndexQueryFn func(conn *gorm.DB, v interface{}) (interface{}, error)
	// IndexQueryCtxFn defines the query method that based on unique indexes with ctx.
	IndexQueryCtxFn func(ctx context.Context, conn *gorm.DB, v interface{}) (interface{}, error)
	// PrimaryQueryFn defines the query method that based on primary keys.
	PrimaryQueryFn func(conn *gorm.DB, v, primary interface{}) error
	// PrimaryQueryCtxFn defines the query method that based on primary keys with ctx.
	PrimaryQueryCtxFn func(ctx context.Context, conn *gorm.DB, v, primary interface{}) error
	// QueryFn defines the query method.
	QueryFn func(conn *gorm.DB, v interface{}) error
	// QueryCtxFn defines the query method with ctx.
	QueryCtxFn func(ctx context.Context, conn *gorm.DB, v interface{}) error
	// TXFn defines the transaction method.
	TXFn func(conn *DBConn) error
	// TXCtxFn defines the transaction method with ctx.
	TXCtxFn func(ctx context.Context, conn *DBConn) error

	DBConn struct {
		db    *gorm.DB
		cache cache.Cache
	}
)

func NewDBConn(db *gorm.DB, c cache.CacheConf, opts ...cache.Option) *DBConn {
	if len(c) == 0 || cache.TotalWeights(c) <= 0 {
		return &DBConn{
			db: db,
		}
	}
	return &DBConn{
		db:    db,
		cache: cache.New(c, exclusiveCalls, stats, sql.ErrNoRows, opts...),
	}
}

// DB *gorm.DB
func (cc *DBConn) DB() *gorm.DB {
	return cc.db
}

// Cache entity
func (cc *DBConn) Cache() cache.Cache {
	return cc.cache
}

// DelCache deletes cache with keys.
func (cc *DBConn) DelCache(keys ...string) error {
	return cc.DelCacheCtx(context.Background(), keys...)
}

// DelCacheCtx deletes cache with keys.
func (cc *DBConn) DelCacheCtx(ctx context.Context, keys ...string) error {
	if cc.cache != nil {
		if err := cc.cache.DelCtx(ctx, keys...); err != nil {
			return perrors.WithStack(err)
		}
	}
	return nil
}

// GetCache unmarshals cache with given key into v.
func (cc *DBConn) GetCache(key string, v interface{}) error {
	return cc.GetCacheCtx(context.Background(), key, v)
}

// GetCacheCtx unmarshals cache with given key into v.
func (cc *DBConn) GetCacheCtx(ctx context.Context, key string, v interface{}) error {
	if cc.cache != nil {
		if err := cc.cache.GetCtx(ctx, key, v); err != nil {
			return perrors.WithStack(err)
		}
		return nil
	}
	return perrors.WithStack(ErrCacheNil)
}

// Exec runs given exec on given keys, and returns execution result.
func (cc *DBConn) Exec(exec ExecFn, keys ...string) (interface{}, error) {
	execCtx := func(_ context.Context, conn *gorm.DB) (interface{}, error) {
		return exec(cc.db)
	}
	return cc.ExecCtx(context.Background(), execCtx, keys...)
}

// ExecCtx runs given exec on given keys, and returns execution result.
func (cc *DBConn) ExecCtx(ctx context.Context, exec ExecCtxFn, keys ...string) (interface{}, error) {
	res, err := exec(ctx, cc.db)
	if err != nil {
		return nil, err
	}
	if err := cc.DelCacheCtx(ctx, keys...); err != nil {
		return nil, perrors.WithStack(err)
	}
	return res, nil
}

// ExecNoCache runs exec with given sql statement, without affecting cache.
func (cc *DBConn) ExecNoCache(q string, args ...interface{}) (*sql.Rows, error) {
	return cc.ExecNoCacheCtx(context.Background(), q, args...)
}

// ExecNoCacheCtx runs exec with given sql statement, without affecting cache.
func (cc *DBConn) ExecNoCacheCtx(ctx context.Context, q string, args ...interface{}) (*sql.Rows, error) {
	rows, err := cc.db.Exec(q, args...).Rows()
	if err != nil {
		return nil, perrors.WithStack(err)
	}
	return rows, nil
}

// QueryRow unmarshals into v with given key and query func.
func (cc *DBConn) QueryRow(v interface{}, key string, query QueryFn) error {
	queryCtx := func(_ context.Context, conn *gorm.DB, v interface{}) error {
		return query(conn, v)
	}
	return cc.QueryRowCtx(context.Background(), v, key, queryCtx)
}

// QueryRowCtx unmarshals into v with given key and query func.
func (cc *DBConn) QueryRowCtx(ctx context.Context, v interface{}, key string, query QueryCtxFn) error {
	if cc.cache == nil {
		if err := cc.QueryRowNoCacheCtx(ctx, v, query); err != nil {
			return err
		}
	}
	err := cc.cache.TakeCtx(ctx, v, key, func(v interface{}) error {
		return query(ctx, cc.db, v)
	})
	if err != nil {
		return perrors.WithStack(err)
	}
	return nil
}

// QueryRowIndex unmarshals into v with given key.
func (cc *DBConn) QueryRowIndex(v interface{}, key string,
	keyer func(primary interface{}) string, indexQuery IndexQueryFn, primaryQuery PrimaryQueryFn) error {
	indexQueryCtx := func(_ context.Context, conn *gorm.DB, v interface{}) (interface{}, error) {
		return indexQuery(conn, v)
	}
	primaryQueryCtx := func(_ context.Context, conn *gorm.DB, v, primary interface{}) error {
		return primaryQuery(conn, v, primary)
	}
	return cc.QueryRowIndexCtx(context.Background(), v, key, keyer, indexQueryCtx, primaryQueryCtx)
}

// QueryRowIndexCtx unmarshals into v with given key.
func (cc *DBConn) QueryRowIndexCtx(ctx context.Context, v interface{}, key string,
	keyer func(primary interface{}) string, indexQuery IndexQueryCtxFn, primaryQuery PrimaryQueryCtxFn) error {
	var primaryKey interface{}
	var found bool

	if cc.cache == nil {
		return cc.QueryRowNoCacheCtx(ctx, v, func(ctx context.Context, conn *gorm.DB, v interface{}) error {
			if _, err := indexQuery(ctx, conn, v); err != nil {
				return err
			}
			return nil
		})
	}

	if err := cc.cache.TakeWithExpireCtx(ctx, &primaryKey, key,
		func(val interface{}, expire time.Duration) (err error) {
			primaryKey, err = indexQuery(ctx, cc.db, v)
			if err != nil {
				return
			}

			found = true
			return cc.cache.SetWithExpireCtx(ctx, keyer(primaryKey), v,
				expire+cacheSafeGapBetweenIndexAndPrimary)
		}); err != nil {
		return perrors.WithStack(err)
	}
	if found {
		return nil
	}

	err := cc.cache.TakeCtx(ctx, v, keyer(primaryKey), func(v interface{}) error {
		return primaryQuery(ctx, cc.db, v, primaryKey)
	})
	if err != nil {
		return perrors.WithStack(err)
	}
	return nil
}

// QueryRowNoCache unmarshals into v with given statement.
func (cc *DBConn) QueryRowNoCache(v interface{}, query QueryFn) error {
	return query(cc.db, v)
}

// QueryRowNoCacheCtx unmarshals into v with given statement.
func (cc *DBConn) QueryRowNoCacheCtx(ctx context.Context, v interface{}, queryCtx QueryCtxFn) error {
	return queryCtx(ctx, cc.db, v)
}

// SetCache sets v into cache with given key.
func (cc *DBConn) SetCache(key string, val interface{}) error {
	return cc.SetCacheCtx(context.Background(), key, val)
}

// SetCacheCtx sets v into cache with given key.
func (cc *DBConn) SetCacheCtx(ctx context.Context, key string, v interface{}) error {
	if cc.cache != nil {
		if err := cc.cache.SetCtx(ctx, key, v); err != nil {
			return perrors.WithStack(err)
		}
	}
	return nil
}

// Transact runs given fn in transaction mode.
func (cc *DBConn) Transact(fn TXFn) error {
	tx := cc.db.Begin()
	if err := fn(cc); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		return perrors.WithStack(err)
	}
	return nil
}

// TransactCtx runs given fn in transaction mode.
func (cc *DBConn) TransactCtx(ctx context.Context, fn TXCtxFn) error {
	tx := cc.db.Begin()
	if err := fn(ctx, cc); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		return perrors.WithStack(err)
	}
	return nil
}

//func (cc *DBConn) InsertIndex(ctx context.Context, insert InsertFn, keys ...string) error {
//	if err := insert(cc.db); err != nil {
//		return err
//	}
//	if err := cc.DelCache(ctx, keys...); err != nil {
//		return err
//	}
//	return nil
//}
//
//func (cc DBConn) DelIndex(ctx context.Context, del DeleteFn, keys ...string) error {
//	if err := del(cc.db); err != nil {
//		return err
//	}
//	if err := cc.DelCache(ctx, keys...); err != nil {
//		return err
//	}
//	return nil
//}
//
//func (cc DBConn) UpdateIndex(ctx context.Context, update UpdateFn, keys ...string) error {
//	if err := update(cc.db); err != nil {
//		return err
//	}
//	if err := cc.DelCache(ctx, keys...); err != nil {
//		return err
//	}
//	return nil
//}
//
//func (cc DBConn) FindIndex(ctx context.Context, v interface{}, key string, find FindFn) error {
//	if cc.cache == nil {
//		return find(cc.db, v)
//	}
//	return cc.cache.TakeWithExpireCtx(ctx, v, key, func(v interface{}, expire time.Duration) error {
//		return find(cc.db, v)
//	})
//}
//
//func (cc DBConn) FindNoCache(v interface{}, find FindFn) error {
//	return find(cc.db, v)
//}
