package config

import (
	"reflect"
	"testing"
)

func Test_firstNonNil(t *testing.T) {
	a := 1
	b := 2
	c := 3
	type args struct {
		maybeNil []interface{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{"default", args{[]interface{}{nil, &b, &c}}, &b},
		{"full", args{[]interface{}{&a, &b, &c}}, &a},
		{"empty", args{[]interface{}{nil, nil, nil}}, nil},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := firstNonNil(tt.args.maybeNil...); !reflect.DeepEqual(
					got, tt.want,
				) {
					t.Errorf("firstNonNil() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
