package preparor

import (
	"fmt"
	"io"
	"os"
)

func common(str1, str2 string) string {

	matches := []string{}

	for index, _ := range str1 {
		for index2, _ := range str2 {
			if str1[index] == str2[index2] {

				if string(str1[index]) == " " {
					continue
				}

				i := index
				j := index2
				match := ""
				matchFound := true

				for matchFound {

					if j >= len(str2) || i >= len(str1) {
						matchFound = false

						if match != "" {
							matches = append(matches, match)
						}
						continue
					}

					if str1[i] != str2[j] {
						matches = append(matches, match)
						matchFound = false
						continue
					}

					match += string(str1[i])
					i += 1
					j += 1

				}
			}
		}
	}

	biggestMatch := ""
	max := -1

	for _, match := range matches {
		if len(match) > max {
			max = len(match)
			biggestMatch = match
		}
	}

	return biggestMatch

}

func contains(char rune, dataList []rune) bool {

	for _, dataChar := range dataList {
		if char == dataChar {
			return true
		}
	}

	return false
}

func cleanString(str string) string {

	result := ""

	ignoreChars := []rune{',', '.', '!', '-', '?'}

	for _, char := range str {
		if !contains(char, ignoreChars) {
			result += string(char)
		}
	}

	return result
}

func CopyFile(src, dst string) error {

	sourceFileStat, err := os.Stat(src)

	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)

	if err != nil {
		return err
	}

	defer source.Close()

	destination, err := os.Create(dst)

	if err != nil {
		return err
	}

	defer destination.Close()
	_, err = io.Copy(destination, source)

	return nil
}

func Exists(path string) bool {

	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		return false
	}
	return true
}

func MakeDir(path string) error {
	return os.Mkdir(path, 0755)
}
