package translator

import (
	"log"

	"github.com/xuri/excelize/v2"
)



func readFile(filename string) {
	file, err := excelize.OpenFile(filename)
	if err != nil {
		log.Fatalf("can not open a file %s: %v", filename, err)
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("can not close the file %s: %v", filename, err)
		}
	}()
}


func Run() {
	readFile("input_dataset.xlsx")
}