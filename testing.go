package testing

import (
	"fmt"
	"github.com/fatih/color"
	"os"
	"strings"
)

type TestContext struct {
	currentBlock       *Block
	topLevelBlocks     []*Block
	currentRunningTest string
	passed             int
	failed             int
}

type Block struct {
	blockType   string
	description string
	parent      *Block
	children    []*Block
	beforeEachs []func()
	afterEachs  []func()
	body        func()
}

var t = TestContext{currentBlock: nil}

func (t *TestContext) AddBlock(block *Block) {
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
	block := Block{blockType: "describe", description: desc, parent: t.currentBlock}

	t.AddBlock(&block)
	t.currentBlock = &block

	processChildBlocks()

	//Reset current block after processing top level block
	if block.parent == nil {
		t.currentBlock = nil
	}
}

func It(desc string, body func()) {
	block := Block{blockType: "it", description: desc, parent: t.currentBlock, body: body}

	t.AddBlock(&block)

	//Reset current block after processing top level block
	if block.parent == nil {
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

func (b Block) Run(testDescriptionPrefix string) {
	testName := strings.TrimSpace(testDescriptionPrefix + " " + b.description)

	if b.blockType == "describe" {
		for _, childBlock := range b.children {

			if childBlock.blockType == "it" {
				for _, before := range b.beforeEachs {
					before()
				}
			}

			childBlock.Run(testName)

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
		b.Run("")
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
