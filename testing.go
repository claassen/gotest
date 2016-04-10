package testing

import (
	"fmt"
	"github.com/fatih/color"
	"os"
	"strings"
)

type testContext struct {
	currentBlock       *block
	topLevelBlocks     []*block
	currentRunningTest string
	passed             int
	failed             int
}

type block struct {
	blockType   string
	description string
	parent      *block
	children    []*block
	beforeEachs []func()
	afterEachs  []func()
	body        func()
}

var t = testContext{currentBlock: nil}

func (t *testContext) addBlock(block *block) {
	if t.currentBlock == nil {
		t.topLevelBlocks = append(t.topLevelBlocks, block)
	} else {
		//Run parent BeforeEachs before child BeforeEachs
		block.beforeEachs = append(t.currentBlock.beforeEachs, block.beforeEachs...)
		//Run child AfterEachs before parent AfterEachs
		block.afterEachs = append(block.afterEachs, t.currentBlock.afterEachs...)
		t.currentBlock.children = append(t.currentBlock.children, block)
		block.parent = t.currentBlock
	}
}

func Describe(desc string, processChildBlocks func()) {
	b := block{blockType: "describe", description: desc, parent: t.currentBlock}

	t.addBlock(&b)
	t.currentBlock = &b

	processChildBlocks()

	//Reset current block after processing top level block
	if b.parent == nil {
		t.currentBlock = nil
	}
}

func It(desc string, body func()) {
	b := block{blockType: "it", description: desc, parent: t.currentBlock, body: body}

	t.addBlock(&b)

	//Reset current block after processing top level block
	if b.parent == nil {
		t.currentBlock = nil
	}
}

func BeforeEach(body func()) {
	if t.currentBlock.blockType == "describe" {
		t.currentBlock.beforeEachs = append(t.currentBlock.beforeEachs, body)
	} else {
		panic("BeforeEach may only be applied inside Describe blocks")
	}
}

func AfterEach(body func()) {
	if t.currentBlock.blockType == "describe" {
		t.currentBlock.afterEachs = append(t.currentBlock.afterEachs, body)
	} else {
		panic("AfterEach may only be applied inside Describe blocks")
	}
}

func runTest(body func(), testName string) {
	defer func() {
		err := recover()

		if err != nil {
			fmt.Println(color.RedString("FAILED:"), testName)
			errStr, ok := err.(string)
			if ok {
				fmt.Println(color.RedString(errStr))
			} else {
				fmt.Println(err)
			}

			t.failed++
		}
	}()

	body()

	fmt.Println(color.GreenString("PASSED:"), testName)

	t.passed++
}

func (b block) run(testDescriptionPrefix string) {
	testName := strings.TrimSpace(testDescriptionPrefix + " " + b.description)

	if b.blockType == "describe" {
		for _, childBlock := range b.children {

			if childBlock.blockType == "it" {
				for _, before := range b.beforeEachs {
					before()
				}
			}

			childBlock.run(testName)

			if childBlock.blockType == "it" {
				for _, after := range b.afterEachs {
					after()
				}
			}
		}
	} else {
		runTest(b.body, testName)
	}
}

func RunTests() {
	fmt.Println("Running tests...")

	for _, b := range t.topLevelBlocks {
		b.run("")
	}

	fmt.Println("-----------")

	if t.failed == 0 {
		fmt.Println("All", t.passed, "tests", color.GreenString("PASSED"))
		os.Exit(0)
	} else {
		fmt.Println(t.passed, "tests", color.GreenString("PASSED"))
		fmt.Println(t.failed, "tests", color.RedString("FAILED"))
		os.Exit(1)
	}
}
