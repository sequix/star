/*
Copyright © 2020 sequix <sequix@163.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/sequix/star/pkg/fs"
	"github.com/sequix/star/pkg/star"
)

// extractCmd represents the extract command
var extractCmd = &cobra.Command{
	Use:   "extract <xxx.star>",
	Aliases: []string{"x"},
	Short: "Extract a star file.",
	Long: `Extract a star file.`,
	Run: func(cmd *cobra.Command, args []string) {
		extractRun(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(extractCmd)
}

func extractRun(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		cmd.Help()
		return
	}

	sfr, err := os.Open(args[0])
	if err != nil {
		fmt.Printf("opening file %q: %s\n", args[0], err)
		return
	}

	sr, err := star.NewReader(sfr)
	if err != nil {
		fmt.Printf("newing star reader: %s\n", err)
		return
	}

	for _, fi := range sr.ListFiles() {
		fr, err := sr.ReaderFor(fi.Name)
		if err != nil {
			fmt.Printf("selecting file read %q: %s\n", fi.Name, err)
			return
		}
		if err := write(fi, fr); err != nil {
			fmt.Printf("writing %q: %s\n", fi.Name, err)
			return
		}
	}
}

func write(fi *star.Info, data io.Reader) error {
	switch fi.Mode & os.ModeType {
	case os.ModeDir:
		return fs.MkdirAll(fi.FileInfo)
	case os.ModeSymlink:
		return fs.Symlink(fi.FileInfo)
	case os.ModeDevice | os.ModeCharDevice, os.ModeDevice:
		return fs.Mknod(fi.FileInfo)
	}

	log.Println("file", fi.Name)

	dir := filepath.Dir(fi.Name)
	if err := os.MkdirAll(dir, 0655); err != nil {
		return fmt.Errorf("mkdirall %q, %s", dir, err)
	}

	fw, err := os.OpenFile(fi.Name, os.O_CREATE|os.O_WRONLY|os.O_EXCL, fi.Mode)
	if err != nil {
		return fmt.Errorf("open file %q, %s", fi.Name, err)
	}
	if err := fs.Chall(fi.FileInfo); err != nil {
		return fmt.Errorf("chall file %q, %s", fi.Name, err)
	}

	n, err := io.Copy(fw, data)
	if err != nil {
		return fmt.Errorf("copy to file %q, copied %d byte, err %s", fi.Name, n, err)
	}
	return nil
}
