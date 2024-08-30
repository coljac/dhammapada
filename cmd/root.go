/*
Copyright © 2024 Colin Jacobs <colin@coljac.space>

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
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"

	"github.com/fatih/color"

	"github.com/spf13/cobra"
)

func getRandomVerse() (string, string, error) {
	verseNumber := rand.Intn(423) + 1
	return getVerse(verseNumber)
}

func getVerse(verseNumber int) (string, string, error) {
	inVerse := false
	verse := ""
	chapter := ""

	scanner := bufio.NewScanner(strings.NewReader(Dhammapada))

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, fmt.Sprintf("%d.", verseNumber)) {
			inVerse = true
		}
		if len(line) == 0 && inVerse {
			break
		}
		if inVerse {
			verse += strings.TrimPrefix(line, fmt.Sprintf("%d. ", verseNumber)) + "\n"
		} else {
			if strings.HasPrefix(line, "Chapter") {
				chapter = line
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return "", "", err
	}
	return verse, fmt.Sprintf("%s, verse %d", chapter, verseNumber), nil
}

var rootCmd = &cobra.Command{
	Use:   "dhammapada [search]",
	Short: "A daily dose of the Dharma",
	Long:  `A random verse from the Dhammapada, translated by F. Max Müller.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var searchString string
		if len(args) > 0 {
			searchString = args[0]
		}
		colour, _ := cmd.Flags().GetBool("colour")
		verseNumber, _ := cmd.Flags().GetInt("verse")
		if searchString != "" {
			printVersesContaining(searchString, colour)
		} else if verseNumber > 0 {
			verse, chapter, err := getVerse(verseNumber)
			if err != nil {
				log.Fatalf("Error getting verse %d: %v", verseNumber, err)
			}
			printVerse(verse, chapter, colour)
		} else {
			verse, chapter, err := getRandomVerse()
			if err != nil {
				log.Fatalf("Error getting random verse: %v", err)
			}
			printVerse(verse, chapter, colour)
		}
	},
}

func printVerse(verse, chapter string, colour bool) {
	if colour {
		c := color.New(color.FgWhite).Add(color.Bold)
		fmt.Println()
		c.Println(verse)
		c = color.New(color.FgBlue).Add(color.Italic)
		c.Println(chapter)
	} else {
		fmt.Println(verse)
		fmt.Println(chapter)
	}
}

func printVersesContaining(searchString string, colour bool) {
	searchString = strings.ToLower(searchString)
	for i := 1; i <= 423; i++ {
		verse, chapter, err := getVerse(i)
		if err != nil {
			log.Printf("Error getting verse %d: %v", i, err)
			continue
		}
		if strings.Contains(strings.ToLower(verse), searchString) {
			printVerse(verse, chapter, colour)
		}
	}
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("colour", "c", false, "Print with colour")
	rootCmd.Flags().IntP("verse", "v", 0, "Specify a verse number (1-423)")
}
