# gotest 
[![Build Status](https://travis-ci.org/claassen/gotest.svg?branch=master)](https://travis-ci.org/claassen/gotest)

**gotest** is a BDD testing framework for Go in the style of the Jasmine BDD testing framework for Javascript.

## Installation

```
go get github.com/claassen/gotest
```

## Writing Tests


Test file names follow the standard Go practice of matching the pattern *_test.go

Test files differ from the requirement for the `go test` tool in that test functions need not follow the Test* pattern and should not accept any arguments. Test functions must be publicly visible (start with a capital letter).

A test file can have multiple functions but typically you would only have one function and instead use `Describe` and `It` blocks to separate tests.

### Simple Example

```go
package mypackage

import(
	. "claassen/gotest/testing"
)

func Test() {
	Describe("A trivial testing example", func() {
		It("should pass", func() {
			//If this block does not panic or make an assertion that fails then the test will pass
		})
		
		It("should fail", func() {
			//A panic will cause the test to fail
			panic("something when wrong")
		})
	})
	
	It("Doesn't need to be nested in a Describe block", func() {
		//Wrapping It blocks in Describe blocks is optional
		
		//It blocks may not contain Describe or It blocks
		//Describe("this will not work", func() {})
		//It("this will not work", func() { })
	})
}
```

### Nested Describes, BeforeEach and AfterEach

```go
package mypackage

import(
	. "claassen/gotest/testing"
)

func Test() {
	Describe("Some initial behaviour", func() {
		BeforeEach(func() {
			//Will be called before each test in the current as well as any child Describe blocks
			//BeforeEach blocks in child Describe blocks will be called after any parent BeforeEach blocks
		})
		
		AfterEach(func() {
			//Will be called after each test in the current and any child Describe blocks
			//AfterEach blocks in child Describe blocks will be called before any parent AfterEach blocks
		})
		
		It("does something for the initial behaviour", func() {
			//Asserts here
		})
		
		Describe("some additional behaviour", func() {
			BeforeEach(func() {
				
			})
			
			It("does something for the initial behaviour plus the additional behaviour", func() {
				//Asserts here
			})
		})
	})
}

```

### Assertions

**gotest** includes a fluent assertions library. The framework depends on assertions to simply panic in order to indicate assertion failure in case you want to try to plugin in a different assertions library. 

```go
package mypackage

import(
	. "claassen/gotest/testing"
	. "claassen/gotest/testing/assertions"
)

func Test() {
	Describe("A trival assertion test example", func() {
		It("should pass", func() {
			AssertThat(true).IsEqualTo(true)
		})
		
		It("should fail", func() {
			Fail("Something bad")
		})
	})
}
```

## Running Tests
Run the `gotest` program providing the package name of the package you wish to test:

```shell
gotest my/package
```

This will find and run tests in the package specified as well tests in any packages which are sub-directories of the specified package.

