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

package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/pufferpanel/pufferpanel/v2/database"
	"github.com/pufferpanel/pufferpanel/v2/logging"
	"github.com/pufferpanel/pufferpanel/v2/response"
	"net/http"
)

func NeedsDatabase(c *gin.Context) {
	db, err := database.GetConnection()

	if err != nil {
		logging.Error().Printf("Database not available: %s", err)
		err = pufferpanel.ErrDatabaseNotAvailable
	}

	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.Set("db", db)

	c.Next()
}

func GetDatabase(c *gin.Context) *gorm.DB {
	db, exist := c.Get("db")
	if !exist {
		return nil
	}
	casted, ok := db.(*gorm.DB)
	if !ok {
		return nil
	}
	return casted
}
