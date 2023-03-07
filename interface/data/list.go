package data

type ListConsumer func(i int, v interface{}) bool

type ListExpected func(a interface{}) bool

type List interface {
	Add(val interface{})
	Get(index int) (val interface{})
	Set(index int, val interface{})
	Insert(index int, val interface{})
	Remove(index int) (val interface{})
	RemoveLast() (val interface{})
	Size() int
	ForEach(consumer ListConsumer)
	Contains(expected ListExpected) bool
	Range(start int, stop int) []interface{}
}
