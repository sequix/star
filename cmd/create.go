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
	"archive/tar"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/sequix/star/pkg/fs"
	"github.com/sequix/star/pkg/star"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use: "create <xxx.star> <files|yyy.tar>",
	Aliases: []string{"c"},
	Short: "Create a star file from files or a regular tar file.",
	Long: `Create a star file from files or a regular tar file.`,
	Run: func(cmd *cobra.Command, args []string) {
		createRun(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().BoolP("force", "f", false, "Overwrite existing file")
}

func createRun(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		cmd.Help()
		return
	}

	var (
		sfn = args[0]
		fsr fs.Reader
		err error
	)

	if len(args) == 2 && filepath.Ext(args[1]) == ".tar" {
		tf, err := os.Open(args[1])
		if err != nil {
			log.Fatalf("opening tar file %q, %s", args[1], err)
			return
		}
		fsr = fs.NewTarReader(tar.NewReader(tf))
	} else {
		fsr, err = fs.NewLocalReader(args[1:]...)
		if err != nil {
			log.Fatalf("newing local reader, %s", err)
		}
	}

	force, err := cmd.Flags().GetBool("force")
	if err != nil {
		log.Fatalf("getting flag --force, %s", err)
	}

	flag := os.O_CREATE | os.O_WRONLY
	if force {
		flag |= os.O_TRUNC
	} else {
		flag |= os.O_EXCL
	}

	sf, err := os.OpenFile(sfn, flag, 0644)
	if err != nil {
		log.Fatalf("opening star file %q, %s", sfn, err)
	}

	if err := star.WriteTo(sf, fsr); err != nil {
		log.Fatalf("creating star file %q, %s", sfn, err)
	}
}


