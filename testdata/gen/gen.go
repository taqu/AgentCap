//go:build ignore

// gen generates fixture files used by reducer tests.
package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	dir := "."
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}

	must(genLsLarge(filepath.Join(dir, "ls-large.txt")))
	must(genFindLarge(filepath.Join(dir, "find-large.txt")))
	must(genRgLarge(filepath.Join(dir, "rg-large.txt")))
	must(genCatLarge(filepath.Join(dir, "cat-large.txt")))
	must(genTreeLarge(filepath.Join(dir, "tree-large.txt")))
	must(genDuLarge(filepath.Join(dir, "du-large.txt")))
}

func must(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func genLsLarge(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintln(f, "total 2097152")
	months := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
	// 151 directories
	for i := 0; i < 151; i++ {
		month := months[i%12]
		fmt.Fprintf(f, "drwxr-xr-x  2 user group    4096 %s %2d 12:00 dir%04d\n", month, (i%28)+1, i)
	}
	// 849 regular files with varying sizes
	for i := 0; i < 849; i++ {
		month := months[i%12]
		size := int64(1024) * int64(i+1)
		fmt.Fprintf(f, "-rw-r--r--  1 user group %7d %s %2d 14:30 file%04d.go\n", size, month, (i%28)+1, i)
	}
	return nil
}

func genFindLarge(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	dirs := []string{"vendor", "src", "tests", "docs", "build", "cmd", "internal", "pkg"}
	exts := []string{".go", ".md", ".json", ".yaml", ".sh", ".txt", ".pb.go", ".sum"}

	for i := 0; i < 5000; i++ {
		d := dirs[i%len(dirs)]
		ext := exts[i%len(exts)]
		subdir := fmt.Sprintf("sub%d", i%50)
		fmt.Fprintf(f, "./%s/%s/file%04d%s\n", d, subdir, i, ext)
	}
	return nil
}

func genRgLarge(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	files := make([]string, 50)
	for i := range files {
		files[i] = fmt.Sprintf("src/module%02d/file.go", i)
	}

	lineNum := 1
	for i := 0; i < 400; i++ {
		file := files[i%len(files)]
		fmt.Fprintf(f, "%s:%d:func Process%d(ctx context.Context) error {\n", file, lineNum, i)
		lineNum += 5
	}
	return nil
}

func genCatLarge(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintln(f, "package main")
	fmt.Fprintln(f)
	fmt.Fprintln(f, "import (")
	fmt.Fprintln(f, `	"context"`)
	fmt.Fprintln(f, `	"fmt"`)
	fmt.Fprintln(f, `	"os"`)
	fmt.Fprintln(f, ")")
	fmt.Fprintln(f)

	for i := 0; i < 2000; i++ {
		fmt.Fprintf(f, "// Function%d does something important\n", i)
		fmt.Fprintf(f, "func Function%d(ctx context.Context) error {\n", i)
		fmt.Fprintln(f, `	fmt.Println("executing function")`)
		fmt.Fprintln(f, `	return nil`)
		fmt.Fprintln(f, "}")
		fmt.Fprintln(f)
	}
	return nil
}

func genTreeLarge(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintln(f, ".")
	dirs := []string{"cmd", "internal", "pkg", "tests", "vendor", "docs", "build", "scripts"}
	totalDirs := 0
	totalFiles := 0

	for di, d := range dirs {
		prefix1 := "├──"
		cont1 := "│"
		if di == len(dirs)-1 {
			prefix1 = "└──"
			cont1 = " "
		}
		fmt.Fprintf(f, "%s %s/\n", prefix1, d)
		totalDirs++
		for j := 0; j < 40; j++ {
			sub := fmt.Sprintf("sub%02d", j)
			prefix2 := "├──"
			cont2 := "│"
			if j == 39 {
				prefix2 = "└──"
				cont2 = " "
			}
			fmt.Fprintf(f, "%s   %s %s/\n", cont1, prefix2, sub)
			totalDirs++
			for k := 0; k < 5; k++ {
				prefix3 := "├──"
				cont3 := "│"
				if k == 4 {
					prefix3 = "└──"
					cont3 = " "
				}
				fmt.Fprintf(f, "%s   %s   %s file%02d.go\n", cont1, cont2, prefix3, k)
				totalFiles++
				// Add one more depth level for some entries
				for m := 0; m < 3; m++ {
					_ = cont3
					fmt.Fprintf(f, "%s   %s   %s   ├── deep%02d.go\n", cont1, cont2, cont3, m)
					totalFiles++
				}
			}
		}
	}
	fmt.Fprintf(f, "\n%d directories, %d files\n", totalDirs, totalFiles)
	return nil
}

func genDuLarge(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	dirs := []string{"vendor", "build", ".git", "testdata", "cmd", "internal", "pkg", "docs"}
	sizes := []int64{2900000, 1100000, 420000, 180000, 50000, 30000, 20000, 10000}

	for i, d := range dirs {
		// Top-level entry
		fmt.Fprintf(f, "%d\t./%s\n", sizes[i]/1024, d)
	}

	// Sub-entries to reach 200 total.
	for i := 0; i < 192; i++ {
		d := dirs[i%len(dirs)]
		sub := fmt.Sprintf("sub%03d", i)
		size := sizes[i%len(sizes)] / int64(i%10+2)
		fmt.Fprintf(f, "%d\t./%s/%s\n", size/1024, d, sub)
	}
	return nil
}
