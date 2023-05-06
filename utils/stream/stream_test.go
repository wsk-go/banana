package stream

import (
	"fmt"
	"testing"
)

func TestStream1(t *testing.T) {
	aaa := []string{"123", "456"}
	Of(aaa).
		Filter(func(s string) bool {
			return true
		}).ToList()
}

func TestStream(t *testing.T) {
	aaa := []string{"123", "456"}
	aa := Map(Of(aaa).
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

	aa := MapStream(
		Of(aaa).
			Filter(func(s string) bool {
				return s == "123" || s == "1234"
			}),
		func(in string) int {
			return 1
		},
	)

	fmt.Println(aa)
}

type User struct {
	id string
}

func TestGroupBy(t *testing.T) {
	aaa := []User{{id: "123"}, {id: "123"}, {id: "1234"}}
	aa := Group(Of(aaa),
		func(in User) string {
			return in.id
		},
	)
	fmt.Println(aa)
}
