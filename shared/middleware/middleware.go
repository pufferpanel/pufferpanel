/*
 Copyright 2019 Padduck, LLC
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
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferpanel/v2/shared"
	"github.com/pufferpanel/pufferpanel/v2/shared/logging"
	"github.com/pufferpanel/pufferpanel/v2/shared/response"
	"net/http"
	"runtime/debug"
)

func ResponseAndRecover(c *gin.Context) {
	defer func() {

		if err := recover(); err != nil {
			if _, ok := err.(error); !ok {
				err = errors.New(shared.ToString(err))
			}
			response.HandleError(c, err.(error), http.StatusInternalServerError)

			logging.Error("Error handling route\n%+v\n%s", err, debug.Stack())
			c.Abort()
		}
	}()

	c.Next()
}

func Recover(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			logging.Error("Error handling route\n%+v\n%s", err, debug.Stack())
			c.Abort()
		}
	}()

	c.Next()
}

func Database(c *gin.Context, db *sql.DB) {
	trans, err := db.Begin()
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := recover(); err != nil {
			_ = trans.Rollback()
			panic(err)
		} else {
			err = trans.Commit()
			if err != nil {
				panic(err)
			}
		}
	}()
	c.Set("database", trans)
	c.Next()
}
