package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"

	"github.com/emmaly/openhab"
)

func main() {
	var baseURL string
	var nameMatch *regexp.Regexp
	var setValue *string
	if len(os.Args) > 1 {
		baseURL = os.Args[1]
		if len(os.Args) > 2 {
			re, err := regexp.Compile("^" + os.Args[2] + "$")
			if err != nil {
				log.Fatal(err)
			}
			nameMatch = re
			if len(os.Args) > 3 {
				setValue = &os.Args[3]
			}
		}
	} else {
		fmt.Printf("%s <baseURL> [<nameMatch> [<newValue>]]\n", os.Args[0])
		os.Exit(1)
	}

	oh := openhab.New(baseURL, nil)
	items, err := oh.Items(nil)
	if err != nil {
		log.Fatal(err)
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].Name < items[j].Name
	})
	for _, item := range items {
		if nameMatch == nil {
			fmt.Printf("%-40s\t%-30s\t%s\n", item.Name, item.Label, item.State)
		} else if nameMatch.Match([]byte(item.Name)) {
			if setValue != nil {
				fmt.Printf("%-40s\t%-30s\t%s → %s\n", item.Name, item.Label, item.State, *setValue)
				item.Set(*setValue, nil)
			} else {
				fmt.Printf("%-40s\t%-30s\t%s\n", item.Name, item.Label, item.State)
			}
		}
	}
}
