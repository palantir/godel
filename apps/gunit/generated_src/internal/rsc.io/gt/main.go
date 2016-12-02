// Copyright 2015 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Gt is a wrapper for ``go test'' that caches test results.
//
// The usage of ``gt'' is nearly identical to that of ``go test.''
//
//	gt [-f] [-l] [arguments for "go test"]
//
// The difference between ``gt'' and ``go test'' is that when testing
// a list of packages, if a package and its dependencies have not changed
// since the last run, ``gt'' reuses the previous result.
//
// The -f flag causes gt to treat all test results as uncached, as does the
// use of any ``go test'' flag other than -short and -v.
//
// The -l flag causes gt to list the uncached tests it would run.
//
// Cached test results are saved in $CACHE/go-test-cache if $CACHE is set,
// or else $HOME/Library/Caches/go-test-cache on OS X
// and $HOME/.cache/go-test-cache on other systems.
// It is always safe to delete these directories if they become too large.
//
// Gt is an experiment in what it would mean and how well it would work
// to cache test results. If the experiment proves successful, the functionality
// may move into the standard go command.
//
// Examples
//
// Run (and cache) the strings test:
//
//	$ gt strings
//	ok  	strings	0.436s
//	$
//
// List tests in str... without cached results:
//
//	$ gt -l str...
//	strconv
//	$
//
// Run str... tests:
//
//	$ gt str...
//	ok  	strconv	1.548s
//	ok  	strings	0.436s (cached)
//	$
//
// Force rerun of both:
//
//	$ gt -f str...
//	ok  	strconv	1.795s
//	ok  	strings	0.629s
//	$
//
package amalgomated

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
)

func usage() {
	fmt.Fprint(os.Stderr, "usage: gt [arguments for \"go test\"]\n")
	os.Exit(2)
}

var (
	flagV		bool
	flagShort	bool
	flagRace	bool
	flagList	bool
	flagForce	bool
	flagTiming	bool
	failed		bool
	cacheDir	string
	start		= time.Now()
)

func AmalgomatedMain() {
	log.SetFlags(0)
	log.SetPrefix("gt: ")

	opts, pkgs := parseFlags()
	if len(pkgs) == 0 {
		if flagList {
			log.Fatal("cannot use -l without package list or with testing flags other than -v and -short")
		}
		cmd := exec.Command("go", append([]string{"test"}, os.Args[1:]...)...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Fatalf("go test: %v", err)
		}
		return
	}

	if flagTiming {
		log.Printf("%.2fs go list", time.Since(start).Seconds())
	}

	// Expand pkg list.
	out, err := exec.Command("go", append([]string{"list"}, pkgs...)...).CombinedOutput()
	if err != nil {
		log.Fatalf("go list: %v", err)
	}
	pkgs = strings.Fields(string(out))

	if flagTiming {
		log.Printf("%.2fs go list -json", time.Since(start).Seconds())
	}

	// Build list of all dependencies.
	readPkgs(pkgs)

	first := true
	next := pkgs
	for {
		var deps []string
		for _, path := range next {
			p := pkgInfo[path]
			if p.Incomplete {
				log.Fatalf("go list: errors loading packages")
			}
			for _, dep := range p.Deps {
				if _, ok := pkgInfo[dep]; !ok {
					pkgInfo[dep] = nil
					deps = append(deps, dep)
				}
			}
			if first {
				for _, dep := range p.TestImports {
					if _, ok := pkgInfo[dep]; !ok {
						pkgInfo[dep] = nil
						deps = append(deps, dep)
					}
				}
				for _, dep := range p.XTestImports {
					if _, ok := pkgInfo[dep]; !ok {
						pkgInfo[dep] = nil
						deps = append(deps, dep)
					}
				}
			}
		}
		if len(deps) == 0 {
			break
		}

		if flagTiming {
			log.Printf("%.2fs go list -json", time.Since(start).Seconds())
		}
		readPkgs(deps)
		next = deps
		first = false
	}

	if env := os.Getenv("CACHE"); env != "" {
		cacheDir = fmt.Sprintf("%s/go-test-cache", env)
	} else if runtime.GOOS == "darwin" {
		cacheDir = fmt.Sprintf("%s/Library/Caches/go-test-cache", os.Getenv("HOME"))
	} else {
		cacheDir = fmt.Sprintf("%s/.cache/go-test-cache", os.Getenv("HOME"))
	}

	if flagTiming {
		log.Printf("%.2fs compute hashes", time.Since(start).Seconds())
	}

	computeStale(pkgs)

	var toRun []string
	for _, pkg := range pkgs {
		if !haveTestResult(pkg) {
			toRun = append(toRun, pkg)
		}
	}

	if flagTiming {
		log.Printf("%.2fs ready to run", time.Since(start).Seconds())
	}

	if flagList {
		for _, pkg := range toRun {
			fmt.Printf("%s\n", pkg)
		}
		return
	}

	var cmd *exec.Cmd
	pr, pw := io.Pipe()
	r := bufio.NewReader(pr)
	if len(toRun) > 0 {
		if err := os.MkdirAll(cacheDir, 0700); err != nil {
			log.Fatal(err)
		}

		args := []string{"test"}
		args = append(args, opts...)
		args = append(args, toRun...)
		cmd = exec.Command("go", args...)
		cmd.Stdout = pw
		cmd.Stderr = pw
		if err := cmd.Start(); err != nil {
			log.Fatalf("go test: %v", err)
		}
	}

	var cmdErr error
	done := make(chan bool)
	go func() {
		if cmd != nil {
			cmdErr = cmd.Wait()
		}
		pw.Close()
		done <- true
	}()

	for _, pkg := range pkgs {
		if len(toRun) > 0 && toRun[0] == pkg {
			readTestResult(r, pkg)
			toRun = toRun[1:]
		} else {
			showTestResult(pkg)
		}
	}

	io.Copy(os.Stdout, r)

	<-done
	if cmdErr != nil && !failed {
		log.Fatalf("go test: %v", cmdErr)
	}

	if flagTiming {
		log.Printf("%.2fs done", time.Since(start).Seconds())
	}

	if failed {
		os.Exit(1)
	}
}

