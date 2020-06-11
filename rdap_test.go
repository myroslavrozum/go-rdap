package rdap

import (
	"net"
	"reflect"
	"testing"
)

func Test_query(t *testing.T) {
	type args struct {
		ipaddr string
	}
	tests := []struct {
		name    string
		args    args
		want    net.IP
		wantErr bool
	}{
		{
			"one",
			args{"173.70.134.162"},
			net.IP{173, 70, 134, 162},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := rdap(tt.args.ipaddr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Bootstrap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("query() = %v, want %v", got, tt.want)
			}
		})
	}
}
