package database

import (
	"database/sql"
	"sync"

	"narou/infrastructure/database/internal"

	"gorm.io/gorm"
)

// DBM is an interface which DB implements
type DBM interface {
	Close() error
	DB() (*sql.DB, error)
	New() (DBM, error)
	Where(query any, args ...any) DBM
	Or(query any, args ...any) DBM
	Not(query any, args ...any) DBM
	Limit(value int) DBM
	Offset(value int) DBM
	Order(value any) DBM
	Select(query any, args ...any) DBM
	Omit(columns ...string) DBM
	Group(query string) DBM
	Having(query string, values ...any) DBM
	Joins(query string, args ...any) DBM
	Scopes(funcs ...func(DBM) DBM) DBM
	Unscoped() DBM
	Attrs(attrs ...any) DBM
	Assign(attrs ...any) DBM
	First(out any, where ...any) DBM
	Last(out any, where ...any) DBM
	Find(out any, where ...any) DBM
	Scan(dest any) DBM
	Row() *sql.Row
	Rows() (*sql.Rows, error)
	ScanRows(rows *sql.Rows, result any) error
	Pluck(column string, value any) DBM
	Count(value *int64) DBM
	FirstOrInit(out any, where ...any) DBM
	FirstOrCreate(out any, where ...any) DBM
	Update(colName string, value any) DBM
	Updates(values any) DBM
	UpdateColumn(colName string, value any) DBM
	UpdateColumns(values any) DBM
	Save(value any) DBM
	Create(value any) DBM
	Delete(value any, where ...any) DBM
	Raw(sql string, values ...any) DBM
	Exec(sql string, values ...any) DBM
	Model(value any) DBM
	Table(name string) DBM
	Debug() DBM
	Begin() DBM
	Commit() DBM
	Rollback() DBM
	CreateTable(values ...any) error
	DropTable(values ...any) error
	HasTable(value any) bool
	AutoMigrate(values ...any) error
	Association(column string) *gorm.Association
	Preload(column string, conditions ...any) DBM
	Set(name string, value any) DBM
	Get(name string) (value any, ok bool)
	AddError(err error) error
	Take(out any, where ...any) DBM
	Error() error
	RowsAffected() int64
}

type gormWrapper struct {
	w *gorm.DB
}

var single = &gormWrapper{}

func GetConn() (DBM, error) {
	if single.w == nil {
		return single, &rdbError{
			message:       "[database.GetConn] invalid connection",
			originalError: ErrNotOpened,
		}
	}

	return single, nil
}

var once sync.Once

// OpenDB is a drop-in replacement for Open().
func OpenDB(path string) error {
	var ret error
	once.Do(
		func() {
			gormDB, err := internal.OpenDB(path)
			single.w = gormDB
			ret = err
		})

	return ret
}

// Wrap wraps gorm.DB in an interface.
func Wrap(db *gorm.DB) DBM {
	return &gormWrapper{db}
}

func (it *gormWrapper) Close() error {
	panic("not implement.")
}

func (it *gormWrapper) DB() (*sql.DB, error) {
	db, err := it.w.DB()
	if err != nil {
		return nil, &rdbError{
			message:       "[database.DB] failed to get db",
			originalError: err,
		}
	}

	return db, nil
}

func (it *gormWrapper) New() (DBM, error) {
	if it.w == nil {
		return nil, &rdbError{
			message:       "[database.New] failed to get db",
			originalError: ErrNotOpened,
		}
	}

	return Wrap(it.w), nil
}

func (it *gormWrapper) Where(query any, args ...any) DBM {
	return Wrap(it.w.Where(query, args...))
}

func (it *gormWrapper) Or(query any, args ...any) DBM {
	return Wrap(it.w.Or(query, args...))
}

func (it *gormWrapper) Not(query any, args ...any) DBM {
	return Wrap(it.w.Not(query, args...))
}

func (it *gormWrapper) Limit(value int) DBM {
	return Wrap(it.w.Limit(value))
}

func (it *gormWrapper) Offset(value int) DBM {
	return Wrap(it.w.Offset(value))
}

func (it *gormWrapper) Order(value any) DBM {
	return Wrap(it.w.Order(value))
}

func (it *gormWrapper) Select(query any, args ...any) DBM {
	return Wrap(it.w.Select(query, args...))
}

func (it *gormWrapper) Omit(columns ...string) DBM {
	return Wrap(it.w.Omit(columns...))
}

func (it *gormWrapper) Group(query string) DBM {
	return Wrap(it.w.Group(query))
}

func (it *gormWrapper) Having(query string, values ...any) DBM {
	return Wrap(it.w.Having(query, values...))
}

func (it *gormWrapper) Joins(query string, args ...any) DBM {
	return Wrap(it.w.Joins(query, args...))
}

