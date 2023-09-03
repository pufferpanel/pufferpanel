package pufferpanel

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUpdateScopesWhereGranted(t *testing.T) {
	type args struct {
		source  []*Scope
		desired []*Scope
		changer []*Scope
	}
	tests := []struct {
		name string
		args args
		want []*Scope
	}{
		{
			name: "Add all",
			args: args{
				source:  []*Scope{},
				desired: allScopes,
				changer: allScopes,
			},
			want: allScopes,
		},
		{
			name: "Have none",
			args: args{
				source:  allScopes,
				desired: allScopes,
				changer: []*Scope{},
			},
			want: allScopes,
		},
		{
			name: "Add one",
			args: args{
				source:  []*Scope{},
				desired: []*Scope{ScopeServerList},
				changer: []*Scope{ScopeServerList},
			},
			want: []*Scope{ScopeServerList},
		},
		{
			name: "Remove one",
			args: args{
				source:  []*Scope{ScopeServerList, ScopeServerStart},
				desired: []*Scope{ScopeServerList},
				changer: []*Scope{ScopeServerList, ScopeServerStart},
			},
			want: []*Scope{ScopeServerList},
		},
		{
			name: "Add but don't have perms",
			args: args{
				source:  []*Scope{ScopeServerList, ScopeServerSftp},
				desired: []*Scope{ScopeServerAdmin},
				changer: []*Scope{ScopeServerReload},
			},
			want: []*Scope{ScopeServerList, ScopeServerSftp},
		},
		{
			name: "Add and source has others",
			args: args{
				source:  []*Scope{ScopeServerList, ScopeServerSftp},
				desired: []*Scope{ScopeServerLogs},
				changer: []*Scope{ScopeServerLogs},
			},
			want: []*Scope{ScopeServerList, ScopeServerSftp, ScopeServerLogs},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, UpdateScopesWhereGranted(tt.args.source, tt.args.desired, tt.args.changer), "UpdateScopesWhereGranted(%v, %v, %v)", tt.args.source, tt.args.desired, tt.args.changer)
		})
	}
}
