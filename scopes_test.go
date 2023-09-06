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
				desired: []*Scope{ScopeServerView},
				changer: []*Scope{ScopeServerView},
			},
			want: []*Scope{ScopeServerView},
		},
		{
			name: "Remove one",
			args: args{
				source:  []*Scope{ScopeServerView, ScopeServerStart},
				desired: []*Scope{ScopeServerView},
				changer: []*Scope{ScopeServerView, ScopeServerStart},
			},
			want: []*Scope{ScopeServerView},
		},
		{
			name: "Add but don't have perms",
			args: args{
				source:  []*Scope{ScopeServerView, ScopeServerSftp},
				desired: []*Scope{ScopeServerAdmin},
				changer: []*Scope{ScopeServerReload},
			},
			want: []*Scope{ScopeServerView, ScopeServerSftp},
		},
		{
			name: "Add and source has others",
			args: args{
				source:  []*Scope{ScopeServerView, ScopeServerSftp},
				desired: []*Scope{ScopeServerConsole},
				changer: []*Scope{ScopeServerConsole},
			},
			want: []*Scope{ScopeServerView, ScopeServerSftp, ScopeServerConsole},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, UpdateScopesWhereGranted(tt.args.source, tt.args.desired, tt.args.changer), "UpdateScopesWhereGranted(%v, %v, %v)", tt.args.source, tt.args.desired, tt.args.changer)
		})
	}
}