func (it *gormWrapper) Scopes(fs ...func(DBM) DBM) DBM {
	for _, f := range fs {
		it = f(it).(*gormWrapper)
	}

	return it
}

func (it *gormWrapper) Unscoped() DBM {
	return Wrap(it.w.Unscoped())
}

func (it *gormWrapper) Attrs(attrs ...any) DBM {
	return Wrap(it.w.Attrs(attrs...))
}

func (it *gormWrapper) Assign(attrs ...any) DBM {
	return Wrap(it.w.Assign(attrs...))
}

func (it *gormWrapper) First(out any, where ...any) DBM {
	return Wrap(it.w.First(out, where...))
}

func (it *gormWrapper) Last(out any, where ...any) DBM {
	return Wrap(it.w.Last(out, where...))
}

func (it *gormWrapper) Find(out any, where ...any) DBM {
	return Wrap(it.w.Find(out, where...))
}

func (it *gormWrapper) Scan(dest any) DBM {
	return Wrap(it.w.Scan(dest))
}

func (it *gormWrapper) Row() *sql.Row {
	return it.w.Row()
}

func (it *gormWrapper) Rows() (*sql.Rows, error) {
	return it.w.Rows()
}

func (it *gormWrapper) ScanRows(rows *sql.Rows, result any) error {
	return it.w.ScanRows(rows, result)
}

func (it *gormWrapper) Pluck(column string, value any) DBM {
	return Wrap(it.w.Pluck(column, value))
}

func (it *gormWrapper) Count(value *int64) DBM {
	return Wrap(it.w.Count(value))
}

func (it *gormWrapper) FirstOrInit(out any, where ...any) DBM {
	return Wrap(it.w.FirstOrInit(out, where...))
}

func (it *gormWrapper) FirstOrCreate(out any, where ...any) DBM {
	return Wrap(it.w.FirstOrCreate(out, where...))
}

func (it *gormWrapper) Update(colName string, value any) DBM {
	return Wrap(it.w.Update(colName, value))
}

func (it *gormWrapper) Updates(values any) DBM {
	return Wrap(it.w.Updates(values))
}

func (it *gormWrapper) UpdateColumn(colName string, value any) DBM {
	return Wrap(it.w.UpdateColumn(colName, value))
}

func (it *gormWrapper) UpdateColumns(values any) DBM {
	return Wrap(it.w.UpdateColumns(values))
}

func (it *gormWrapper) Save(value any) DBM {
	return Wrap(it.w.Save(value))
}

func (it *gormWrapper) Create(value any) DBM {
	return Wrap(it.w.Create(value))
}

func (it *gormWrapper) Delete(value any, where ...any) DBM {
	return Wrap(it.w.Delete(value, where...))
}

func (it *gormWrapper) Raw(sql string, values ...any) DBM {
	return Wrap(it.w.Raw(sql, values...))
}

func (it *gormWrapper) Exec(sql string, values ...any) DBM {
	return Wrap(it.w.Exec(sql, values...))
}

func (it *gormWrapper) Model(value any) DBM {
	return Wrap(it.w.Model(value))
}

func (it *gormWrapper) Table(name string) DBM {
	return Wrap(it.w.Table(name))
}

func (it *gormWrapper) Debug() DBM {
	return Wrap(it.w.Debug())
}

func (it *gormWrapper) Begin() DBM {
	return Wrap(it.w.Begin())
}

func (it *gormWrapper) Commit() DBM {
	return Wrap(it.w.Commit())
}

func (it *gormWrapper) Rollback() DBM {
	return Wrap(it.w.Rollback())
}

func (it *gormWrapper) CreateTable(values ...any) error {
	return it.w.Migrator().CreateTable(values...)
}

func (it *gormWrapper) DropTable(values ...any) error {
	return it.w.Migrator().DropTable(values...)
}

func (it *gormWrapper) HasTable(value any) bool {
	return it.w.Migrator().HasTable(value)
}

func (it *gormWrapper) AutoMigrate(values ...any) error {
	return it.w.AutoMigrate(values...)
}

func (it *gormWrapper) Association(column string) *gorm.Association {
	return it.w.Association(column)
}

func (it *gormWrapper) Preload(column string, conditions ...any) DBM {
	return Wrap(it.w.Preload(column, conditions...))
}

func (it *gormWrapper) Set(name string, value any) DBM {
	return Wrap(it.w.Set(name, value))
}

func (it *gormWrapper) Get(name string) (any, bool) {
	return it.w.Get(name)
}

func (it *gormWrapper) AddError(err error) error {
	return it.w.AddError(err)
}

func (it *gormWrapper) RowsAffected() int64 {
	return it.w.RowsAffected
}

func (it *gormWrapper) Error() error {
	return it.w.Error
}

func (it *gormWrapper) Take(out any, where ...any) DBM {
	return Wrap(it.w.Take(out, where...))
}
