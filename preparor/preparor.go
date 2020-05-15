package preparor

import (
	"path/filepath"
	"os"
	"fmt"
	"strings"
)

type Preparor struct {
}

func (p *Preparor) Process(path string , minLabelCount , minLabelNameSize int) {
	
	filepaths := p.getFilePaths(path)
	labels := p.getLabels(filepaths , minLabelCount , minLabelNameSize)

	for index , label  := range labels {
		fmt.Printf("%d = %s\n" , index , label)
	}


}

func (p *Preparor) getLabels(filepaths []string, minLabelCount , minLabelNameSize int) []string{


	possibleLabel := ""
	count := 0
	labels := [] string {}

	for _ , path := range filepaths {
		currentLabel := p.getFileFrom(path)

		if possibleLabel == "" {
			possibleLabel = currentLabel
			count += 1
			continue
		}

		if currentLabel == possibleLabel {
			count += 1
			continue
		}


		commonStr := common(possibleLabel , currentLabel)


		if commonStr == ""  || len(commonStr) < minLabelNameSize  {
			if count >= minLabelCount {
				labels = append(labels , possibleLabel)
				count = 0
				possibleLabel = ""
				continue
			}
		}

		possibleLabel = commonStr
		count += 1

	}

	if count >= minLabelCount {
		labels = append(labels , possibleLabel)
	}


	for i := 0 ; i < len(labels) ; i++ {
		labels[i] = cleanString(labels[i])
	}

	return labels
}

func (p *Preparor) getFileFrom(path string) string  {

	data := strings.Split(path , "/")
	filename := data[len(data) - 1]

	return filename
}


func (p *Preparor) getFilePaths(path string) [] string  {

	paths := []string {}

	err := filepath.Walk(path , func(currentPath string , info os.FileInfo, err error) error {
		paths = append(paths , currentPath)
		return nil
	})

	if err != nil {
		fmt.Println(err)
	}

	return paths[1 : len(paths)]
}

func NewPreparor() *Preparor {
	preparor := Preparor{}
	return &preparor
}
