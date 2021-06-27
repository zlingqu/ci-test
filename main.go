package main

import (
	r "ci-test/router"
)

func main() {

	r := r.SetupRouter()
	r.Run(":" + "80") // listen and serve on 0.0.0.0:80
}
