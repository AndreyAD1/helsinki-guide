package translator

import (
	"fmt"
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
	for _, cell := range []string{"C2", "T2", "U2"} {
		value, err := file.GetCellValue("Lauttasaari", cell)
		if err != nil {
			log.Printf("can not read the cell T2: %v", err)
			continue
		}
		fmt.Println(value)
	}
}


func Run() {
	readFile("input_dataset.xlsx")
}