package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

func main() {
	args := make([]string, 0, len(os.Args))
	for i := 0; i < len(os.Args); i++ {
		arg := os.Args[i]
		switch arg {
		case "-H":
			i++
			head := os.Args[i]
			if !removeHeader(head) {
				args = append(args, arg, formatArg(head))
			}
		case "--data-binary":
			i++
			data := getFormData(os.Args[i])
			args = append(args, data...)
		default:
			if strings.HasPrefix(arg, "-") {
				args = append(args, arg)
			} else {
				args = append(args, formatArg(arg))
			}
		}
	}

	args[0] = "curl"

	fmt.Print(strings.Join(args, " "))
}

func removeHeader(header string) bool {
	switch {
	case strings.HasPrefix(header, "content-type") && strings.Contains(header, "boundary="):
		return true
		// case strings.HasPrefix(header, "dnt:"):
		// 	return true
	}
	return false
}

var nameRegex = regexp.MustCompile("name=\"(.+?)\"")

func getFormData(data string) (formArgs []string) {
	data = strings.Replace(data, "\\r", "", -1)
	data = strings.Replace(data, "\\n", "\n", -1)
	lines := strings.Split(data, "\n")
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if strings.HasPrefix(line, "Content-Disposition") {
			matches := nameRegex.FindAllStringSubmatch(line, -1)
			if len(matches) > 0 {
				i += 2
				formArgs = append(formArgs, "-F", formatArg(matches[0][1]+"="+lines[i]))
			}
		}
	}

	return
}

func formatArg(arg string) string {
	return fmt.Sprintf("\"%s\"", arg)
}
