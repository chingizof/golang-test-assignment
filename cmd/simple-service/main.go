package main

import (
	"fmt"
	"os"
	"time"
)

type item struct {
	name           string
	price          int
	income_outcome string
	comment        string
	category       string
	date           string
}

func main() {
	var name string
	var price int             //сами высчитываем
	var income_outcome string //возмодно bool, расход либо доход
	var comment string
	var category string
	date := time.Now()

	fmt.Fscan(os.Stdin, &name)
	fmt.Fscan(os.Stdin, &price)
	fmt.Fscan(os.Stdin, &income_outcome)
	fmt.Fscan(os.Stdin, &comment)
	fmt.Fscan(os.Stdin, &category)
	fmt.Println(date.Day(), date.Month(), date.Year())
}
