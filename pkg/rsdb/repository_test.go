package rsdb

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/realsangil/apimonitor/pkg/rsmodel"
)

type mockFunc func(m sqlmock.Sqlmock)

type MockModel struct {
	Id         int  `json:"id" gorm:"primary_key"`
	Count      int  `json:"count"`
	isValidate bool `gorm:"-"`
}

func (m *MockModel) IsValidated() bool {
	return m.isValidate
}

func (m *MockModel) Validate() error {
	m.isValidate = true
	return nil
}

func (m MockModel) TableName() string {
	return "test_models"
}

func createMockDB() (*gorm.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	gormDB, gerr := gorm.Open("test", db)
	if gerr != nil {
		return nil, nil, gerr
	}
	gormDB.LogMode(true)
	return gormDB, mock, nil
}

func assertMockDatabase(t *testing.T, mock sqlmock.Sqlmock) {
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %s", err)
	}
}

func TestHandleSQLError(t *testing.T) {
	errUnexpected := errors.New("unexpected")

	type args struct {
		err error
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "already exists",
			args: args{
				err: &mysql.MySQLError{
					Number: 1062,
				},
			},
			wantErr: ErrDuplicateData,
		},
		{
			name: "not exists",
			args: args{
				err: &mysql.MySQLError{
					Number: 1452,
				},
			},
			wantErr: ErrForeignKeyConstraint,
		},
		{
			name: "invalid data test",
			args: args{
				err: &mysql.MySQLError{
					Number: 1406,
				},
			},
			wantErr: ErrInvalidData,
		},
		{
			name: "record not found",
			args: args{
				err: gorm.ErrRecordNotFound,
			},
			wantErr: ErrRecordNotFound,
		},
		{
			name: "unexpected error",
			args: args{
				err: errUnexpected,
			},
			wantErr: errUnexpected,
		},
		{
			name: "nil error",
			args: args{
				err: nil,
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := HandleSQLError(tt.args.err)
			assert.Equal(t, errors.Cause(gotErr), errors.Cause(tt.wantErr))
		})
	}
}

func TestDefaultRepository_Create(t *testing.T) {
	type args struct {
		expectQuery func(mock sqlmock.Sqlmock)
		mockData    rsmodel.ValidatedObject
	}
	cases := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "Success Create mockModel",
			args: args{
				expectQuery: func(mock sqlmock.Sqlmock) {
					mock.ExpectBegin()
					mock.ExpectExec(`INSERT`).WithArgs(1, 2).WillReturnResult(sqlmock.NewResult(1, 1))
					mock.ExpectCommit()
				},
				mockData: &MockModel{Id: 1, Count: 2, isValidate: true},
			},
			wantErr: nil,
		},
		{
			name: "유효성 검사되지 않은 모델 에러",
			args: args{
				expectQuery: func(mock sqlmock.Sqlmock) {

				},
				mockData: &MockModel{Id: 1, Count: 2, isValidate: false},
			},
			wantErr: ErrInvalidModel,
		},
		{
			name: "SQL 에러",
			args: args{
				expectQuery: func(mock sqlmock.Sqlmock) {
					mock.ExpectBegin()
					mock.ExpectExec(`INSERT`).WithArgs(2, 2).WillReturnError(fmt.Errorf("oops"))
					mock.ExpectRollback()
				},
				mockData: &MockModel{Id: 2, Count: 2, isValidate: true},
			},
			wantErr: fmt.Errorf("oops"),
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			gormDB, mock, err := createMockDB()
			if err != nil {
				t.Fatalf("Failed to connect mock database: %v", err)
			}
			defer gormDB.Close()
			defer assertMockDatabase(t, mock)

			tx := NewTransaction(gormDB)

			tt.args.expectQuery(mock)

			repo := &DefaultRepository{}
			gotErr := repo.Create(tx, tt.args.mockData)
			assert.Equal(t, tt.wantErr, errors.Cause(gotErr))
		})
	}
}

