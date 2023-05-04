package utils

import (
	"fmt"
	"testing"
)

func TestStream(t *testing.T) {
	aaa := []string{"123", "456"}
	aa := Map(Stream(aaa).
		Filter(func(s string) bool {
			return s == "123"
		}),
		func(in string) int {
			return 1
		},
	)
	fmt.Println(aa)
}

func TestMapStream(t *testing.T) {
	aaa := []string{"123", "123", "1234", "456"}
	aa := MapStream(Stream(aaa).
		Filter(func(s string) bool {
			return s == "123" || s == "1234"
		}),
		func(in string) int {
			return 1
		},
	).Filter(func(i int) bool {
		return i > 1
	}).ToList()

	fmt.Println(aa)
}

type User struct {
	id string
}

func TestGroupBy(t *testing.T) {
	aaa := []User{{id: "123"}, {id: "123"}, {id: "1234"}}
	aa := Group(Stream(aaa),
		func(in User) string {
			return in.id
		},
	)
	fmt.Println(aa)
}
