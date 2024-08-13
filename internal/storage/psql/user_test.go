package psql

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/awleory/medodstest/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestUsersCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("err while open mock db: %v", err)
	}
	defer db.Close()

	type args struct {
		ctx  context.Context
		user domain.User
	}

	r := NewUsers(db)

	tests := []struct {
		name    string
		mock    func()
		args    args
		wantErr bool
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				user: domain.User{
					Email:    "test@email.ru",
					Password: "test",
				},
			},
			mock: func() {
				mock.ExpectExec("INSERT INTO users").WillReturnResult(sqlmock.NewResult(0, 1))
			},
		},
		{

			name: "NOT OK",
			args: args{
				ctx:  context.Background(),
				user: domain.User{},
			},
			mock: func() {
				mock.ExpectExec("INSERT INTO users")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := r.Create(tt.args.ctx, tt.args.user)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
