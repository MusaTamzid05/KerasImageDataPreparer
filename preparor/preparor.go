package preparor

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Preparor struct {
}

func (p *Preparor) Process(path, dst string, minLabelCount, minLabelNameSize int, trainSize float64) {

	filepaths := p.getFilePaths(path)

	if len(filepaths) == 0 {
		log.Printf("Not files found in %s.\n", path)
	}

	labels := p.getLabels(filepaths, minLabelCount, minLabelNameSize)
	labelCountInfo, labelPaths := p.getUniqueLabelCount(labels, filepaths)

	for label, size := range labelCountInfo {
		fmt.Printf("%s = %d\n", label, size)
	}

	pathData := p.getDstPath(dst, labelCountInfo, labelPaths, trainSize)

	for key, value := range pathData {
		fmt.Printf("%s = %s\n", key, value)
	}

	var err error

	err = p.initDir(dst, labels)

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Preparing files")
	for srcPath, dstPath := range pathData {

		fmt.Printf("%s => %s\n", srcPath, dstPath)
		err = CopyFile(srcPath, dstPath)

		if err != nil {
			fmt.Println(err)
		}
	}

}

func (p *Preparor) initDir(dst string, labels []string) error {

	if Exists(dst) {
		log.Fatalf("The dst dir '%s' already exists", dst)
	}

	err := MakeDir(dst)

	if err != nil {
		log.Fatalln(err)
	}

	var currentPath string

	for _, folderName := range []string{"train", "test", "val"} {
		currentPath = dst + string(os.PathSeparator) + folderName
		err = MakeDir(currentPath)

		if err != nil {
			log.Fatalln(err)
		}

		rootDirPath := currentPath
		for _, label := range labels {
			currentDirPath := rootDirPath + string(os.PathSeparator) + label
			err = MakeDir(currentDirPath)

			if err != nil {
				log.Fatalln(err)
			}
		}
	}

	return nil
}

func (p *Preparor) getUniqueLabelCount(labels, paths []string) (map[string]int, map[string][]string) {

	/*
		TODO:This should be done when unique labels are found
	*/
	labelInfo := make(map[string]int)
	labelPath := make(map[string][]string)

	for _, label := range labels {
		labelInfo[label] = 0
		labelPath[label] = []string{}
	}

	for _, path := range paths {

		for _, label := range labels {
			filename := p.getFileFrom(path)
			if strings.Contains(filename, label) {
				labelInfo[label] += 1
				labelPath[label] = append(labelPath[label], path)
			}
		}
	}

	return labelInfo, labelPath
}

func (p *Preparor) getLabels(filepaths []string, minLabelCount, minLabelNameSize int) []string {

	possibleLabel := ""
	count := 0
	labels := []string{}

	for _, path := range filepaths {
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

		commonStr := common(possibleLabel, currentLabel)

		if commonStr == "" || len(commonStr) < minLabelNameSize {
			if count >= minLabelCount {
				labels = append(labels, possibleLabel)
				count = 0
				possibleLabel = ""
				continue
			}
		}

		possibleLabel = commonStr
		count += 1

	}

	if count >= minLabelCount {
		labels = append(labels, possibleLabel)
	}

	for i := 0; i < len(labels); i++ {
		labels[i] = cleanString(labels[i])
	}

	return labels
}

func (p *Preparor) getFileFrom(path string) string {

	data := strings.Split(path, "/")
	filename := data[len(data)-1]

	return filename
}

func (p *Preparor) getFilePaths(path string) []string {

	paths := []string{}

	err := filepath.Walk(path, func(currentPath string, info os.FileInfo, err error) error {
		paths = append(paths, currentPath)
		return nil
	})

	if err != nil {
		fmt.Println(err)
	}

	return paths[1:len(paths)]
}

func (p *Preparor) getDstPath(dst string, labelCountInfo map[string]int, labelPaths map[string][]string, trainSize float64) map[string]string {

	labelDstPathData := make(map[string]string)

	labels := getKeysFrom(labelCountInfo)

	for label, size := range labelCountInfo {
		trainIndex := (int(trainSize*100) * size) / 100
		valTestSize := 1.0 - trainSize
		valSize := valTestSize / 2

		valIndex := ((int(valSize*100) * size) / 100) + trainIndex

		labelType := ""

		temp := ""

		for index, path := range labelPaths[label] {

			if index < trainIndex {
				labelType = "train"
			} else if index < valIndex {
				labelType = "val"
			} else {
				labelType = "test"
			}

			if labelType != temp {
				temp = labelType
				fmt.Printf("%d = %s,%s\n", index, temp, label)
			}

			filename := p.getFileFrom(path)
			fileFound := false

			for _, fileLabel := range labels {

				if strings.Contains(filename, fileLabel) {
					labelDst := dst + string(os.PathSeparator) + labelType // create dir here
					labelDst += string(os.PathSeparator) + fileLabel + string(os.PathSeparator) + filename
					labelDstPathData[path] = labelDst
					fileFound = true

				}
			}

			if !fileFound {
				fmt.Println("No label found for ", path)

			}

		}
	}

	return labelDstPathData

}

func NewPreparor() *Preparor {
	preparor := Preparor{}
	return &preparor
}
