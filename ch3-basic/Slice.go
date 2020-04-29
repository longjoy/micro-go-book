package main

import "fmt"

func main()  {


	source := [...]int{1,2,3}
	sli := source[0:1]

	fmt.Printf("sli value is %v\n", sli)
	fmt.Printf("sli len is %v\n", len(sli))
	fmt.Printf("sli cap is %v\n", cap(sli))

	sli[0] = 4
	fmt.Printf("sli value is %v\n", sli)
	fmt.Printf("source value is %v\n", source)


    sli = make([]int, 2, 4)
	fmt.Printf("sli value is %v\n", sli)
	fmt.Printf("sli len is %v\n", len(sli))
	fmt.Printf("sli cap is %v\n", cap(sli))

	ex := []int{1,2,3}
	fmt.Printf("ex value is %v\n", ex)
	fmt.Printf("ex len is %v\n", len(ex))
	fmt.Printf("ex cap is %v\n", cap(ex))


}
