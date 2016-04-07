package assert

import (
	. "github.com/claassen/gotest"
)

type TestObj struct {
	i int
	s string
}

func Test() {

	Describe("When using IsNil and IsNotNil", func() {
		
		//Nil
		It("nil is nil", func() {
			AssertThat(nil).IsNil()
		})

		//Bool
		It("true is not nil", func() {
			AssertThat(true).IsNotNil()
		})

		It("false is not nil", func() {
			AssertThat(false).IsNotNil()
		})

		It("an uninitialized bool is not nil", func() {
			var x bool

			AssertThat(x).IsNotNil()
		})

		//Number
		It("zero is not nil", func() {
			AssertThat(0).IsNotNil()
		})

		//String
		It("a string is not nil", func() {
			AssertThat("abc").IsNotNil()
		})

		It("an uninitialized string is not nil", func() {
			var x string

			AssertThat(x).IsNotNil()
		})

		//Struct
		It("a struct is not nil", func() {
			x := TestObj{}

			AssertThat(x).IsNotNil()
		})

		It("an uninitialized struct is not nil", func() {
			var x TestObj

			AssertThat(x).IsNotNil()
		})
		
		//Channel
		It("a chan is not nil", func(){
			x := make(chan int)

			AssertThat(x).IsNotNil()
		})

		It("an uninitialized chan is nil", func() {
			var x chan int

			AssertThat(x).IsNil()
		})

		It("a chan assigned to nil is nill", func() {
			x := make(chan int)
			x = nil

			AssertThat(x).IsNil()
		})

		//Func
		It("a func is not nil", func(){
			x := func() int { return 42 }

			AssertThat(x).IsNotNil()
		})

		It("an uninitialized func is nil", func() {
			var x func()

			AssertThat(x).IsNil()
		})

		It("a func assigned to nil is nill", func() {
			x := func() int { return 42 }
			x = nil

			AssertThat(x).IsNil()
		})

		//Interface
		It("an interface is not nil", func(){
			var x interface{}
			x = 5

			AssertThat(x).IsNotNil()
		})

		It("an uninitialized interface is nil", func() {
			var x interface{}

			AssertThat(x).IsNil()
		})

		It("an interface assigned to nil is nill", func() {
			var x interface{}
			x = nil

			AssertThat(x).IsNil()
		})

		//Map
		It("a map is not nil", func(){
			x := make(map[string]int)

			AssertThat(x).IsNotNil()
		})

		It("an uninitialized map is nil", func() {
			var x map[string]int

			AssertThat(x).IsNil()
		})

		It("a map assigned to nil is nill", func() {
			x := make(map[string]int)
			x = nil

			AssertThat(x).IsNil()
		})

		//Array
		It("an array is not nil", func() {
			x := make([]int, 10)

			AssertThat(x).IsNotNil()
		})

		It("an uninitialized array is nil", func() {
			var x []int

			AssertThat(x).IsNil()
		})

		It("an array assigned to nil is nil", func() {
			x := make([]int, 10)
			x = nil

			AssertThat(x).IsNil()
		})

		//Slice
		It("a slice is not nil", func(){
			x := make([]int, 10)
			y := x[:1]

			AssertThat(y).IsNotNil()
		})

		//Uninitialized slice?

		It("a slice assigned to nil is nill", func() {
			x := make([]int, 10)
			y := x[:1]
			y = nil

			AssertThat(y).IsNil()
		})

		//Ptr
		It("a pointer to a struct is not nil", func() {
			x := &TestObj{}

			AssertThat(x).IsNotNil()
		})

		It("a pointer to a struct assigned nil is nil", func(){
			x := &TestObj{}
			x = nil

			AssertThat(x).IsNil()
		})
	})

	Describe("When using IsEqualTo and IsNotEqualTo", func() {

		//Nil already covered as IsNil and IsNotNil just use IsEqualTo under the covers

		//Bool
		It("an uninitialized bool is false", func() {
			var x bool

			AssertThat(x).IsEqualTo(false)
		})

		It("true equals true", func() {
			AssertThat(true).IsEqualTo(true)
		})

		It("true does not equal false", func() {
			AssertThat(true).IsNotEqualTo(false)
		})

		//Number
		It("an uninitialized int is equal to 0", func() {
			var x int

			AssertThat(x).IsEqualTo(0)
		})

		It("an uninitialized float is equal to 0", func() {
			var x float32

			AssertThat(x).IsEqualTo(0.0)
		})

		It("a number is equal to itself", func() {
			AssertThat(42).IsEqualTo(42)
		})

		It("an int and a float of the same value are equal", func() {
			AssertThat(1.0).IsEqualTo(1)
		})

		It("two different number are not equal", func() {
			AssertThat(42).IsNotEqualTo(123)
		})

		//String
		It("an uninitialized string is equal to an empty string", func() {
			var x string
			y := ""

			AssertThat(x).IsEqualTo(y)
		})

		It("a string is equal to the same string", func() {
			x := "abc"
			y := "abc"

			AssertThat(x).IsEqualTo(y)
		})

		//Struct
		It("two different equivalent structs are equal", func() {
			x := TestObj{i: 1, s: "abc"}
			y := TestObj{i: 1, s: "abc"}

			AssertThat(x).IsEqualTo(y)
		})

		It("two different structs are not equal", func() {
			x := TestObj{i: 1, s: "abc"}
			y := TestObj{i: 1, s: "xyz"}

			AssertThat(x).IsNotEqualTo(y)
		})

		It("an uninitialized struct is equal to a default constructed struct", func() {
			var x TestObj
			y := TestObj{}

			AssertThat(x).IsEqualTo(y)
		})

		//Channel
		It("two references to the same channel are equal", func() {
			x := make(chan int)
			y := x

			AssertThat(x).IsEqualTo(y)
		})

		It("two references to two different channels are not equal", func() {
			x := make(chan int)
			y := make(chan int)

			AssertThat(x).IsNotEqualTo(y)
		})

		//Func
		It("two even equivalent funcs are not equal", func() {
			x := func() int { return 42 }
			y := func() int { return 42 }

			AssertThat(x).IsNotEqualTo(y)
		})

		It("two different funcs are not equal", func() {
			x := func() int { return 1 }
			y := func() int { return 2 }

			AssertThat(x).IsNotEqualTo(y)
		})

		It("a func is equal to the same func", func() {
			x := func() int { return 42 }
			y := x

			AssertThat(x).IsEqualTo(y)
		})

		//Interface
		It("two unitialized interfaces are equal", func() {
			var x interface{}
			var y interface{}

			AssertThat(x).IsEqualTo(y)
		})

		It("two interfaces assigned different values are not equal", func() {
			var x interface{}
			var y interface{}

			x = 1
			y = 2

			AssertThat(x).IsNotEqualTo(y)
		})

		//Map
		It("two equivalent maps are equal", func() {
			x := make(map[string]int)
			y := make(map[string]int)

			x["a"] = 1
			y["a"] = 1

			AssertThat(x).IsEqualTo(y)
		})

		It("two different maps are not equal", func() {
			x := make(map[string]int)
			y := make(map[string]int)

			x["a"] = 1
			y["a"] = 2

			AssertThat(x).IsNotEqualTo(y)
		})

		//Array
		It("two different equivalent arrays are equal", func() {
			x := make([]int, 10)
			x[0] = 42
			y := make([]int, 10)
			y[0] = 42

			AssertThat(x).IsEqualTo(y)
		})

		It("two different arrays are not equal", func() {
			x := make([]int, 10)
			x[0] = 1
			y := make([]int, 10)
			y[0] = 2

			AssertThat(x).IsNotEqualTo(y)
		})		

		//Slice
		It("two different equivalent slices are equal", func() {
			x := make([]int, 10)
			x[0] = 42
			y := x[:5]

			q := make([]int, 10)
			q[0] = 42
			r := q[:5]

			AssertThat(y).IsEqualTo(r)
		})

		It("two different slices are not equal", func() {
			x := make([]int, 10)
			x[0] = 1
			y := x[:5]

			q := make([]int, 10)
			q[0] = 2
			r := q[:5]

			AssertThat(y).IsNotEqualTo(r)
		})

		It("two different equivalent slices of different lengths are not equal", func() {
			x := make([]int, 10)
			x[0] = 42
			y := make([]int, 11)
			y[0] = 42

			AssertThat(x).IsNotEqualTo(y)	
		})

		//Ptr
		It("two references to the same struct are equal", func() {
			o := TestObj{i: 1, s: "abc"}
			x := &o
			y := &o

			AssertThat(x).IsEqualTo(y)
		})

		It("two different references to two different equivalent structs are equal", func() {
			x := &TestObj{i: 1, s: "abc"}
			y := &TestObj{i: 1, s: "abc"}

			AssertThat(x).IsEqualTo(y)
		})

		It("two different references to two different equivalent structs are equal when dereferenced", func() {
			x := &TestObj{i: 1, s: "abc"}
			y := &TestObj{i: 1, s: "abc"}

			AssertThat(*x).IsEqualTo(*y)
		})

		It("two different references to two different structs are not equal", func() {
			x := &TestObj{i: 1}
			y := &TestObj{i: 2}

			AssertThat(x).IsNotEqualTo(y)
		})		
	})

	Describe("When using Is and IsNot", func() {

		//Nil
		It("nil is the same object as nil", func() {
			AssertThat(nil).Is(nil)
		})

		It("a variable assigned nil is the same object as itself", func() {
			var x interface{}
			x = nil

			AssertThat(x).Is(x)
		})

		It("two different variables assigned nil are the same object", func() {
			var x interface{}
			var y interface{}

			x = nil
			y = nil

			AssertThat(x).Is(y)
		})

		It("nil is not a value", func() {
			AssertThat(nil).IsNot(42)
		})

		//Bool
		It("a bool value is not another equivalent bool value", func() {
			AssertThat(true).IsNot(true)
		})

		//Number
		It("a number is not the same object as itself", func() {
			AssertThat(42).IsNot(42)
		})

		It("two variables assigned the same number value are not the same object", func() {
			x := 42
			y := 42

			AssertThat(x).IsNot(y)
		})

		It("two copies of the same number are not the same object", func() {
			x := 42
			y := x

			AssertThat(x).IsNot(y)
		})

		//String
		It("a string is not the same object as itself", func() {
			AssertThat("abc").IsNot("abc")
		})

		It("two variables assigned the same string value are not the same object", func() {
			x := "abc"
			y := "abc"

			AssertThat(x).IsNot(y)
		})

		It("two copies of the same string are not the same object", func() {
			x := "abc"
			y := x

			AssertThat(x).IsNot(y)
		})

		//Struct
		It("two different structs are not the same object", func() {
			x := TestObj{i: 1, s: "abc"}
			y := TestObj{i: 2, s: "xyz"}

			AssertThat(x).IsNot(y)
		})
		
		It("two copies of the same struct are not the same object", func() {
			x := TestObj{i: 1, s: "abc"}
			y := x

			AssertThat(x).IsNot(y)
		})

		//Channel
		It("two references to the same channel are the same object", func() {
			x := make(chan int)
			y := x

			AssertThat(x).Is(y)
		})

		It("two references to two different channels are not the same object", func() {
			x := make(chan int)
			y := make(chan int)

			AssertThat(x).IsNot(y)
		})

		//Func
		It("two references to the same func are the same object", func() {
			x := func() string { 
				return "hello" 
			}

			y := x

			AssertThat(x).Is(y)
		})

		It("two references to two different equivalent funcs are not the same object", func() {
			x := func() string { 
				return "hello" 
			}

			y := func() string { 
				return "hello" 
			}

			AssertThat(x).IsNot(y)
		})

		//Interface
		// It("two copies of the same interface are not the same object", func() {
		// 	var x interface{}
		// 	x = 5
		// 	var y = x

		// 	AssertThat(x).IsNot(y)
		// })

		// It("two uninitialized interfaces are not the same object", func() {
		// 	var x interface{}
		// 	var y interface{}

		// 	AssertThat(x).IsNot(y)
		// })

		It("an interface is the same object as itself", func() {
			var x interface{}

			AssertThat(x).Is(x)
		})

		//Map
		It("two different equivalent maps are not the same object", func() {
			x := make(map[string]int)
			y := make(map[string]int)

			x["a"] = 1
			x["a"] = 1

			AssertThat(x).IsNot(y)
		})

		It("two copies of the same map are not the same object", func() {
			x := make(map[string]int)
			x["a"] = 1
			y := x

			AssertThat(x).Is(y)
		})

		//Array
		It("two different equivalent arrays are not the same object", func() {
			x := make([]int, 10)
			x[0] = 42
			y := make([]int, 10)
			y[0] = 42

			AssertThat(x).IsNot(y)
		})

		It("two copies of the same array are not the same object", func() {
			x := make([]int, 10)
			x[0] = 42
			y := x

			AssertThat(x).IsNot(y)
		})

		//Slice
		It("two different equivalent slices are not the same object", func() {
			x := make([]int, 10)
			x[0] = 42
			y := x[:5]

			q := make([]int, 10)
			q[0] = 42
			r := q[:5]

			AssertThat(y).IsNot(r)
		})

		It("two copies of the same slice are not the same object", func() {
			x := make([]int, 10)
			x[0] = 42
			y := x[:5]
			z := y

			AssertThat(y).IsNot(z)
		})

		//Ptr
		It("two pointers to the same object are the same object", func() {
			x := &TestObj{i: 1, s: "abc"}
			y := x

			AssertThat(x).Is(y)
		})

		It("two pointers to two copies of the same object are not the same object", func() {
			x := &TestObj{i: 1, s: "abc"}
			y := &TestObj{i: 1, s: "abc"}

			AssertThat(x).IsNot(y)
		})

		It("two pointers to different objects are not the same object", func() {
			x := &TestObj{i: 1, s: "abc"}
			y := &TestObj{i: 2, s: "xyz"}

			AssertThat(x).IsNot(y)			
		})
	})

	Describe("Assert that using Panics and DoesNotPanic", func() {

		It("panicing function panics", func() {
			AssertThat(func() {
				panic("oops")
			}).Panics()
		})

		It("non panicing function does not panic", func() {
			AssertThat(func() {
				
			}).DoesNotPanic()
		})

		It("cannot assert that non function panics", func() {
			AssertThat(func() {
				AssertThat(true).Panics()
			}).Panics()
		})

		It("cannot assert that non function does not panic", func() {
			AssertThat(func() {
				AssertThat(true).DoesNotPanic()
			}).Panics()
		})
	})
}
