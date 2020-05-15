package main

import (
	"data_preparor/preparor"
)

func main() {

	dataPreparor := preparor.NewPreparor()
	dataPreparor.Process("/home/musa/data/images/cat_and_dogs/dogs-vs-cats/train" , "./data" ,  50 , 2)
}
