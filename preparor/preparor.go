package preparor

import (
	"path/filepath"
	"os"
	"fmt"
	"strings"
)

type Preparor struct {
}

func (p *Preparor) Process(path , dst string , minLabelCount , minLabelNameSize int , trainSize float64) {
	
	filepaths := p.getFilePaths(path)
	labels := p.getLabels(filepaths , minLabelCount , minLabelNameSize)
	dstPaths := p.getDstPath(dst , labels , filepaths , trainSize)

	for _ , path := range dstPaths {
		fmt.Println(path)
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

func (p *Preparor) getDstPath(dst string , labels , paths []string , trainSize float64) []string {

	labelDstPath := [] string {}
	trainIndex := (int(trainSize * 100 )* len(paths))  / 100
	valTestSize := 1.0 - trainSize
	valSize := valTestSize / 2

	valIndex := ((int(valSize * 100 )* len(paths))  / 100) + trainIndex

	var pathFound bool
	labelType := ""

	for index  , path := range paths {

		if index < trainIndex {
			labelType = "train"
		} else if index < valIndex {
			labelType = "val"
		} else {
			labelType = "test"
		}

		pathFound = false
		for _ , label := range labels {
			filename := p.getFileFrom(path)

			if strings.Contains(filename, label) {
				labelDst := dst + string(os.PathSeparator) +  labelType // create dir here
				labelDst += string(os.PathSeparator) + label 

				labelDstPath = append(labelDstPath ,  labelDst + string(os.PathSeparator) +  filename )
				pathFound = true
				break

			}
		}

		if !pathFound {
			fmt.Printf("Could not get label for %s\n" , path)
		}


	}

	return labelDstPath

}

func NewPreparor() *Preparor {
	preparor := Preparor{}
	return &preparor
}
