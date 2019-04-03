package main

import "fmt"

var set = []string{
	"Hom", "Het", "Hemi", "NA",
}

var sep = ";"
var dep = 3

func main() {
	loop("", dep-1)
}

func loop(str string, i int) {
	if str != "" {
		str = str + sep
	}
	if i == 0 {
		for t := range set {
			fmt.Println(str + set[t])
		}
	} else {
		for t := range set {
			loop(str+set[t], i-1)
		}
	}
}
