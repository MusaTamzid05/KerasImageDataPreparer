package preparor

import (
	"path/filepath"
	"os"
	"log"
)

type Preparor struct {
}

func (p *Preparor) Process(path string) {
	
	filepaths := p.getFilePaths(path)

	for index , currentPath := range filepaths {
		log.Printf("%d.%s\n" , index , currentPath)
	}

}


func (p *Preparor) getFilePaths(path string) [] string  {

	paths := []string {}

	err := filepath.Walk(path , func(currentPath string , info os.FileInfo, err error) error {
		paths = append(paths , currentPath)
		return nil
	})

	if err != nil {
		log.Println(err)
	}

	return paths
}

func NewPreparor() *Preparor {
	preparor := Preparor{}
	return &preparor
}
