package main

import "github.com/pufferpanel/pufferpanel/v2/shared/logging"

func main() {
	defer logging.Close()

	Execute()
}
