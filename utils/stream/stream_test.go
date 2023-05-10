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

type FlatMapValue struct {
	Code int
	List []string
}

//func TestFlatMapStream(t *testing.T) {
//	aaa := []FlatMapValue{{Code: 1, List: []string{"123", "456"}}, {Code: 1, List: []string{"666", "888"}}}
//	list := FlatMapStream(Of(aaa), func(val FlatMapValue) []string {
//		return val.List
//	}).ToList()
//	fmt.Println(list)
//}

func TestStream_ForEach(t *testing.T) {
	aaa := []string{"111", "222", "333", "444"}
	Of(aaa).ForEach(func(s string) {
		println(s)
	})
}

func TestStream_Reduce(t *testing.T) {
	aaa := []int{100, 200, 500, 600}
	v := Of(aaa).Reduce(func(all int, current int) int {
		return all + current
	})

	println(v)
}

func TestMapStream(t *testing.T) {

	aaa := []string{"111", "222", "333", "444"}

	aa := MapStream(
		Of(aaa),
		func(in string) string {
			return in + "hello"
		},
	).Distinct()

	fmt.Println(aa)
}

func TestDistinct(t *testing.T) {

	aaa := []string{"123", "123", "1234", "456"}

	aa := MapStream(
		Of(aaa).
			Filter(func(s string) bool {
				return true
			}),
		func(in string) string {
			return in
		},
	).Distinct().ToList()

	fmt.Println(aa)
}

func TestDistinctByKey(t *testing.T) {

	aaa := []User{{id: "123"}, {id: "123"}, {id: "1234"}}

	aa := Of(aaa).
		DistinctByKey(func(s User) any {
			return s.id
		}).ToList()

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
