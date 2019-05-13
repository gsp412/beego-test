package main

import (
	"beego-test/models"
	"fmt"
)

func main()  {

	pwd, err := models.CreatePassword("12345678")
	if nil != err {
		fmt.Println("fail:", err)
		return
	} else {
		fmt.Println("succ:", pwd)
	}


	isMatch, err := models.MatchPassword("12345678", pwd)
	if nil != err {
		fmt.Println("fail:", err)
		return
	} else {
		fmt.Println("succ:", isMatch)
	}
}

