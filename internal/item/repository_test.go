package item_test

import (
	"strings"
	"testing"

	"github.com/pkg/errors"

	"github.com/texazcowboy/warehouse/internal/foundation/database"

	. "github.com/texazcowboy/warehouse/internal/item"

	"github.com/stretchr/testify/assert"

	"github.com/DATA-DOG/go-sqlmock"
)

// todo: implement tests for all methods

func TestRepository_Create(t *testing.T) {

	tests := []struct {
		name         string
		mockFunc     func(mock sqlmock.Sqlmock)
		givenData    Item
		expectedData Item
		wantErr      bool
	}{
		{
			name: "ok",
			mockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectPrepare("INSERT INTO item\\(name\\)").
					ExpectQuery().
					WithArgs("test").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).
						AddRow(1))
				mock.ExpectCommit()
			},
			givenData:    Item{Name: "test"},
			expectedData: Item{ID: 1, Name: "test"},
		},
		{
			name:      "constraint violation",
			givenData: Item{Name: strings.Repeat("a", 51)},
			mockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectPrepare("INSERT INTO item\\(name\\)").
					ExpectQuery().
					WithArgs(strings.Repeat("a", 51)).
					WillReturnError(database.NewSQLError(errors.New("constraint violation")))
				mock.ExpectCommit()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup mock
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Error %s was not expected during stud db creation", err)
			}
			defer db.Close()

			tt.mockFunc(mock)

			// given
			givenRepo := NewRepository(db)

			// when
			err = givenRepo.Create(&tt.givenData)

			// then
			if err != nil && tt.wantErr {
				if err, ok := err.(*database.SQLError); !ok {
					t.Errorf("got unexpected error %T: %v", err, err)
				}
				// todo: add assertions regarding error data
				return
			}
			if err != nil && !tt.wantErr {
				t.Fatalf("error not expected: %v", err)
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedData.Name, tt.givenData.Name)
			assert.Equal(t, tt.expectedData.ID, tt.givenData.ID)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestRepository_GetByID(t *testing.T) {

	tests := []struct {
		name        string
		id          int64
		mockClosure func(mock sqlmock.Sqlmock)
		want        *Item
		wantErr     bool
	}{
		{
			name: "ok",
			id:   1,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectPrepare("SELECT \\* FROM item").
					ExpectQuery().
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
						AddRow(1, "test"))
				mock.ExpectCommit()
			},
			want: &Item{
				ID:   1,
				Name: "test",
			},
		},
		{
			name: "not found",
			id:   1,
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectPrepare("SELECT \\* FROM item").
					ExpectQuery().
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name"}))
				mock.ExpectCommit()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// given
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Error %s was not expected during stud db creation", err)
			}
			defer db.Close()

			tt.mockClosure(mock)
			givenRepo := NewRepository(db)

			// when
			result, err := givenRepo.GetByID(tt.id)

			// then
			if err != nil && tt.wantErr {
				err, ok := err.(*database.SQLError)
				if !ok {
					t.Errorf("got unexpected error %T: %v", err, err)
				}
				// todo: implement
				//assert.Equal(t, "get", err.Op)
				//assert.Equal(t, "item", err.Entity)
				assert.Contains(t, []string{"sql: no rows in result set"}, err.Message)
				return
			}
			if err != nil && !tt.wantErr {
				t.Errorf("error not expected: %v", err)
				return
			}

			assert.Equal(t, tt.want.ID, result.ID)
			assert.Equal(t, tt.want.Name, result.Name)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestRepository_GetAll(t *testing.T) {
}

func TestRepository_Update(t *testing.T) {

}

func TestRepository_DeleteByID(t *testing.T) {
}
