package tasktype

import (
	"context"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"testing"
)

func TestRunGolang(t *testing.T) {
	pythoncode := `package main
import "fmt"
func main() {
	fmt.Println("testgolang")
}
`
	cmd, codepath, err := rungolang(context.Background(), pythoncode)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.Remove(codepath)
	}()
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(output), "testgolang") {
		t.Errorf("run golang failed,want res:testgolang,but get res:%s", output)
	}
}

func TestRunShell(t *testing.T) {
	shellcode := `sleep 1
echo testshell`
	cmd, codepath, err := runshell(context.Background(), shellcode)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.Remove(codepath)
	}()
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(output), "testshell") {
		t.Errorf("run golang failed,want res:testshell,but get res:%s", output)
	}
}

func TestRunPython(t *testing.T) {
	shellcode := `sleep 1
echo testpython`
	cmd, codepath, err := runshell(context.Background(), shellcode)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.Remove(codepath)
	}()
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(output), "testpython") {
		t.Errorf("run golang failed,want res:testshell,but get res:%s", output)
	}
}

func TestImport(t *testing.T) {
	cmd := exec.Command("go", "version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}

	pattern := `[0-1]\.[0-9]{1,2}`
	re := regexp.MustCompile(pattern)
	t.Log(re.FindString(string(out)) > "1.11")
}