func TestDefaultRepository_GetById(t *testing.T) {
	type args struct {
		id          rsmodel.ValidatedObject
		expectQuery func(mock sqlmock.Sqlmock)
	}
	tests := []struct {
		name    string
		args    args
		want    rsmodel.ValidatedObject
		wantErr error
	}{
		{
			name: "pass",
			args: args{
				id: &MockModel{Id: 1},
				expectQuery: func(mock sqlmock.Sqlmock) {
					rows := sqlmock.NewRows([]string{"id", "count"}).AddRow(1, 10)
					mock.ExpectQuery(`SELECT`).WithArgs(1).WillReturnRows(rows)
				},
			},
			want:    &MockModel{Id: 1, Count: 10},
			wantErr: nil,
		},
		{
			name: "not found",
			args: args{
				id: &MockModel{Id: 10},
				expectQuery: func(mock sqlmock.Sqlmock) {
					mock.ExpectQuery("SELECT").WithArgs(10).WillReturnError(gorm.ErrRecordNotFound)
				},
			},
			want:    &MockModel{Id: 10},
			wantErr: ErrRecordNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gormDB, mock, err := createMockDB()
			if err != nil {
				t.Fatal(err)
			}
			defer gormDB.Close()
			defer assertMockDatabase(t, mock)

			tt.args.expectQuery(mock)

			tx := NewTransaction(gormDB)

			repo := &DefaultRepository{}
			gotErr := repo.GetById(tx, tt.args.id)

			assert.Equal(t, tt.want, tt.args.id)
			assert.Equal(t, tt.wantErr, errors.Cause(gotErr))
		})
	}
}

func TestDefaultRepository_DeleteById(t *testing.T) {
	type args struct {
		id          rsmodel.ValidatedObject
		expectQuery func(mock sqlmock.Sqlmock)
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "pass",
			args: args{
				id: &MockModel{Id: 1},
				expectQuery: func(mock sqlmock.Sqlmock) {
					mock.ExpectBegin()
					mock.ExpectExec(`DELETE`).WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
					mock.ExpectCommit()
				},
			},
			wantErr: nil,
		},
		{
			name: "unexpected error",
			args: args{
				id: &MockModel{Id: 10},
				expectQuery: func(mock sqlmock.Sqlmock) {
					mock.ExpectBegin()
					mock.ExpectExec("DELETE").WithArgs(10).WillReturnError(fmt.Errorf("oops"))
					mock.ExpectRollback()
				},
			},
			wantErr: fmt.Errorf("oops"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gormDB, mock, err := createMockDB()
			if err != nil {
				t.Fatalf("Failed to connect mock database: %v", err)
			}
			defer gormDB.Close()
			defer assertMockDatabase(t, mock)

			tx := NewTransaction(gormDB)

			tt.args.expectQuery(mock)

			repo := &DefaultRepository{}
			gotErr := repo.DeleteById(tx, tt.args.id)

			assert.Equal(t, tt.wantErr, errors.Cause(gotErr))
		})
	}
}

func TestDefaultRepository_Save(t *testing.T) {
	type args struct {
		src         rsmodel.ValidatedObject
		expectQuery func(mock sqlmock.Sqlmock)
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "pass",
			args: args{
				src: &MockModel{Id: 1, Count: 100},
				expectQuery: func(mock sqlmock.Sqlmock) {
					mock.ExpectBegin()
					mock.ExpectExec(`UPDATE`).WithArgs(100, 1).WillReturnResult(sqlmock.NewResult(1, 1))
					mock.ExpectCommit()
				},
			},
			wantErr: nil,
		},
		{
			name: "invalid data",
			args: args{
				src: &MockModel{Id: 10, Count: 1000},
				expectQuery: func(mock sqlmock.Sqlmock) {
					mock.ExpectBegin()
					mock.ExpectExec("UPDATE").WithArgs(1000, 10).WillReturnError(&mysql.MySQLError{Number: 1406})
					mock.ExpectRollback()
				},
			},
			wantErr: ErrInvalidData,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gormDB, mock, err := createMockDB()
			if err != nil {
				t.Fatal(err)
			}
			defer gormDB.Close()
			defer assertMockDatabase(t, mock)

			tt.args.expectQuery(mock)

			tx := NewTransaction(gormDB)

			repo := &DefaultRepository{}
			gotErr := repo.Save(tx, tt.args.src)

			assert.Equal(t, tt.wantErr, errors.Cause(gotErr))
		})
	}
}

