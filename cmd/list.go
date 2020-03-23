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
	"os"
	"path/filepath"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/sequix/star/pkg/star"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list <xxx.star>",
	Aliases: []string{"t"},
	Short:   "List content of a star file.",
	Long:    `List content of a star file.`,
	Run: func(cmd *cobra.Command, args []string) {
		listRun(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolP("long", "l", false, "Print in ls-style")
	listCmd.Flags().BoolP("basename", "b", false, "Print basename instead of full path")
	listCmd.Flags().BoolP("human", "s", false, "Print size in human-friendly format")
	listCmd.Flags().BoolP("ctime", "c", false, "Print ctime instead of mtime")
	listCmd.Flags().BoolP("atime", "a", false, "Print atime instead of atime")
}

func listRun(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		cmd.Help()
		return
	}

	sfn := args[0]
	sf, err := os.Open(sfn)
	if err != nil {
		fmt.Printf("opening star file %q, %s\n", sfn, err)
		return
	}

	sfr, err := star.NewReader(sf)
	if err != nil {
		fmt.Printf("newing star reader, %s\n", err)
		return
	}

	flags := cmd.Flags()
	long, err := flags.GetBool("long")
	if err != nil {
		fmt.Printf("getting flag `long`: %s\n", err)
		return
	}

	basename, err := flags.GetBool("basename")
	if err != nil {
		fmt.Printf("getting flag `basename`: %s\n", err)
		return
	}

	human, err := flags.GetBool("human")
	if err != nil {
		fmt.Printf("getting flag `human`: %\n", err)
		return
	}

	ctime, err := flags.GetBool("ctime")
	if err != nil {
		fmt.Printf("getting flag `ctime`: %s\n", err)
		return
	}

	atime, err := flags.GetBool("atime")
	if err != nil {
		fmt.Printf("getting flag `atime`: %s\n", err)
		return
	}

	var time string
	switch {
	case atime:
		time = "atime"
	case ctime:
		time = "ctime"
	default:
		time = "mtime"
	}

	if !long {
		for _, fn := range sfr.ListNames() {
			if basename {
				fn = filepath.Base(fn)
			}
			fmt.Println(fn)
		}
	} else {
		for _, fi := range sfr.ListFiles() {
			printFileInfo(fi, human, basename, time)
		}
	}
}

func printFileInfo(fi *star.Info, human, basename bool, timeKind string) {
	isDevice := false
	switch fileType := fi.Mode & os.ModeType; fileType {
	case os.ModeSymlink:
		fmt.Print("l")
	case os.ModeDir:
		fmt.Print("d")
	case os.ModeDevice:
		isDevice = true
		fmt.Print("b")
	case os.ModeDevice | os.ModeCharDevice:
		isDevice = true
		fmt.Print("c")
	default:
		fmt.Print("-")
	}
	 for i, mod := 8, fi.Mode.Perm(); i >= 0; i-- {
		 switch i % 3 {
		 case 2:
		 	if ((mod >> i) & 1) != 0 {
		 		fmt.Print("r")
			} else {
				fmt.Print("-")
			}
		 case 1:
			 if ((mod >> i) & 1) != 0 {
				 fmt.Print("w")
			 } else {
				 fmt.Print("-")
			 }
		 case 0:
			 if ((mod >> i) & 1) != 0 {
				 fmt.Print("x")
			 } else {
				 fmt.Print("-")
			 }
		 }
	 }
	 fmt.Printf(" %d %d\t", fi.Uid, fi.Gid)

	 if isDevice {
	 	fmt.Printf("%d,%d\t", fi.Major, fi.Minor)
	 } else {
		 if human {
			 fmt.Printf("%s\t", humanize.IBytes(fi.Size))
		 } else {
			 fmt.Printf("%d\t", fi.Size)
		 }
	 }

	switch timeKind {
	case "mtime":
		fmt.Printf("%s ", fi.Mtime.Format("2006-01-02 15:04:05"))
	case "ctime":
		fmt.Printf("%s ", fi.Ctime.Format("2006-01-02 15:04:05"))
	case "atime":
		fmt.Printf("%s ", fi.Atime.Format("2006-01-02 15:04:05"))
	}

	if basename {
		fmt.Printf("%s", filepath.Base(fi.Name))
	} else {
		fmt.Printf("%s", fi.Name)
	}

	if len(fi.Linkname) > 0 {
		fmt.Printf(" -> %s", fi.Linkname)
	}
	fmt.Println()
}
