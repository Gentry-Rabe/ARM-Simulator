package main

import (
	"bufio"
	"flag"
	"log"
	"os"
)

type Instruction struct {
	opCode      int
	destination int
	r1          int
	r2          int
	immediate   int
	shamt       int
	offset      int
	opString    string
	result      int
}
type Data struct {
	address int
	info    [8]int
}

// dynamic array keeping track of opcodes, registers, and immediate values
var asmCode = make([]Instruction, 0)

// dynamic array keeping track of data
var dataset = make([]Data, 0)

// static 2D array keeping track of all 32 registers
var registers = [32]int{
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
}

// boolean keeping track of opCode versus data read
var readingData = false

// array for reading in machine code
var txtlines []string

func main() {
	//grab file name pointers
	var InputFileName = flag.String("i", "", "Gets the input file name")
	var OutputFileName = flag.String("o", "", "Gets the output file name")
	flag.Parse()

	if flag.NArg() != 0 {
		log.Fatalf("Inappropriate number of arguments")
	}

	//scan by opening file from pointer
	file, err := os.Open(*InputFileName)
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}

	//scan into txtlines
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}

	//opening display output file
	outFile, err := os.Create(*OutputFileName + "_dis.txt")
	if err != nil {
		log.Fatalf("Failed to open display output file.")
	}

	//opening simulator output file
	simFile, err := os.Create(*OutputFileName + "_sim.txt")
	if err != nil {
		log.Fatalf("Failed to open simulator output file.")
	}

	//opening pipeline simulator output file
	pipelineFile, err := os.Create(*OutputFileName + "_pipeline.txt")
	if err != nil {
		log.Fatalf("Failed to open pipeline output file.")
	}

	//DISPLAY CODE
	dis(outFile)

	//SIMULATOR CODE
	sim(simFile)

	//PIPELINE CODE
	pipeline(pipelineFile)

	//CLOSING FILES
	err = outFile.Close()
	if err != nil {
		log.Fatalf("Failed to close output file.")
	}
	err = file.Close()
	if err != nil {
		log.Fatalf("Failed to close input file.")
	}
	err = simFile.Close()
	if err != nil {
		log.Fatalf("Failed to close simulator file.")
	}
	err = pipelineFile.Close()
	if err != nil {
		log.Fatalf("Failed to close pipeline file.")
	}
}

/*

func findIndex(target int, dataset []Data) int {
	retInt := -1
	for i := range dataset {
		if dataset[i].address == target {
			retInt = i
		}
	}
	return retInt
}

//prints data
_, err = io.WriteString(simFile, "Data:\n")
if err != nil {
	return
}
for i := range dataset {
	if i%8 == 0 {
		_, err = io.WriteString(simFile, strconv.Itoa(dataset[i].address)+": ")
	}
	_, err = io.WriteString(simFile, "\t"+strconv.Itoa(dataset[i].value))
	if i%8 == 7 {
		_, err = io.WriteString(simFile, "\n")
	}
	if err != nil {
		return
	}
}
*/