func TestDefaultRepository_Patch(t *testing.T) {
	type args struct {
		src         rsmodel.ValidatedObject
		data        rsmodel.ValidatedObject
		expectQuery func(mock sqlmock.Sqlmock)
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
		want    rsmodel.ValidatedObject
	}{
		{
			name: "pass",
			args: args{
				src:  &MockModel{Id: 1, Count: 10},
				data: &MockModel{Count: 1000, isValidate: true},
				expectQuery: func(mock sqlmock.Sqlmock) {
					mock.ExpectBegin()
					mock.ExpectExec("UPDATE").WithArgs(1000, 1).WillReturnResult(sqlmock.NewResult(1, 1))
					mock.ExpectCommit()
				},
			},
			want:    &MockModel{Id: 1, Count: 1000},
			wantErr: nil,
		},
		{
			name: "failed",
			args: args{
				src:  &MockModel{Id: 1, Count: 10},
				data: &MockModel{Count: 1000, isValidate: true},
				expectQuery: func(mock sqlmock.Sqlmock) {
					mock.ExpectBegin()
					mock.ExpectExec("UPDATE").WithArgs(1000, 1).WillReturnError(&mysql.MySQLError{Number: 1452})
					mock.ExpectRollback()
				},
			},
			want:    &MockModel{Id: 1, Count: 1000},
			wantErr: ErrForeignKeyConstraint,
		},
		{
			name: "유효성 검사에 실패한 data 값 오류",
			args: args{
				src:  &MockModel{Id: 1, Count: 10},
				data: &MockModel{Count: 1000, isValidate: false},
				expectQuery: func(mock sqlmock.Sqlmock) {
					//mock.ExpectExec("UPDATE").WithArgs(1000, 1).WillReturnError(&mysql.MySQLError{Number: 1452})
				},
			},
			want:    &MockModel{Id: 1, Count: 10},
			wantErr: ErrInvalidModel,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gormDB, mock, err := createMockDB()
			if err != nil {
				t.Fatal(err)
			}

			defer gormDB.Close()
			defer assertMockDatabase(t, mock)

			tt.args.expectQuery(mock)

			tx := NewTransaction(gormDB)

			repo := &DefaultRepository{}
			gotErr := repo.Patch(tx, tt.args.src, tt.args.data)

			assert.Equal(t, tt.wantErr, errors.Cause(gotErr))
			assert.Equal(t, tt.want, tt.args.src)
		})
	}
}

func TestDefaultRepository_List(t *testing.T) {
	type args struct {
		src         interface{}
		filter      ListFilter
		expectQuery func(mock sqlmock.Sqlmock)
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr error
	}{
		{
			name: "pass",
			args: args{
				src:    &[]MockModel{},
				filter: NewListFilter(1, 20, nil, FieldCondition{"site_id", 1}),
				expectQuery: func(mock sqlmock.Sqlmock) {
					countRows := sqlmock.NewRows([]string{"count(*)"}).AddRow(1)
					mock.ExpectQuery("SELECT count").WithArgs(1).WillReturnRows(countRows)

					rows := sqlmock.NewRows([]string{"id", "from"}).AddRow(1, 2)
					mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(rows)
				},
			},
			want: 1,
		},
		{
			name: "src not pointer",
			args: args{
				src:    []MockModel{},
				filter: NewListFilter(1, 20, nil, FieldCondition{"site_id", 1}),
				expectQuery: func(mock sqlmock.Sqlmock) {
				},
			},
			wantErr: ErrInvalidData,
		},
		{
			name: "src not array or slice",
			args: args{
				src:    &MockModel{},
				filter: NewListFilter(1, 20, nil, FieldCondition{"site_id", 1}),
				expectQuery: func(mock sqlmock.Sqlmock) {
				},
			},
			wantErr: ErrInvalidData,
		},
		{
			name: "unexpected sql error",
			args: args{
				src:    &[]MockModel{},
				filter: NewListFilter(1, 20, nil, FieldCondition{"site_id", 1}),
				expectQuery: func(mock sqlmock.Sqlmock) {
					countRows := sqlmock.NewRows([]string{"count(*)"}).AddRow(0)
					mock.ExpectQuery("SELECT count").WithArgs(1).WillReturnRows(countRows)
					mock.ExpectQuery("SELECT").WithArgs(1).WillReturnError(fmt.Errorf("sql error"))
				},
			},
			wantErr: fmt.Errorf("sql error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gormDB, mock, err := createMockDB()
			if err != nil {
				t.Fatal(err)
			}
			defer gormDB.Close()
			defer assertMockDatabase(t, mock)

			tt.args.expectQuery(mock)

			tx := NewTransaction(gormDB)

			repo := &DefaultRepository{}
			got, gotErr := repo.List(tx, tt.args.src, tt.args.filter)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, errors.Cause(gotErr))
		})
	}
}

