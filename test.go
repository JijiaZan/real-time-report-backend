package main

import (
	"fmt"
	"log"
)


func haha(ids ...int) {
	if len(ids) > 0 {
		fmt.Printf("xx")
	}
	var a string
	a = "s"
	a += "sb" + "haha"
	log.Fatal("exit")
	fmt.Println(a)
}

func main() {
	haha()
	haha(1,2)
}