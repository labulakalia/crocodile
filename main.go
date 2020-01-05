package main

import "fmt"


func main() {
	a := 0
	for i:=0;i<10;i++ {
		a = a ^ 1
		fmt.Println(a)
	}
}