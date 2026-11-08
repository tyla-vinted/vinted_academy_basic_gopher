// Command steps unlocks the workshop steps one at a time.
//
//	go run ./internal/steps check        test the steps you have
//	go run ./internal/steps next         if they pass, unlock the next step
//	go run ./internal/steps next -force  unlock the next step without the tests
//	go run ./internal/steps hint         a hint for the first failing test
//	go run ./internal/steps final        pass only when all steps are done
//	go run ./internal/steps goto <n>     unlock up to step n, no tests (for the host)
//	go run ./internal/steps gen          rewrite avatar/locked.go
package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const (
	avatarDir = "avatar"
	stepsDir  = ".steps"
	lockedGo  = "avatar/locked.go"
	lastStep  = 9
)

var titles = map[int]string{
	1: "functions, single and multiple return values",
	2: "flow control and enums: walking",
	3: "for loops: facing left and right",
	4: "if, strings and maps: chat and emotes",
	5: "structs and methods: high-fives",
	6: "pointers: your status and emote cooldown",
	7: "generics: the same code for any type",
	8: "errors: why a high-five failed",
	9: "testing: greetings",
}

var emptyAfterBrace = regexp.MustCompile(`\{\n\s*\n`)

func main() {
	if len(os.Args) < 2 {
		usage()
	}
	var err error
	switch os.Args[1] {
	case "check":
		err = check()
	case "next":
		err = next(len(os.Args) > 2 && os.Args[2] == "-force")
	case "hint":
		err = hint()
	case "final":
		err = final()
	case "goto":
		if len(os.Args) < 3 {
			usage()
		}
		n, convErr := strconv.Atoi(os.Args[2])
		if convErr != nil || n < 1 || n > lastStep {
			usage()
		}
		err = unlockTo(n)
	case "gen":
		err = gen(current())
	default:
		usage()
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: steps check | next [-force] | hint | final | goto <1-9> | gen")
	os.Exit(2)
}

// current returns the highest step that has files in avatar/.
func current() int {
	n := 1
	for k := 1; k <= lastStep; k++ {
		if m, _ := filepath.Glob(filepath.Join(avatarDir, fmt.Sprintf("step%d_*.go", k))); len(m) > 0 {
			n = k
		}
	}
	return n
}

func runTests() bool {
	cmd := exec.Command("go", "test", "./"+avatarDir)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	return cmd.Run() == nil
}

func check() error {
	n := current()
	fmt.Printf("Step %d of %d: %s\n\n", n, lastStep, titles[n])
	if !runTests() {
		return fmt.Errorf("\nStep %d is not done yet. Fix the failing tests above, then run: make check", n)
	}
	if n == lastStep {
		fmt.Println("\nAll steps are done. Join the lobby: make final SERVER=<url>")
	} else {
		fmt.Printf("\nStep %d is done. Unlock step %d: make next\n", n, n+1)
	}
	return nil
}

func next(force bool) error {
	n := current()
	if n == lastStep {
		fmt.Println("You have all the steps already.")
		return nil
	}
	if !force {
		fmt.Printf("Checking step %d: %s\n\n", n, titles[n])
		if !runTests() {
			return fmt.Errorf("\nFinish step %d first: make the failing tests above pass", n)
		}
	}
	files, err := unlock(n + 1)
	if err != nil {
		return err
	}
	if err := gen(n + 1); err != nil {
		return err
	}
	fmt.Printf("\nStep %d unlocked: %s\n", n+1, titles[n+1])
	for _, f := range files {
		fmt.Println("  " + f)
	}
	fmt.Println("\nWhen you are done: make check")
	return nil
}

func final() error {
	if n := current(); n < lastStep {
		return fmt.Errorf("You are on step %d of %d. Finish all the steps first: make check, then make next", n, lastStep)
	}
	if !runTests() {
		return fmt.Errorf("\nYour tests are your ticket to the lobby. Make them pass first")
	}
	return nil
}

func unlockTo(target int) error {
	for k := current() + 1; k <= target; k++ {
		if _, err := unlock(k); err != nil {
			return err
		}
	}
	return gen(max(target, current()))
}

// unlock copies the files of step k into avatar/ and returns their paths.
func unlock(k int) ([]string, error) {
	src := filepath.Join(stepsDir, fmt.Sprintf("step%d", k))
	entries, err := os.ReadDir(src)
	if err != nil {
		return nil, fmt.Errorf("step %d: %w", k, err)
	}
	var out []string
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".go") {
			continue // hints.txt stays in .steps
		}
		data, err := os.ReadFile(filepath.Join(src, e.Name()))
		if err != nil {
			return nil, err
		}
		dst := filepath.Join(avatarDir, e.Name())
		if _, err := os.Stat(dst); err == nil {
			return nil, fmt.Errorf("%s exists already; not overwriting it", dst)
		}
		if err := os.WriteFile(dst, data, 0o644); err != nil {
			return nil, err
		}
		out = append(out, dst)
	}
	return out, nil
}

