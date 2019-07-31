package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/1lann/staticice"
)

func main() {
	c := staticice.NewClient(http.DefaultClient)
	res, err := c.Search(
		staticice.RegionAU,
		staticice.NewSearchQuery().Query("samsung 970 evo").MaxPrice(200),
	)
	if err != nil {
		panic(err)
	}

	d, _ := json.MarshalIndent(res, "", "\t")
	fmt.Println(string(d))
}
