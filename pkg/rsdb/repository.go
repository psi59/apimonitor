package rsdb

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/VividCortex/mysqlerr"
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/lalaworks/kmessenger-webserver-go/pkg/lalalog"
	"github.com/pkg/errors"

	sierrors "github.com/realsangil/apimonitor/pkg/rserrors"
	"github.com/realsangil/apimonitor/pkg/rsmodel"
)

const (
	ErrRecordNotFound       sierrors.Error = "record not found"
	ErrDuplicateData        sierrors.Error = "duplicated data"
	ErrInvalidData          sierrors.Error = "invalid data"
	ErrInvalidModel         sierrors.Error = "invalid model"
	ErrForeignKeyConstraint sierrors.Error = "foreign key constraint fails"
)

type Repository interface {
	Create(tx Transaction, src rsmodel.ValidatedObject) error
	FirstOrCreate(tx Transaction, src rsmodel.ValidatedObject) error
	DeleteById(tx Transaction, id rsmodel.ValidatedObject) error
	GetById(tx Transaction, id rsmodel.ValidatedObject) error
	Save(tx Transaction, src rsmodel.ValidatedObject) error
	Patch(tx Transaction, src rsmodel.ValidatedObject, data rsmodel.ValidatedObject) error
	List(tx Transaction, items interface{}, filter ListFilter, orders Orders) (totalCount int, err error)
	CreateTable(tx Transaction) error
}

type DefaultRepository struct{}

func (repo DefaultRepository) CreateTable(tx Transaction) error {
	return nil
}

func checkItems(items interface{}) error {
	typeOfItems := reflect.TypeOf(items)
	srcTypeErr := errors.Wrap(ErrInvalidData, "src must be pointer of slice or array")
	if typeOfItems.Kind() != reflect.Ptr {
		return srcTypeErr
	}
	switch typeOfItems.Elem().Kind() {
	case reflect.Array, reflect.Slice:
	default:
		return srcTypeErr
	}
	return nil
}

func (repo *DefaultRepository) Patch(tx Transaction, src rsmodel.ValidatedObject, data rsmodel.ValidatedObject) error {
	if !data.IsValidated() {
		return ErrInvalidModel
	}
	err := tx.Tx().Model(src).Updates(data).Error
	return HandleSQLError(err)
}

func (repo *DefaultRepository) Save(tx Transaction, src rsmodel.ValidatedObject) error {
	err := tx.Tx().Save(src).Error
	return HandleSQLError(err)
}

func (repo *DefaultRepository) DeleteById(tx Transaction, id rsmodel.ValidatedObject) error {
	err := tx.Tx().Debug().Delete(id).Error
	return HandleSQLError(err)
}

func (repo *DefaultRepository) GetById(tx Transaction, id rsmodel.ValidatedObject) error {
	err := tx.Tx().First(id).Error
	return HandleSQLError(err)
}

func (repo *DefaultRepository) Create(tx Transaction, src rsmodel.ValidatedObject) error {
	if !src.IsValidated() {
		return ErrInvalidModel
	}
	if err := tx.Tx().Create(src).Error; err != nil {
		return HandleSQLError(err)
	}
	return nil
}

func (repo *DefaultRepository) FirstOrCreate(tx Transaction, src rsmodel.ValidatedObject) error {
	if !src.IsValidated() {
		return ErrInvalidModel
	}
	if err := tx.Tx().FirstOrCreate(src).Error; err != nil {
		return HandleSQLError(err)
	}
	return nil
}

func (repo *DefaultRepository) List(tx Transaction, items interface{}, filter ListFilter, orders Orders) (totalCount int, err error) {
	if e := checkItems(items); e != nil {
		err = errors.WithStack(e)
		return
	}

	query := tx.Tx()

	if filter.Conditions != nil {
		query = query.Where(filter.Conditions)
	}

	if e := query.Model(items).Count(&totalCount).Error; e != nil {
		err = HandleSQLError(e)
		return
	}

	if orders != nil {
		query = query.Order(orders.String())
	}

	if e := query.Offset(filter.NumItem * (filter.Page - 1)).Limit(filter.NumItem).Find(items).Error; e != nil {
		err = HandleSQLError(e)
		return
	}

	return
}

func NewDefaultRepository() Repository {
	return &DefaultRepository{}
}

func HandleSQLError(err error) error {
	if err == nil {
		return nil
	}
	if err == gorm.ErrRecordNotFound {
		return ErrRecordNotFound
	}
	mysqlError, ok := err.(*mysql.MySQLError)
	if ok {
		switch mysqlError.Number {
		case mysqlerr.ER_DUP_ENTRY:
			return ErrDuplicateData
		case mysqlerr.ER_DATA_TOO_LONG:
			return ErrInvalidData
		case mysqlerr.ER_NO_REFERENCED_ROW_2:
			return ErrForeignKeyConstraint
		}
	}

	_, fn, line, _ := runtime.Caller(1)
	lalalog.Errorf("repository error=%v, function=%v, line=%v", err, fn, line)
	return err
}

type ListFilter struct {
	Page       int
	NumItem    int
	Conditions map[string]interface{}
}

type Orders []Order

func (orders Orders) String() string {
	ordersStringArray := make([]string, 0)
	for _, o := range orders {
		ordersStringArray = append(ordersStringArray, o.String())
	}
	return strings.Join(ordersStringArray, ", ")
}

type Order struct {
	Field string
	IsASC bool
}

func (o Order) String() string {
	if o.IsASC {
		return fmt.Sprintf("%s ASC", o.Field)
	} else {
		return fmt.Sprintf("%s DESC", o.Field)
	}
}

func CreateTables(repos ...Repository) error {
	tx, err := GetConnection().Begin()
	if err != nil {
		return errors.WithStack(err)
	}

	for _, r := range repos {
		if err := r.CreateTable(tx); err != nil {
			return errors.WithStack(err)
		}
	}

	if err := tx.Commit(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