func TestDefaultRepository_FirstOrCreate(t *testing.T) {
	type args struct {
		expectQuery func(mock sqlmock.Sqlmock)
		mockData    rsmodel.ValidatedObject
	}
	cases := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "Success Create mockModel",
			args: args{
				expectQuery: func(mock sqlmock.Sqlmock) {
					// rows := sqlmock.NewRows([]string{"id", "count"})
					// mock.ExpectBegin()
					mock.ExpectQuery("SELECT").WithArgs(1).WillReturnError(gorm.ErrRecordNotFound)
					mock.ExpectExec(`INSERT`).WithArgs(1, 2).WillReturnResult(sqlmock.NewResult(1, 1))
					mock.ExpectCommit()
				},
				mockData: &MockModel{Id: 1, Count: 2, isValidate: true},
			},
			wantErr: nil,
		},
		{
			name: "이미 존재하는 row",
			args: args{
				expectQuery: func(mock sqlmock.Sqlmock) {
					rows := sqlmock.NewRows([]string{"id", "count"}).AddRow(1, 2)
					mock.ExpectBegin()
					mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(rows)
					mock.ExpectCommit()
				},
				mockData: &MockModel{Id: 1, Count: 2, isValidate: true},
			},
			wantErr: nil,
		},
		{
			name: "유효성 검사되지 않은 모델 에러",
			args: args{
				expectQuery: func(mock sqlmock.Sqlmock) {

				},
				mockData: &MockModel{Id: 1, Count: 2, isValidate: false},
			},
			wantErr: ErrInvalidModel,
		},
		{
			name: "SQL 에러",
			args: args{
				expectQuery: func(mock sqlmock.Sqlmock) {
					rows := sqlmock.NewRows([]string{"id", "count"})
					mock.ExpectBegin()
					mock.ExpectQuery("SELECT").WithArgs(2).WillReturnRows(rows)
					mock.ExpectExec(`INSERT`).WithArgs(2, 2).WillReturnError(fmt.Errorf("oops"))
					mock.ExpectRollback()
				},
				mockData: &MockModel{Id: 2, Count: 2, isValidate: true},
			},
			wantErr: fmt.Errorf("oops"),
		},
		{
			name: "SQL 에러",
			args: args{
				expectQuery: func(mock sqlmock.Sqlmock) {
					mock.ExpectBegin()
					mock.ExpectQuery("SELECT").WithArgs(2).WillReturnError(fmt.Errorf("oops"))
					mock.ExpectRollback()
				},
				mockData: &MockModel{Id: 2, Count: 2, isValidate: true},
			},
			wantErr: fmt.Errorf("oops"),
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			gormDB, mock, err := createMockDB()
			if err != nil {
				t.Fatalf("Failed to connect mock database: %v", err)
			}
			defer gormDB.Close()
			defer assertMockDatabase(t, mock)

			tx := NewTransaction(gormDB)

			tt.args.expectQuery(mock)

			repo := &DefaultRepository{}
			gotErr := repo.FirstOrCreate(tx, tt.args.mockData)
			assert.Equal(t, tt.wantErr, errors.Cause(gotErr))
		})
	}
}
