package gordap

import (
	"reflect"
	"testing"
)

func TestBootstrap(t *testing.T) {
	type args struct {
		resourceType BootstrapType
	}
	tests := []struct {
		name    string
		args    args
		want    *BootstrapRecord
		wantErr bool
	}{
		{
			"one",
			args{IPV4},
			nil,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := bootstrap(tt.args.resourceType)
			if (err != nil) != tt.wantErr {
				t.Errorf("Bootstrap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			//var ipaddr = net.IPv4(173, 70, 134, 162)
			//log.Println(got.GetEndpoints(ipaddr).HTTPS)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Bootstrap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBootstrapIP(t *testing.T) {
	type args struct {
		ipaddress string
	}
	tests := []struct {
		name    string
		args    args
		want    Refs
		wantErr bool
	}{
		{
			"one",
			args{"173.70.134.162"},
			Refs{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BootstrapIP(tt.args.ipaddress)
			if (err != nil) != tt.wantErr {
				t.Errorf("BootstrapIP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BootstrapIP() = %v, want %v", got, tt.want)
			}
		})
	}
}
