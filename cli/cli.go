package cli

import (
	"bufio"
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

const pragma string = "//gee:"

var (
	dir     string
	inPlace bool

	// Constants for the error handling
	errWithoutMsg string = "%sif err != nil {\n%sreturn err\n%s}"
	errWithMsg    string = "%sif err != nil {\n%sreturn fmt.Errorf(\"%s: %%w\", err)\n%s}"
)

var rootCmd = &cobra.Command{
	Use:   "gee",
	Short: "Go Error Expander (GEE) injects error handling into Go code",
	Long:  "Go Error Expander (GEE) injects error handling into Go code",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			log.Fatal("Too many arguments")
		}
		// If the directory path is provided, use it
		// instead of running on the current directory
		if len(args) == 1 {
			dir = args[0]
		} else {
			dir = "."
		}
		// Fetch all `.go` files in the directory
		var files []string
		if err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if filepath.Ext(path) == ".go" && !strings.HasSuffix(path, ".gen.go") {
				files = append(files, path)
			}
			return nil
		}); err != nil {
			log.Fatal(fmt.Errorf("failed to discover Go files: %w", err))
		}

		// For each file, read it line by line.
		// Store the line with the comment prefix
		// if it is present above the line with the
		// question mark, and use it to generate
		// the error string.
		var (
			errMsg   string
			isPragma bool
			inject   bool
		)
		for _, file := range files {
			f, err := os.Open(file)
			if err != nil {
				log.Println(fmt.Errorf("failed to open file %s: %w", file, err))
				continue
			}
			defer f.Close()

			log.Println("Processing file", file)

			// The output buffer to write to
			w := bytes.NewBuffer([]byte{})

			s := bufio.NewScanner(f)
			// Read the file line by line
			for s.Scan() {
				line := s.Text()
				// If it is an empty line, skip it
				if len(line) < 1 {
					fmt.Fprintln(w)
					continue
				}
				// If it is a line declaring the error variable,
				// Replace the underscore and proceed
				if strings.TrimSpace(line) == "var _ error" {
					fmt.Fprintln(w, strings.Replace(line, "_", "err", 1))
					continue
				}
				if !inject {
					// If the line starts with the pragma,
					// add error handling into the next line
					if errMsg, isPragma = strings.CutPrefix(strings.TrimSpace(line), pragma); isPragma {
						inject = true
						continue
					}
				}
				if inject {
					tokens := strings.Split(line, " ")
					var errIndex int
					for i, token := range tokens {
						if token == ":=" || token == "=" {
							errIndex = i - 1
							break
						}
					}

					// Identify the indentation level of the line
					var indent string
					for _, c := range tokens[0] {
						if c == '\t' {
							indent += "\t"
						}
					}

					// If there is no other assignment, then add the `err` variable with indentation
					if errIndex == 0 {
						if tokens[errIndex] != indent+"_" {
							log.Println(fmt.Errorf("invalid placeholder in line: %s", line))
							continue
						}
						tokens[errIndex] = indent + "err"
					} else {
						if tokens[errIndex] != "_" {
							log.Println(fmt.Errorf("invalid placeholder in line: %s", line))
							continue
						}
						tokens[errIndex] = "err"
					}
					errLine := strings.Join(tokens, " ")

					// Add the lines with unwrapped error handling
					var errHandler string
					if errMsg == "" {
						errHandler = fmt.Sprintf(errWithoutMsg, indent, indent+"\t", indent)
					} else {
						errHandler = fmt.Sprintf(errWithMsg, indent, indent+"\t", errMsg, indent)
					}

					// Write the modified line and error handler
					fmt.Fprintln(w, errLine)
					fmt.Fprintln(w, errHandler)

					// Reset the flag
					inject = false
				} else {
					fmt.Fprintln(w, s.Text())
				}
			}

			if inPlace {
				os.WriteFile(file, w.Bytes(), 0o644)
			} else {
				fmt.Print(w.String())
			}
		}
	},
}

func Init() {
	// Disable completion generation
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.Flags().BoolVarP(&inPlace, "in-place", "i", false, "Modify files in-place")
}

func Run() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(fmt.Errorf("error running CLI: %w", err))
	}
}
