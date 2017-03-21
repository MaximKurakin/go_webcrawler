/**
	Go Multi Thread Web Parser. 
	Main file
*/

package main

import (
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	fmt.Println("Starting Web server... ", time.Since(start))
	simple_http("65123")

	fmt.Println("Web server started ", time.Since(start))
	for ; ; {
	
	}
}