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

package services

import (
	"github.com/pufferpanel/pufferpanel/v3/logging"
)

func SendEmailViaDebug(to, subject, body string, async bool) error {
	logging.Debug.Println("DEBUG EMAIL TO " + to + "\n" + subject + "\n" + body)
	return nil
}
