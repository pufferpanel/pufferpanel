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
			name: "java 16",
			args: args{
				data: []byte(`5236:
 garbage-first heap   total 1086464K, used 973437K [0x0000000080000000, 0x0000000100000000)
  region size 1024K, 41 young (41984K), 4 survivors (4096K)
 Metaspace       used 385449K, committed 388928K, reserved 1441792K
  class space    used 47768K, committed 49472K, reserved 1048576K
`),
			},
			want: &JvmStats{HeapUsed: 973437 * 1024, HeapTotal: 1086464 * 1024, MetaspaceUsed: 385449 * 1024, MetaspaceTotal: 388928 * 1024},
		},
		{
			name: "java 21",
			args: args{
				data: []byte(`0:
 def new generation   total 104576K, used 36054K [0x0000000080000000, 0x0000000087170000, 0x00000000aaaa0000)
  eden space 92992K,  36% used [0x0000000080000000, 0x00000000821627e0, 0x0000000085ad0000)
  from space 11584K,  16% used [0x0000000085ad0000, 0x0000000085ca3138, 0x0000000086620000)
  to   space 11584K,   0% used [0x0000000086620000, 0x0000000086620000, 0x0000000087170000)
 tenured generation   total 232236K, used 145245K [0x00000000aaaa0000, 0x00000000b8d6b000, 0x0000000100000000)
   the space 232236K,  62% used [0x00000000aaaa0000, 0x00000000b3877760, 0x00000000b3877800, 0x00000000b8d6b000)
 Metaspace       used 79381K, committed 80384K, reserved 1179648K
  class space    used 12067K, committed 12480K, reserved 1048576K`),
			},
			want: &JvmStats{HeapUsed: (36054 + 145245) * 1024, HeapTotal: (104576 + 232236) * 1024, MetaspaceUsed: 79381 * 1024, MetaspaceTotal: 80384 * 1024},
		},
		{
			name: "test 2",
			args: args{
				data: []byte(`0:
 garbage-first heap   total 342016K, used 248394K [0x0000000080000000, 0x0000000100000000)
  region size 1024K, 98 young (100352K), 19 survivors (19456K)
 Metaspace       used 80690K, committed 81728K, reserved 1179648K
  class space    used 12159K, committed 12608K, reserved 1048576K
`),
			},
			want: &JvmStats{HeapUsed: 248394 * 1024, HeapTotal: 342016 * 1024, MetaspaceUsed: 80690 * 1024, MetaspaceTotal: 1179648 * 1024},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, ParseJCMDResponse(tt.args.data), "ParseJCMDResponse(%v)", tt.args.data)
		})
	}
}
