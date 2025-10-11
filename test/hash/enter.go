package main

import (
	"fmt"
	"sk_mall/utils"
)

func main() {
	str := "123234515"
	hashstr := utils.GetHashStr(str)
	fmt.Println(hashstr)
	
	fmt.Println(utils.CheckHashStr(str,"47829b2a13c39a81cfa2d7c50bb6b6a063354a410e647889f3debd851232a641"))
}