// gen writes avatar/locked.go: the code of every step after step n, without
// comments, so that the arena compiles before those steps are unlocked.
func gen(n int) error {
	if n >= lastStep {
		if err := os.Remove(lockedGo); err != nil && !os.IsNotExist(err) {
			return err
		}
		return nil
	}

	fset := token.NewFileSet()
	imports := map[string]bool{}
	var decls []ast.Decl
	for k := n + 1; k <= lastStep; k++ {
		paths, _ := filepath.Glob(filepath.Join(stepsDir, fmt.Sprintf("step%d", k), "*.go"))
		sort.Strings(paths)
		for _, p := range paths {
			if strings.HasSuffix(p, "_test.go") {
				continue
			}
			f, err := parser.ParseFile(fset, p, nil, parser.SkipObjectResolution)
			if err != nil {
				return err
			}
			for _, d := range f.Decls {
				if g, ok := d.(*ast.GenDecl); ok && g.Tok == token.IMPORT {
					for _, s := range g.Specs {
						imports[s.(*ast.ImportSpec).Path.Value] = true
					}
					continue
				}
				decls = append(decls, d)
			}
		}
	}

	var body bytes.Buffer
	for _, d := range decls {
		if err := printer.Fprint(&body, fset, d); err != nil {
			return err
		}
		body.WriteString("\n\n")
	}

	var src bytes.Buffer
	src.WriteString("// Code generated by internal/steps. DO NOT EDIT.\n\n")
	src.WriteString("// The steps you have not unlocked yet, so that the arena compiles.\n")
	src.WriteString("// make next replaces them with the real files, one step at a time.\n\n")
	src.WriteString("package avatar\n\n")
	if len(imports) > 0 {
		paths := make([]string, 0, len(imports))
		for p := range imports {
			paths = append(paths, p)
		}
		sort.Strings(paths)
		src.WriteString("import (\n")
		for _, p := range paths {
			src.WriteString("\t" + p + "\n")
		}
		src.WriteString(")\n\n")
	}
	src.Write(emptyAfterBrace.ReplaceAll(body.Bytes(), []byte("{\n")))

	out, err := format.Source(src.Bytes())
	if err != nil {
		return fmt.Errorf("format %s: %w", lockedGo, err)
	}
	return os.WriteFile(lockedGo, out, 0o644)
}

const hintsState = ".hints"

var failLine = regexp.MustCompile(`--- FAIL: TestStep(\d)_([A-Za-z0-9]+)`)

// hint runs the tests and prints the next hint for the first failing test.
// Each call for the same test shows the next hint; .hints remembers where
// each learner is.
func hint() error {
	out, _ := exec.Command("go", "test", "./"+avatarDir).CombinedOutput()
	if strings.Contains(string(out), "[build failed]") || strings.Contains(string(out), "[setup failed]") {
		return fmt.Errorf("Your code does not compile yet. Run make check to see the error, then fix that first")
	}
	m := failLine.FindStringSubmatch(string(out))
	if m == nil {
		n := current()
		if n == lastStep {
			fmt.Println("All tests pass. Join the lobby: make final SERVER=<url>")
		} else {
			fmt.Printf("All tests pass. Unlock step %d: make next\n", n+1)
		}
		return nil
	}
	step, name := m[1], m[2]
	hints := readHints(filepath.Join(stepsDir, "step"+step, "hints.txt"), name)
	if len(hints) == 0 {
		fmt.Printf("TestStep%s_%s fails, and there is no hint for it. Read the test, or ask the host.\n", step, name)
		return nil
	}

	key := "step" + step + "/" + name
	seen := readState()
	i := min(seen[key], len(hints)-1)
	fmt.Printf("Hint %d of %d for %s (step %s):\n\n  %s\n\n", i+1, len(hints), name, step, hints[i])
	if i+1 < len(hints) {
		fmt.Println("Still stuck? Run make hint again for the next hint.")
	} else {
		fmt.Println("That was the last hint. If you are still stuck, ask the host.")
	}
	seen[key] = i + 1
	return writeState(seen)
}

// readHints returns the hints under "## <name>" in a hints.txt file.
func readHints(path, name string) []string {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var hints []string
	in := false
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "## ") {
			in = strings.TrimSpace(strings.TrimPrefix(line, "## ")) == name
			continue
		}
		if in && line != "" {
			hints = append(hints, line)
		}
	}
	return hints
}

func readState() map[string]int {
	seen := map[string]int{}
	data, _ := os.ReadFile(hintsState)
	for _, line := range strings.Split(string(data), "\n") {
		key, n, ok := strings.Cut(line, " ")
		if v, err := strconv.Atoi(n); ok && err == nil {
			seen[key] = v
		}
	}
	return seen
}

func writeState(seen map[string]int) error {
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var b strings.Builder
	for _, k := range keys {
		fmt.Fprintf(&b, "%s %d\n", k, seen[k])
	}
	return os.WriteFile(hintsState, []byte(b.String()), 0o644)
}
