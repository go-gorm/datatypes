package datatypes

import (
	"database/sql/driver"
	"math"
	"reflect"
	"testing"
	"time"
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

func TestNullIntegerParameters(t *testing.T) {
	tests := []struct {
		value driver.Valuer
		want  int64
	}{
		{NewNull(int(-1)), -1},
		{NewNull(int8(math.MinInt8)), math.MinInt8},
		{NullInt16{V: math.MinInt16, Valid: true}, math.MinInt16},
		{NullInt32{V: math.MaxInt32, Valid: true}, math.MaxInt32},
		{NullByte{V: math.MaxUint8, Valid: true}, math.MaxUint8},
		{NewNull(uint16(math.MaxUint16)), math.MaxUint16},
		{NewNull(uint32(math.MaxUint32)), math.MaxUint32},
		{NewNull(uint(0)), 0},
		{NewNull(uint64(math.MaxInt64)), math.MaxInt64},
	}
	for _, tt := range tests {
		got, err := driver.DefaultParameterConverter.ConvertValue(tt.value)
		if err != nil || got != tt.want {
			t.Errorf("binding %T: got (%v, %v), want (%d, nil)", tt.value, got, err, tt.want)
		}
	}
}

func TestNullPreservesDriverSpecificValues(t *testing.T) {
	for _, value := range []any{uint64(math.MaxUint64), time.Minute, []string{"paid", "pending"}} {
		got, err := NewNull(value).Value()
		if err != nil || !reflect.DeepEqual(got, value) {
			t.Errorf("Value() for %T: got (%#v, %v), want (%#v, nil)", value, got, err, value)
		}
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
