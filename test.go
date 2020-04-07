package main

import ("fmt" 
		//"chord/src/chord")

func main() {
	ipRing, ipNot := CheckRing()
	fmt.Println("ip in Ring: ", ipRing)
	fmt.Println("Ip not in ring: ", ipNot)
}