var (
	pkgInfo		= map[string]*Package{}
	outOfSync	bool
)

type Package struct {
	Dir		string
	ImportPath	string
	Standard	bool
	Goroot		bool
	Stale		bool
	GoFiles		[]string
	CgoFiles	[]string
	CFiles		[]string
	CXXFiles	[]string
	MFiles		[]string
	HFiles		[]string
	SFiles		[]string
	SwigFiles	[]string
	SwigCXXFiles	[]string
	SysoFiles	[]string
	Imports		[]string
	Deps		[]string
	Incomplete	bool
	TestGoFiles	[]string
	TestImports	[]string
	XTestGoFiles	[]string
	XTestImports	[]string

	testHash	string
	pkgHash		string
}

func readPkgs(pkgs []string) {
	out, err := exec.Command("go", append([]string{"list", "-json"}, pkgs...)...).CombinedOutput()
	if err != nil {
		log.Fatalf("go list: %v\n%s", err, out)
	}

	dec := json.NewDecoder(bytes.NewReader(out))
	for {
		var p Package
		if err := dec.Decode(&p); err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("reading go list output: %v", err)
		}
		pkgInfo[p.ImportPath] = &p
	}
}

func computeStale(pkgs []string) {
	for _, pkg := range pkgs {
		computeTestHash(pkgInfo[pkg])
	}
}

func computeTestHash(p *Package) {
	if p.testHash != "" {
		return
	}
	p.testHash = "cycle"
	computePkgHash(p)
	h := sha1.New()
	fmt.Fprintf(h, "test\n")
	if flagRace {
		fmt.Fprintf(h, "-race\n")
	}
	if flagShort {
		fmt.Fprintf(h, "-short\n")
	}
	if flagV {
		fmt.Fprintf(h, "-v\n")
	}
	fmt.Fprintf(h, "pkg %s\n", p.pkgHash)
	for _, imp := range p.TestImports {
		p1 := pkgInfo[imp]
		computePkgHash(p1)
		fmt.Fprintf(h, "testimport %s\n", p1.pkgHash)
	}
	for _, imp := range p.XTestImports {
		p1 := pkgInfo[imp]
		computePkgHash(p1)
		fmt.Fprintf(h, "xtestimport %s\n", p1.pkgHash)
	}
	hashFiles(h, p.Dir, p.TestGoFiles)
	hashFiles(h, p.Dir, p.XTestGoFiles)
	p.testHash = fmt.Sprintf("%x", h.Sum(nil))
}

