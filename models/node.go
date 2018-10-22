/*
 Copyright 2018 Padduck, LLC
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

type Node struct {
	ID          uint   `gorm:"PRIMARY_KEY;AUTO_INCREMENT" json:"-"`
	Name        string `gorm:"size:100;UNIQUE;NOT NULL" json:"-"`
	PublicHost  string `gorm:"size:100;NOT NULL" json:"-"`
	PrivateHost string `gorm:"size:100;NOT NULL" json:"-"`
	PublicPort  int    `gorm:"DEFAULT:5656;NOT NULL" json:"-"`
	PrivatePort int    `gorm:"DEFAULT:5656;NOT NULL" json:"-"`
	SFTPPort    int    `gorm:"DEFAULT:5657;NOT NULL" json:"-"`

	//CreatedAt time.Time `json:"-"`
	//UpdatedAt time.Time `json:"-"`
}

type Nodes []*Node
