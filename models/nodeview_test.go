/*
 Copyright 2020 Padduck, LLC
  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at
  	http://www.apache.org/licenses/LICENSE-2.0
  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package models

import (
	"reflect"
	"testing"
)

func TestFromNodes(t *testing.T) {
	type args struct {
		n *Nodes
	}

	sourceNode := make(Nodes, 1)
	sourceNode[0] = &Node{
		ID:          5,
		Name:        "node",
		PublicHost:  "localhost",
		PrivateHost: "127.0.0.1",
		PublicPort:  8080,
		PrivatePort: 5658,
		SFTPPort:    5657,
		Secret:      "somesecret",
	}

	desired := make(NodesView, 1)
	desired[0] = &NodeView{
		Id:          sourceNode[0].ID,
		Name:        sourceNode[0].Name,
		PublicHost:  sourceNode[0].PublicHost,
		PrivateHost: sourceNode[0].PrivateHost,
		PublicPort:  sourceNode[0].PublicPort,
		PrivatePort: sourceNode[0].PrivatePort,
		SFTPPort:    sourceNode[0].SFTPPort,
	}

	tests := []struct {
		name string
		args args
		want *NodesView
	}{
		{
			name: "conversion",
			args: args{
				n: &sourceNode,
			},
			want: &desired,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FromNodes(tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromNodes() = %v, want %v", got, tt.want)
			}
		})
	}
}
