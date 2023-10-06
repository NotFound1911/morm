package morm

import (
	errs "github.com/NotFound1911/morm/internal/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegistry_get(t *testing.T) {
	mp := func() any {
		val := &TestModel{}
		return &val
	}
	testCases := []struct {
		name      string
		val       any
		wantModel *model
		wantErr   error
	}{
		{
			name:    "test model",
			val:     TestModel{},
			wantErr: errs.NewErrPointerOnly(TestModel{}),
		},
		{
			name: "pointer",
			val:  &TestModel{},
			wantModel: &model{
				tableName: "test_model",
				fieldMap: map[string]*field{
					"Id": {
						colName: "id",
					},
					"FirstName": {
						colName: "first_name",
					},
					"Age": {
						colName: "age",
					},
					"LastName": {
						colName: "last_name",
					},
				},
			},
		},
		{
			name:    "multiple pointer",
			val:     mp,
			wantErr: errs.NewErrPointerOnly(mp),
		},
		{
			name:    "map",
			val:     map[string]any{},
			wantErr: errs.NewErrPointerOnly(map[string]any{}),
		},
		{
			name:    "slice",
			val:     []int{},
			wantErr: errs.NewErrPointerOnly([]int{}),
		},
		{
			name:    "basic type",
			val:     0,
			wantErr: errs.NewErrPointerOnly(0),
		},
	}
	r := &registry{}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := r.get(tc.val)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantModel, m)
		})
	}
}
