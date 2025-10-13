package main

import (
	"fmt"
)

func main() {
	teamNotActivatedList := make([]int, 0)
	//teamNotActivatedList = append(teamNotActivatedList, 1)
	fmt.Println("Group Count:", len(teamNotActivatedList))
	groupCount := (len(teamNotActivatedList) + 1) / 2
	fmt.Println("Group Count:", groupCount)
}
