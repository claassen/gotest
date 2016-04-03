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
	testPackages []TestPackageInfo
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

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Usage: gotest <package name>")
	}

	goPath := os.ExpandEnv("$GOPATH")
	goPathSrc := filepath.Join(goPath, "src")

	rootPackageName := os.Args[1]

	//fmt.Println("Testing package: ", rootPackageName)

	rootPackageFullPath := os.ExpandEnv(filepath.Join(goPathSrc, rootPackageName))

	//fmt.Println("Root package path: ", rootPackageFullPath)

	filepath.Walk(rootPackageFullPath, func(path string, fileInfo os.FileInfo, err error) (e error) {
		if fileInfo.IsDir() {

			if strings.HasSuffix(filepath.Base(path), "__test") || strings.HasSuffix(filepath.Base(path), "__testmain") {
				os.RemoveAll(path)
				return
			}

			//fmt.Println("Processing package directory: ", path)

			testPackageInfo := TestPackageInfo{}
			testPackageInfo.originalPackageName = filepath.Base(path)
			testPackageInfo.originalPackagePath = path
			testPackageInfo.testPackageName = filepath.Base(path) + "__test"
			testPackageInfo.testPackageFullName = strings.TrimLeft(strings.Replace(path, goPathSrc, "", 1)+"/"+testPackageInfo.testPackageName, "/")
			testPackageInfo.testPackagePath = filepath.Join(testPackageInfo.originalPackagePath, testPackageInfo.testPackageName)

			includePackage := false

			files, _ := ioutil.ReadDir(path)
			for _, f := range files {
				if !f.IsDir() && filepath.Ext(f.Name()) == ".go" {

					//fmt.Println("Processing file: ", f.Name())

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

							//fmt.Println("TEST FUNCTION:", fnDecl.Name.String())

							testPackageInfo.testFuncNames = append(testPackageInfo.testFuncNames, fnDecl.Name.String())
						}
					}
				}
			}

			if includePackage {
				context.testPackages = append(context.testPackages, testPackageInfo)
			}
		}

		return err
	})

	//fmt.Println("")

	//for _, p := range context.testPackages {
	//fmt.Println(p.originalPackageName)
	//}

	//Create test package directories, copy .go files, rename _test.go files so we can build them
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

	//Create testmain package dir and main package file for running tests
	testMainPackageDir := filepath.Join(rootPackageFullPath, "__testmain")
	if err := os.MkdirAll(testMainPackageDir, os.ModePerm); err != nil {
		fmt.Println("Error creating __testmain directory: ", err)
	}

	testMainFilePath := filepath.Join(testMainPackageDir, "testmain.go")
	testMainFile, err := os.Create(testMainFilePath)

	if err != nil {
		fmt.Println("Error creating testmain.go:", err)
	}
	defer testMainFile.Close()

	testMainWriter := bufio.NewWriter(testMainFile)

	fmt.Fprintln(testMainWriter, "package main")
	fmt.Fprintln(testMainWriter, "import(")
	fmt.Fprintln(testMainWriter, "\"claassen/gotest/testing\"")

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

	//testMainExePath := filepath.Join(testMainPackageDir, "testmain")

	//buildCmd := exec.Command("go", "build", "-o", testMainExePath, rootPackageName+"/__testmain")
	//if err := buildCmd.Run(); err != nil {
	//	fmt.Println(err)
	//}

	runCmd := exec.Command("go", "run", testMainFilePath)
	runCmd.Stdout = os.Stdout
	runCmd.Stderr = os.Stderr

	if err := runCmd.Run(); err != nil {
		fmt.Println(err)
	}

	//Cleanup
	os.RemoveAll(testMainPackageDir)

	for _, p := range context.testPackages {
		os.RemoveAll(p.testPackagePath)
	}
}
