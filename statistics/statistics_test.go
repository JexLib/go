package statistics

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_getNote(t *testing.T) {
	type args struct {
		tag []string
	}
	tests := []struct {
		name string
		args args
		want StatisticsNote
	}{
		// TODO: Add test cases.
		{name: "",
			args: args{
				tag: []string{"sadfa", "sadf", "333"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getNote(tt.args.tag...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getNote() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewStatistics(t *testing.T) {
	type args struct {
		tag []string
	}
	tests := []struct {
		name string
		args args
		want StatisticsNote
	}{
		// TODO: Add test cases.
		{name: "",
			args: args{
				tag: []string{"ssss", "ddd", "111", "222"},
			},
		},
	}
	for i, tt := range tests {
		fmt.Println(i, len(tests))
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println(tt)
			if got := NewStatistics(tt.args.tag...); !reflect.DeepEqual(got, tt.want) {
				fmt.Println(got)
				t.Errorf("NewStatistics() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetStatistics(t *testing.T) {
	tests := []struct {
		name string
		want map[string]interface{}
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetStatistics(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetStatistics() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStatisticsNote_Add(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		st   StatisticsNote
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.st.Add(tt.args.key)
		})
	}
}

func TestStatistics_Clean(t *testing.T) {
	tests := []struct {
		name string
		st   *Statistics
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.st.Clean()
		})
	}
}

func TestStatistics_write(t *testing.T) {
	tests := []struct {
		name string
		st   *Statistics
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.st.write()
		})
	}
}

func TestStatistics_read(t *testing.T) {
	tests := []struct {
		name    string
		st      *Statistics
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.st.read(); (err != nil) != tt.wantErr {
				t.Errorf("Statistics.read() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
