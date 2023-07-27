package utils

import (
	"reflect"
	"testing"
)

func TestConcatBytes(t *testing.T) {
	type args struct {
		slices [][]byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
		{
			name: "case1",
			args: args{slices: [][]byte{{1}, {2, 3}, {4}}},
			want: []byte{1, 2, 3, 4},
		},
		{
			name: "case2",
			args: args{slices: [][]byte{{4}, {3, 2}, {1}}},
			want: []byte{4, 3, 2, 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConcatBytes(tt.args.slices); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConcatBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