func computePkgHash(p *Package) {
	if p.pkgHash != "" {
		return
	}
	p.pkgHash = "cycle"
	h := sha1.New()
	fmt.Fprintf(h, "pkg\n")
	for _, imp := range p.Deps {
		p1 := pkgInfo[imp]
		if p1 == nil {
			log.Fatalf("lost package: %v for %v", imp, p.ImportPath)
		}
		computePkgHash(p1)
		fmt.Fprintf(h, "import %s\n", p1.pkgHash)
	}
	hashFiles(h, p.Dir, p.GoFiles)
	hashFiles(h, p.Dir, p.CgoFiles)
	hashFiles(h, p.Dir, p.CFiles)
	hashFiles(h, p.Dir, p.CXXFiles)
	hashFiles(h, p.Dir, p.MFiles)
	hashFiles(h, p.Dir, p.HFiles)
	hashFiles(h, p.Dir, p.SFiles)
	hashFiles(h, p.Dir, p.SwigFiles)
	hashFiles(h, p.Dir, p.SwigCXXFiles)
	hashFiles(h, p.Dir, p.SysoFiles)

	p.pkgHash = fmt.Sprintf("%x", h.Sum(nil))
}

func hashFiles(h io.Writer, dir string, files []string) {
	for _, file := range files {
		f, err := os.Open(filepath.Join(dir, file))
		if err != nil {
			fmt.Fprintf(h, "%s error\n", file)
			continue
		}
		fmt.Fprintf(h, "file %s\n", file)
		n, _ := io.Copy(h, f)
		fmt.Fprintf(h, "%d bytes\n", n)
		f.Close()
	}
}

func cacheFile(p *Package) string {
	return filepath.Join(cacheDir, fmt.Sprintf("%s/%s.test", p.testHash[:3], p.testHash[3:]))
}

func haveTestResult(path string) bool {
	if flagForce {
		return false
	}
	p := pkgInfo[path]
	if p.testHash == "cycle" {
		return false
	}
	fi, err := os.Stat(cacheFile(pkgInfo[path]))
	return err == nil && fi.Mode().IsRegular()
}

var fail = []byte("FAIL")

func showTestResult(path string) {
	p := pkgInfo[path]
	if p.testHash == "cycle" {
		return
	}
	data, err := ioutil.ReadFile(cacheFile(pkgInfo[path]))
	if err != nil {
		fmt.Printf("%v\n", err)
		fmt.Printf("FAIL\t%s\t(cached)\n", path)
		return
	}
	os.Stdout.Write(data)
	data = bytes.TrimSpace(data)
	i := bytes.LastIndex(data, []byte{'\n'})
	line := data[i+1:]
	if bytes.HasPrefix(line, fail) {
		failed = true
	}
}

var endRE = regexp.MustCompile(`\A(\?|ok|FAIL) ? ? ?\t([^ \t]+)\t([0-9.]+s|\[.*\])\n\z`)

func readTestResult(r *bufio.Reader, path string) {
	var buf bytes.Buffer
	for {
		line, err := r.ReadString('\n')
		os.Stdout.WriteString(line)
		if err != nil {
			log.Fatalf("reading test output for %s: %v", path, err)
		}
		if outOfSync {
			continue
		}
		m := endRE.FindStringSubmatch(line)
		if m == nil {
			buf.WriteString(line)
			continue
		}

		if m[1] == "FAIL" {
			failed = true
		}

		fmt.Fprintf(&buf, "%s (cached)\n", strings.TrimSuffix(line, "\n"))
		file := cacheFile(pkgInfo[path])
		if err := os.MkdirAll(filepath.Dir(file), 0700); err != nil {
			log.Print(err)
		} else if err := ioutil.WriteFile(file, buf.Bytes(), 0600); err != nil {
			log.Print(err)
		}

		break
	}
}

func parseFlags() (opts, pkgs []string) {
	donePkgs := false
	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		if !strings.HasPrefix(arg, "-") {
			if donePkgs {
				// additional arguments after pkg list ended
				return nil, nil
			}
			pkgs = append(pkgs, arg)
			continue
		}
		donePkgs = len(pkgs) > 0
		if strings.HasPrefix(arg, "--") && arg != "--" && !strings.HasPrefix(arg, "---") {
			arg = arg[1:]
		}
		if arg == "-gt.timing" {
			flagTiming = true
			continue
		}
		if arg == "-v" {
			flagV = true
			opts = append(opts, arg)
			donePkgs = len(pkgs) > 0
			continue
		}
		if arg == "-race" {
			flagRace = true
			opts = append(opts, arg)
			continue
		}
		if arg == "-short" {
			flagShort = true
			opts = append(opts, arg)
			continue
		}
		if arg == "-f" {
			flagForce = true
			continue
		}
		if arg == "-l" {
			flagList = true
			continue
		}
		// unrecognized flag
		return nil, nil
	}
	return opts, pkgs
}
