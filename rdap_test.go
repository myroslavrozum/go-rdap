package gordap

import (
	"net"
	"testing"
)

func Benchmark_query(b *testing.B) {
	type args struct {
		ipaddr string
	}
	test := struct {
		name string
		ip   string
	}{
		name: "one",
		ip:   "173.70.134.162",
		//ip: "178.210.203.50",
		//ip: "181.1.140.66",
		//ip: "193.19.84.177",
		//ip: "94.100.180.200",
	}

	for i := 0; i < b.N; i++ {
		Rdap2Vcard(test.ip)
	}
	b.ReportAllocs()
}

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
			name: "one",
			want: net.IP{173, 70, 134, 162},
			//want: net.IP{178, 210, 203, 50},
			//want: net.IP{181, 1, 140, 66},
			//want: net.IP{193, 19, 84, 177},
			//want:    net.IP{94, 100, 180, 200},
			wantErr: true,
			args:    args{"173.70.134.162"},
			//args: args{"178.210.203.50"},
			//args: args{"181.1.140.66"},
			//args: args{"193.19.84.177"},
			//args: args{"94.100.180.200"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Rdap2Vcard(tt.args.ipaddr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Bootstrap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			if !got.Equal(tt.want) {
				t.Errorf("query() = %v, want %v", got, tt.want)
			}
		})
	}
}
