package pufferpanel

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseJCMDResponse(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name string
		args args
		want *JvmStats
	}{
		{
			name: "expected",
			args: args{
				data: []byte(`5236:
 garbage-first heap   total 1086464K, used 973437K [0x0000000080000000, 0x0000000100000000)
  region size 1024K, 41 young (41984K), 4 survivors (4096K)
 Metaspace       used 385449K, committed 388928K, reserved 1441792K
  class space    used 47768K, committed 49472K, reserved 1048576K
`),
			},
			want: &JvmStats{HeapUsed: 973437 * 1024, HeapTotal: 1086464 * 1024, MetaspaceUsed: 385449 * 1024},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, ParseJCMDResponse(tt.args.data), "ParseJCMDResponse(%v)", tt.args.data)
		})
	}
}
