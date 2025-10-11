package main

import (
	"fmt"
	"sk_mall/utils"
	"time"
)

func main(){
	
	c,_:=utils.NewSafeSnowFlakeCreater(0,3000*time.Millisecond)
	fmt.Println(c.Generate())

}
