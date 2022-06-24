package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

func filterArgs(args []string) ([]string, []string) {

	foundFiles := make([]string, 0)
	foundRegex := make([]string, 0)

	// flag regex -smth and -m=#
	flagRegex := regexp.MustCompile("-(color|n|c|v)|-m=[0-9]*")

	// filename regex something.wxyz
	filenameRegex := regexp.MustCompile("[\\w]+\\.[A-Za-z]{2,4}")

	// regex regex any word only gonna use first one
	regexRegex := regexp.MustCompile("\\w")

	for _, element := range args {
		if flagRegex.MatchString(element) {
			// so flags arent confused w/ regex filter these out
		} else if filenameRegex.MatchString(element) {
			foundFiles = append(foundFiles, element)
		} else if regexRegex.MatchString(element) {
			foundRegex = append(foundRegex, element)
		}
	}

	return foundFiles, foundRegex
}

func searchAndPrint(files []string, regexString string, flags []*bool, maxLines *int) {

	regex, err := regexp.Compile(regexString)

	if err != nil {
		fmt.Println("Pattern failed to compile")
		return
	}

	for _, filename := range files {
		file, err := os.OpenFile(filename, os.O_RDONLY, 0777)

		if err != nil {
			fmt.Println("go-grep: " + filename + ": No such file")
			continue // go to next file/iterate loop
		}

		defer file.Close()

		scanner := bufio.NewScanner(file)
		matchCount := 0
		lineCount := 0

		for scanner.Scan() {
			// exit if we've read enough lines, go to next file
			if matchCount >= *maxLines {
				break
			}

			lineText := scanner.Text()

			// XOR
			if regex.MatchString(lineText) != *flags[3] {

				// count matches flag
				if !*flags[2] {
					if len(files) > 1 {
						fmt.Printf(filename + ":")
					}

					// line #'s flag
					if *flags[1] {
						fmt.Printf(strconv.Itoa(lineCount) + ":")
					}

					// text highlight/color flag
					if *flags[0] {
						printLineColor(regex, lineText)
					} else {
						fmt.Println(lineText)
					}
				}
				matchCount++
			}
			lineCount++
		}

		// count flag
		if *flags[2] {
			if len(files) > 1 {
				fmt.Printf(filename + ":" + strconv.Itoa(matchCount) + "\n")
			} else {
				fmt.Println(matchCount)
			}
		}
	}
	fmt.Printf("\n")

}

func printLineColor(regex *regexp.Regexp, currentLineText string) {
	matchIndices := regex.FindAllStringIndex(currentLineText, -1)
	splitText := regex.Split(currentLineText, -1)

	fmt.Printf(splitText[0])
	for i := range matchIndices {
		fmt.Printf("\033[31m" + currentLineText[matchIndices[i][0]:matchIndices[i][1]] + "\033[0m")
		fmt.Printf(splitText[i+1])
	}

	fmt.Printf("\n")
}

func main() {
	files, regex := filterArgs(os.Args[1:])

	colorPtr := flag.Bool("color", false, "Highlight matches")
	numPtr := flag.Bool("n", false, "Display line numbers")
	countPtr := flag.Bool("c", false, "Show number of lines that contain at least one occurrence of the pattern")
	invertPtr := flag.Bool("v", false, "Show a count of all lines not containing pattern")
	maxPtr := flag.Int("m", 1000000, "Limit the maximum number of output lines per file")

	flag.Parse()
	inputFlags := []*bool{}
	inputFlags = append(inputFlags, colorPtr)
	inputFlags = append(inputFlags, numPtr)
	inputFlags = append(inputFlags, countPtr)
	inputFlags = append(inputFlags, invertPtr)

	if len(files) <= 0 || len(regex) <= 0 {
		fmt.Println("Usage: go-grep [OPTION]... PATTERN [FILE]...\nTry 'gogrep -h' for more information.\n")
		return
	}

	searchAndPrint(files, regex[0], inputFlags, maxPtr)
}
