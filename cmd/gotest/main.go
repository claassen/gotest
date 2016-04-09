package main

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type TestContext struct {
	rootPackageName     string
	rootPackageFullPath string
	testPackages        []TestPackageInfo
	testMainPackageDir  string
	testMainFilePath    string
}

type TestPackageInfo struct {
	originalPackageName string
	originalPackagePath string
	testPackageName     string
	testPackageFullName string
	testPackagePath     string
	goFileNames         []string
	testFuncNames       []string
}

var context = TestContext{}
var goPath = strings.Split(os.ExpandEnv("$GOPATH"), ":")[0]
var goPathSrc = filepath.Join(goPath, "src")

func processDir(path string) {

	testPackageInfo := TestPackageInfo{}
	testPackageInfo.originalPackageName = filepath.Base(path)
	testPackageInfo.originalPackagePath = path
	testPackageInfo.testPackageName = filepath.Base(path) + "__test"
	testPackageInfo.testPackageFullName = strings.TrimLeft(strings.Replace(path, goPathSrc, "", 1)+"/"+testPackageInfo.testPackageName, "/")
	testPackageInfo.testPackagePath = filepath.Join(testPackageInfo.originalPackagePath, testPackageInfo.testPackageName)

	includePackage := false

	files, err := ioutil.ReadDir(path)

	if err != nil {
		fmt.Println(err)
		return
	}

	for _, f := range files {
		if f.IsDir() {
			if f.Name() != "vendor" && f.Name() != ".git" {
				processDir(filepath.Join(path, f.Name()))
			}
		} else if filepath.Ext(f.Name()) == ".go" {
			testPackageInfo.goFileNames = append(testPackageInfo.goFileNames, f.Name())

			if strings.HasSuffix(f.Name(), "_test.go") {

				//Read Test function names
				fset := token.NewFileSet()
				f, parseErr := parser.ParseFile(fset, filepath.Join(path, f.Name()), nil, 0)

				if parseErr != nil {
					fmt.Println(parseErr)
				}

				if f.Name.Name == "main" {
					break
				}

				includePackage = true

				for _, decl := range f.Decls {
					//See if decl is a func decl
					fnDecl, ok := decl.(*ast.FuncDecl)

					if !ok {
						continue
					}

					//Skip instance function
					if fnDecl.Recv != nil {
						continue
					}

					if fnDecl.Type.Results != nil {
						fmt.Println("Skipping func:", fnDecl.Name.String(), ". Test functions cannot have a return type.")
						continue
					}

					if fnDecl.Type.Params.List != nil {
						fmt.Println("Skipping func:", fnDecl.Name.String(), ". Test functions cannot have parameters.")
						continue
					}

					testPackageInfo.testFuncNames = append(testPackageInfo.testFuncNames, fnDecl.Name.String())
				}
			}
		}
	}

	if includePackage {
		context.testPackages = append(context.testPackages, testPackageInfo)
	}
}

func createTestPackages() {

	for _, p := range context.testPackages {

		if err := os.MkdirAll(p.testPackagePath, os.ModePerm); err != nil {
			fmt.Println("Error creating temp test package path ", p.testPackagePath, ":", err)
		}

		for _, fileName := range p.goFileNames {
			originalFileFullPath := filepath.Join(p.originalPackagePath, fileName)
			copiedFileFullPath := filepath.Join(p.testPackagePath, strings.Replace(fileName, "_test.go", "_testx.go", 1))

			originalFile, err1 := os.Open(originalFileFullPath)
			if err1 != nil {
				fmt.Println("Error opening file ", originalFileFullPath, ":", err1)
			}
			defer originalFile.Close()

			copiedFile, err2 := os.Create(copiedFileFullPath)
			if err2 != nil {
				fmt.Println("Error creating file copy ", copiedFileFullPath, ":", err2)
			}
			defer copiedFile.Close()

			originalFileScanner := bufio.NewScanner(originalFile)
			copiedFileWriter := bufio.NewWriter(copiedFile)

			for originalFileScanner.Scan() {
				line := originalFileScanner.Text()

				if strings.Contains(line, "package "+p.originalPackageName) {
					line = "package " + p.testPackageName
				}

				fmt.Fprintln(copiedFileWriter, line)
			}

			copiedFileWriter.Flush()
		}
	}
}

func createTestMainPackage() {

	context.testMainPackageDir = filepath.Join(context.rootPackageFullPath, "__testmain")

	if err := os.MkdirAll(context.testMainPackageDir, os.ModePerm); err != nil {
		fmt.Println("Error creating __testmain directory: ", err)
	}

	context.testMainFilePath = filepath.Join(context.testMainPackageDir, "testmain.go")
	testMainFile, err := os.Create(context.testMainFilePath)

	if err != nil {
		fmt.Println("Error creating testmain.go:", err)
	}
	defer testMainFile.Close()

	testMainWriter := bufio.NewWriter(testMainFile)

	fmt.Fprintln(testMainWriter, "package main")
	fmt.Fprintln(testMainWriter, "import(")
	fmt.Fprintln(testMainWriter, "\"github.com/claassen/gotest\"")

	for _, p := range context.testPackages {
		fmt.Fprintln(testMainWriter, "\""+p.testPackageFullName+"\"")
	}

	fmt.Fprintln(testMainWriter, ")")
	fmt.Fprintln(testMainWriter, "func main() {")

	for _, p := range context.testPackages {
		for _, fn := range p.testFuncNames {
			fmt.Fprintln(testMainWriter, p.testPackageName+"."+fn+"()")
		}
	}

	fmt.Fprintln(testMainWriter, "testing.RunTests()")

	fmt.Fprintln(testMainWriter, "}")

	testMainWriter.Flush()
}

func runTests() {

	runCmd := exec.Command("go", "run", context.testMainFilePath)
	runCmd.Stdout = os.Stdout
	runCmd.Stderr = os.Stderr

	if err := runCmd.Run(); err != nil {
		fmt.Println(err)
	}
}

func cleanup() {

	os.RemoveAll(context.testMainPackageDir)

	for _, p := range context.testPackages {
		os.RemoveAll(p.testPackagePath)
	}
}

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Usage: gotest <package name>")
	}

	context.rootPackageName = os.Args[1]
	context.rootPackageFullPath = os.ExpandEnv(filepath.Join(goPathSrc, context.rootPackageName))

	processDir(context.rootPackageFullPath)

	createTestPackages()

	createTestMainPackage()

	runTests()

	cleanup()
}
