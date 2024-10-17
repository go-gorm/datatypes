package datatypes

import (
	"database/sql/driver"
	"reflect"
	"testing"
)

func TestNull_Scan(t *testing.T) {
	type args struct {
		value any
	}
	type testCase[T any] struct {
		name    string
		n       Null[T]
		args    args
		wantErr bool
	}
	tests := []testCase[int64]{
		{
			name:    "test",
			n:       Null[int64]{},
			args:    args{value: "test"},
			wantErr: true,
		}, {
			name:    "test2",
			n:       Null[int64]{},
			args:    args{value: "6"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.n.Scan(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNull_Value(t *testing.T) {
	type testCase[T any] struct {
		name    string
		n       Null[T]
		want    driver.Value
		wantErr bool
	}
	var (
		v1 int64 = 1
		v2 int64 = 2
	)
	tests := []testCase[int64]{
		{
			name:    "test",
			n:       Null[int64]{V: v1, Valid: true},
			want:    v1,
			wantErr: false,
		}, {
			name:    "test",
			n:       Null[int64]{V: v2, Valid: false},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Value() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNullInt64_Value(t *testing.T) {
	type testCase[T any] struct {
		name    string
		n       NullInt64
		want    driver.Value
		wantErr bool
	}
	var (
		v1 int64 = 1
		v2 int64 = 2
	)
	tests := []testCase[int64]{
		{
			name:    "test",
			n:       NullInt64{V: v1, Valid: true},
			want:    v1,
			wantErr: false,
		}, {
			name:    "test",
			n:       NullInt64{V: v2, Valid: false},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("Int64_Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Int64_Value() got = %v, want %v", got, tt.want)
			}
		})
	}
}
