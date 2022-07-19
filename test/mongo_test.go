package main

import (
	"MyGameServer/mogodb"
	"fmt"
	"testing"
)

func BenchmarkFunc1(b *testing.B) {
	helper := mogodb.NewMongoHelper("test", "Students")
	helper.Connect("mongodb://127.0.0.1:27017")
	for i := 0; i < b.N; i++ {
		_, err := helper.FindAll()
		if err != nil {
			fmt.Print(err)
		}
	}
}

func BenchmarkFunc2(b *testing.B) {
	helper := mogodb.NewMongoHelper("test", "Students")
	helper.Connect("mongodb://127.0.0.1:27017")
	for i := 0; i < b.N; i++ {
		_, err := helper.TestFind()
		if err != nil {
			fmt.Print(err)
		}
	}
}
