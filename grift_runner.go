
package main

import _ "github.com/pufferpanel/pufferpanel/grifts"
import "os"
import "log"
import "github.com/markbates/grift/grift"

func main() {
	grift.CommandName = "buffalo task"
	err := grift.Exec(os.Args[1:], false)
	if err != nil {
		log.Fatal(err)
	}
}