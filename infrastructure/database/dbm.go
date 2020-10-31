package database

import (
	"database/sql"

	"narou/infrastructure/database/internal"

	"gorm.io/gorm"
)

// DBM is an interface which DB implements
type DBM interface {
	Close() error
	DB() (*sql.DB, error)
	New() (DBM, error)
	Where(query interface{}, args ...interface{}) DBM
	Or(query interface{}, args ...interface{}) DBM
	Not(query interface{}, args ...interface{}) DBM
	Limit(value int) DBM
	Offset(value int) DBM
	Order(value interface{}) DBM
	Select(query interface{}, args ...interface{}) DBM
	Omit(columns ...string) DBM
	Group(query string) DBM
	Having(query string, values ...interface{}) DBM
	Joins(query string, args ...interface{}) DBM
	Scopes(funcs ...func(DBM) DBM) DBM
	Unscoped() DBM
	Attrs(attrs ...interface{}) DBM
	Assign(attrs ...interface{}) DBM
	First(out interface{}, where ...interface{}) DBM
	Last(out interface{}, where ...interface{}) DBM
	Find(out interface{}, where ...interface{}) DBM
	Scan(dest interface{}) DBM
	Row() *sql.Row
	Rows() (*sql.Rows, error)
	ScanRows(rows *sql.Rows, result interface{}) error
	Pluck(column string, value interface{}) DBM
	Count(value *int64) DBM
	FirstOrInit(out interface{}, where ...interface{}) DBM
	FirstOrCreate(out interface{}, where ...interface{}) DBM
	Update(colName string, value interface{}) DBM
	Updates(values interface{}) DBM
	UpdateColumn(colName string, value interface{}) DBM
	UpdateColumns(values interface{}) DBM
	Save(value interface{}) DBM
	Create(value interface{}) DBM
	Delete(value interface{}, where ...interface{}) DBM
	Raw(sql string, values ...interface{}) DBM
	Exec(sql string, values ...interface{}) DBM
	Model(value interface{}) DBM
	Table(name string) DBM
	Debug() DBM
	Begin() DBM
	Commit() DBM
	Rollback() DBM
	CreateTable(values ...interface{}) error
	DropTable(values ...interface{}) error
	HasTable(value interface{}) bool
	AutoMigrate(values ...interface{}) error
	Association(column string) *gorm.Association
	Preload(column string, conditions ...interface{}) DBM
	Set(name string, value interface{}) DBM
	Get(name string) (value interface{}, ok bool)
	AddError(err error) error

	// extra
	Take(out interface{}, where ...interface{}) DBM
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

// OpenDB is a drop-in replacement for Open().
func OpenDB(path string) (err error) {
	gormDB, err := internal.OpenDB(path)
	single.w = gormDB

	return err
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

func (it *gormWrapper) Where(query interface{}, args ...interface{}) DBM {
	return Wrap(it.w.Where(query, args...))
}

func (it *gormWrapper) Or(query interface{}, args ...interface{}) DBM {
	return Wrap(it.w.Or(query, args...))
}

func (it *gormWrapper) Not(query interface{}, args ...interface{}) DBM {
	return Wrap(it.w.Not(query, args...))
}

func (it *gormWrapper) Limit(value int) DBM {
	return Wrap(it.w.Limit(value))
}

func (it *gormWrapper) Offset(value int) DBM {
	return Wrap(it.w.Offset(value))
}

func (it *gormWrapper) Order(value interface{}) DBM {
	return Wrap(it.w.Order(value))
}

func (it *gormWrapper) Select(query interface{}, args ...interface{}) DBM {
	return Wrap(it.w.Select(query, args...))
}

func (it *gormWrapper) Omit(columns ...string) DBM {
	return Wrap(it.w.Omit(columns...))
}

func (it *gormWrapper) Group(query string) DBM {
	return Wrap(it.w.Group(query))
}

func (it *gormWrapper) Having(query string, values ...interface{}) DBM {
	return Wrap(it.w.Having(query, values...))
}

func (it *gormWrapper) Joins(query string, args ...interface{}) DBM {
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

func (it *gormWrapper) Attrs(attrs ...interface{}) DBM {
	return Wrap(it.w.Attrs(attrs...))
}

func (it *gormWrapper) Assign(attrs ...interface{}) DBM {
	return Wrap(it.w.Assign(attrs...))
}

func (it *gormWrapper) First(out interface{}, where ...interface{}) DBM {
	return Wrap(it.w.First(out, where...))
}

func (it *gormWrapper) Last(out interface{}, where ...interface{}) DBM {
	return Wrap(it.w.Last(out, where...))
}

func (it *gormWrapper) Find(out interface{}, where ...interface{}) DBM {
	return Wrap(it.w.Find(out, where...))
}

func (it *gormWrapper) Scan(dest interface{}) DBM {
	return Wrap(it.w.Scan(dest))
}

func (it *gormWrapper) Row() *sql.Row {
	return it.w.Row()
}

func (it *gormWrapper) Rows() (*sql.Rows, error) {
	return it.w.Rows()
}

func (it *gormWrapper) ScanRows(rows *sql.Rows, result interface{}) error {
	return it.w.ScanRows(rows, result)
}

func (it *gormWrapper) Pluck(column string, value interface{}) DBM {
	return Wrap(it.w.Pluck(column, value))
}

func (it *gormWrapper) Count(value *int64) DBM {
	return Wrap(it.w.Count(value))
}

func (it *gormWrapper) FirstOrInit(out interface{}, where ...interface{}) DBM {
	return Wrap(it.w.FirstOrInit(out, where...))
}

func (it *gormWrapper) FirstOrCreate(out interface{}, where ...interface{}) DBM {
	return Wrap(it.w.FirstOrCreate(out, where...))
}

func (it *gormWrapper) Update(colName string, value interface{}) DBM {
	return Wrap(it.w.Update(colName, value))
}

func (it *gormWrapper) Updates(values interface{}) DBM {
	return Wrap(it.w.Updates(values))
}

func (it *gormWrapper) UpdateColumn(colName string, value interface{}) DBM {
	return Wrap(it.w.UpdateColumn(colName, value))
}

func (it *gormWrapper) UpdateColumns(values interface{}) DBM {
	return Wrap(it.w.UpdateColumns(values))
}

func (it *gormWrapper) Save(value interface{}) DBM {
	return Wrap(it.w.Save(value))
}

func (it *gormWrapper) Create(value interface{}) DBM {
	return Wrap(it.w.Create(value))
}

func (it *gormWrapper) Delete(value interface{}, where ...interface{}) DBM {
	return Wrap(it.w.Delete(value, where...))
}

func (it *gormWrapper) Raw(sql string, values ...interface{}) DBM {
	return Wrap(it.w.Raw(sql, values...))
}

func (it *gormWrapper) Exec(sql string, values ...interface{}) DBM {
	return Wrap(it.w.Exec(sql, values...))
}

func (it *gormWrapper) Model(value interface{}) DBM {
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

func (it *gormWrapper) CreateTable(values ...interface{}) error {
	return it.w.Migrator().CreateTable(values...)
}

func (it *gormWrapper) DropTable(values ...interface{}) error {
	return it.w.Migrator().DropTable(values...)
}

func (it *gormWrapper) HasTable(value interface{}) bool {
	return it.w.Migrator().HasTable(value)
}

func (it *gormWrapper) AutoMigrate(values ...interface{}) error {
	return it.w.AutoMigrate(values...)
}

func (it *gormWrapper) Association(column string) *gorm.Association {
	return it.w.Association(column)
}

func (it *gormWrapper) Preload(column string, conditions ...interface{}) DBM {
	return Wrap(it.w.Preload(column, conditions...))
}

func (it *gormWrapper) Set(name string, value interface{}) DBM {
	return Wrap(it.w.Set(name, value))
}

func (it *gormWrapper) Get(name string) (interface{}, bool) {
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

func (it *gormWrapper) Take(out interface{}, where ...interface{}) DBM {
	return Wrap(it.w.Take(out, where...))
}
