package assert

import (
	. "github.com/claassen/gotest/testing"
)

type TestObj struct {
	i int
	s string
}

func Test() {

	Describe("Assert that", func() {

		It("true equals true", func() {
			AssertThat(true).IsEqualTo(true)
		})

		It("true does not equal false", func() {
			AssertThat(true).IsNotEqualTo(false)
		})

		It("two copies of an identical object are equal", func() {
			x := TestObj{i: 1, s: "abc"}
			y := TestObj{i: 1, s: "abc"}

			AssertThat(x).IsEqualTo(y)
		})

		It("two copies of different objects are not equal", func() {
			x := TestObj{i: 1, s: "abc"}
			y := TestObj{i: 1, s: "xyz"}

			AssertThat(x).IsNotEqualTo(y)
		})

		It("two references to the same object are equal", func() {
			o := TestObj{i: 1, s: "abc"}
			x := &o
			y := &o

			AssertThat(x).IsEqualTo(y)
		})

		It("two different references to copies of the same object are not equal", func() {
			x := &TestObj{i: 1, s: "abc"}
			y := &TestObj{i: 1, s: "abc"}

			AssertThat(x).IsNotEqualTo(y)
		})

		It("two different references to copies of the same object are equal when dereferenced", func() {
			x := &TestObj{i: 1, s: "abc"}
			y := &TestObj{i: 1, s: "abc"}

			AssertThat(*x).IsEqualTo(*y)
		})

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
