package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

// removes spaces from integer string, checks if string contains only 1s and 0s
func cleanString(input string) string {
	grab := []rune(input)
	var runes []rune
	for i := 0; i < len(grab); i++ {
		if grab[i] == '1' || grab[i] == '0' {
			runes = append(runes, grab[i])
		}
	}
	if len(runes) == 32 {
		return string(runes)
	} else {
		fmt.Print("Error: bad string length\n")
		return "e"
	}
}

// converts binary to decimal
func binToDec(input string) int {
	chars := strings.SplitAfter(input, "")
	if len(chars) == 0 {
		return 0
	}
	firstNum, err := strconv.Atoi(chars[0])
	if err != nil {
		log.Fatalf("Failed to convert binary character %s", err)
	}

	return firstNum*int(math.Pow(2, float64(len(chars)-1))) + binToDec(strings.Join(chars[1:], ""))

}

func twosComplementToDec(input string) int {
	chars := strings.SplitAfter(input, "")
	num, _ := strconv.Atoi(chars[0])
	secondNum, _ := strconv.Atoi(chars[1])
	num = (num * -2) + secondNum
	for i := 1; i < len(chars)-1; i++ {
		nextNum, _ := strconv.Atoi(chars[i+1])
		num = (num * 2) + nextNum
	}
	return num
}

// gets opCode from binary
func getOpCode(input string, count int, readingData *bool, asmCode *[]Line) string {

	retString := ""

	//6-digit opCode
	//two print lines per case; one for binary and one for opcode + registers
	opcode := binToDec(input[0:6])
	switch {
	case opcode == 5:
		retString = input[0:6] + " " + input[6:32] + "\t" + strconv.Itoa(count) + "\tB\t#" + strconv.Itoa(twosComplementToDec(input[6:32])) + "\n"
		*asmCode = append(*asmCode, Line{5, 0, 0, 0, twosComplementToDec(input[6:32]), 0, 0, "B\t#" + strconv.Itoa(twosComplementToDec(input[6:32]))})

	default:
		opcode = binToDec(input[0:8])
		switch {
		case opcode == 180:
			retString = input[0:8] + " " + input[8:27] + " " + input[27:32] + "\t" + strconv.Itoa(count) + "\tCBZ\tR" + strconv.Itoa(binToDec(input[27:32])) + ", #" + strconv.Itoa(twosComplementToDec(input[8:27])) + "\n"
			*asmCode = append(*asmCode, Line{180, 0, binToDec(input[27:32]), 0, twosComplementToDec(input[8:27]), 0, 0, "CBZ\tR" + strconv.Itoa(binToDec(input[27:32])) + ", #" + strconv.Itoa(twosComplementToDec(input[8:27]))})

		case opcode == 181:
			retString = input[0:8] + " " + input[8:27] + " " + input[27:32] + "\t" + strconv.Itoa(count) + "\tCBNZ\tR" + strconv.Itoa(binToDec(input[27:32])) + ", #" + strconv.Itoa(twosComplementToDec(input[8:27])) + "\n"
			*asmCode = append(*asmCode, Line{181, 0, binToDec(input[27:32]), 0, twosComplementToDec(input[8:27]), 0, 0, "CBNZ\tR" + strconv.Itoa(binToDec(input[27:32])) + ", #" + strconv.Itoa(twosComplementToDec(input[8:27]))})

		default:
			opcode = binToDec(input[0:9])
			switch {
			case opcode == 421:
				retString = input[0:9] + " " + input[9:11] + " " + input[11:27] + " " + input[27:32] + "\t" + strconv.Itoa(count) + "\tMOVZ\tR" + strconv.Itoa(binToDec(input[27:32])) + ", " + strconv.Itoa(binToDec(input[11:27])) + ", LSL " + strconv.Itoa(binToDec(input[9:11])*16) + "\n"
				*asmCode = append(*asmCode, Line{421, binToDec(input[27:32]), 0, 0, binToDec(input[11:27]), 0, binToDec(input[9:11]) * 16, "MOVZ\tR" + strconv.Itoa(binToDec(input[27:32])) + ", " + strconv.Itoa(binToDec(input[11:27])) + ", LSL " + strconv.Itoa(binToDec(input[9:11])*16)})

			case opcode == 485:
				retString = input[0:9] + " " + input[9:11] + " " + input[11:27] + " " + input[27:32] + "\t" + strconv.Itoa(count) + "\tMOVK\tR" + strconv.Itoa(binToDec(input[27:32])) + ", " + strconv.Itoa(binToDec(input[11:27])) + ", LSL " + strconv.Itoa(binToDec(input[9:11])*16) + "\n"
				*asmCode = append(*asmCode, Line{485, binToDec(input[27:32]), 0, 0, binToDec(input[11:27]), 0, binToDec(input[9:11]) * 16, "MOVZ\tR" + strconv.Itoa(binToDec(input[27:32])) + ", " + strconv.Itoa(binToDec(input[11:27])) + ", LSL " + strconv.Itoa(binToDec(input[9:11])*16)})

			default:
				opcode = binToDec(input[0:10])
				switch {
				case opcode == 580:
					retString = input[0:10] + " " + input[10:22] + " " + input[22:27] + " " + input[27:32] + "\t" + strconv.Itoa(count) + "\tADDI\tR" + strconv.Itoa(binToDec(input[27:32])) + ", R" + strconv.Itoa(binToDec(input[22:27])) + ", #" + strconv.Itoa(binToDec(input[10:22])) + "\n"
					*asmCode = append(*asmCode, Line{580, binToDec(input[27:32]), binToDec(input[22:27]), 0, binToDec(input[10:22]), 0, 0, "ADDI\tR" + strconv.Itoa(binToDec(input[27:32])) + ", R" + strconv.Itoa(binToDec(input[22:27])) + ", #" + strconv.Itoa(binToDec(input[10:22]))})

				case opcode == 836:
					retString = input[0:10] + " " + input[10:22] + " " + input[22:27] + " " + input[27:32] + "\t" + strconv.Itoa(count) + "\tSUBI\tR" + strconv.Itoa(binToDec(input[27:32])) + ", R" + strconv.Itoa(binToDec(input[22:27])) + ", #" + strconv.Itoa(binToDec(input[10:22])) + "\n"
					*asmCode = append(*asmCode, Line{836, binToDec(input[27:32]), binToDec(input[22:27]), 0, binToDec(input[10:22]), 0, 0, "SUBI\tR" + strconv.Itoa(binToDec(input[27:32])) + ", R" + strconv.Itoa(binToDec(input[22:27])) + ", #" + strconv.Itoa(binToDec(input[10:22]))})

				default:
					opcode = binToDec(input[0:11])
					switch {
					case opcode == 1104:
						retString = input[0:11] + " " + input[11:16] + " " + input[16:22] + " " + input[22:27] + " " + input[27:32] + "\t" + strconv.Itoa(count) + "\tAND\tR" + strconv.Itoa(binToDec(input[27:32])) + ", R" + strconv.Itoa(binToDec(input[22:27])) + ", R" + strconv.Itoa(binToDec(input[11:16])) + "\n"
						*asmCode = append(*asmCode, Line{1104, binToDec(input[27:32]), binToDec(input[22:27]), binToDec(input[11:16]), 0, 0, 0, "AND\tR" + strconv.Itoa(binToDec(input[27:32])) + ", R" + strconv.Itoa(binToDec(input[22:27])) + ", R" + strconv.Itoa(binToDec(input[11:16]))})

					case opcode == 1112:
						retString = input[0:11] + " " + input[11:16] + " " + input[16:22] + " " + input[22:27] + " " + input[27:32] + "\t" + strconv.Itoa(count) + "\tADD\tR" + strconv.Itoa(binToDec(input[27:32])) + ", R" + strconv.Itoa(binToDec(input[22:27])) + ", R" + strconv.Itoa(binToDec(input[11:16])) + "\n"
						*asmCode = append(*asmCode, Line{1112, binToDec(input[27:32]), binToDec(input[22:27]), binToDec(input[11:16]), 0, 0, 0, "ADD\tR" + strconv.Itoa(binToDec(input[27:32])) + ", R" + strconv.Itoa(binToDec(input[22:27])) + ", R" + strconv.Itoa(binToDec(input[11:16]))})

					case opcode == 1360:
						retString = input[0:11] + " " + input[11:16] + " " + input[16:22] + " " + input[22:27] + " " + input[27:32] + "\t" + strconv.Itoa(count) + "\tORR\tR" + strconv.Itoa(binToDec(input[27:32])) + ", R" + strconv.Itoa(binToDec(input[22:27])) + ", R" + strconv.Itoa(binToDec(input[11:16])) + "\n"
						*asmCode = append(*asmCode, Line{1360, binToDec(input[27:32]), binToDec(input[22:27]), binToDec(input[11:16]), 0, 0, 0, "ORR\tR" + strconv.Itoa(binToDec(input[27:32])) + ", R" + strconv.Itoa(binToDec(input[22:27])) + ", R" + strconv.Itoa(binToDec(input[11:16]))})

					case opcode == 1624:
						retString = input[0:11] + " " + input[11:16] + " " + input[16:22] + " " + input[22:27] + " " + input[27:32] + "\t" + strconv.Itoa(count) + "\tSUB\tR" + strconv.Itoa(binToDec(input[27:32])) + ", R" + strconv.Itoa(binToDec(input[22:27])) + ", R" + strconv.Itoa(binToDec(input[11:16])) + "\n"
						*asmCode = append(*asmCode, Line{1624, binToDec(input[27:32]), binToDec(input[22:27]), binToDec(input[11:16]), 0, 0, 0, "SUB\tR" + strconv.Itoa(binToDec(input[27:32])) + ", R" + strconv.Itoa(binToDec(input[22:27])) + ", R" + strconv.Itoa(binToDec(input[11:16]))})

					case opcode == 1690:
						retString = input[0:11] + " " + input[11:16] + " " + input[16:22] + " " + input[22:27] + " " + input[27:32] + "\t" + strconv.Itoa(count) + "\tLSR\tR" + strconv.Itoa(binToDec(input[27:32])) + ", R" + strconv.Itoa(binToDec(input[22:27])) + ", #" + strconv.Itoa(binToDec(input[16:22])) + "\n"
						*asmCode = append(*asmCode, Line{1690, binToDec(input[27:32]), binToDec(input[22:27]), 0, binToDec(input[16:22]), 0, 0, "LSR\tR" + strconv.Itoa(binToDec(input[27:32])) + ", R" + strconv.Itoa(binToDec(input[22:27])) + ", #" + strconv.Itoa(binToDec(input[16:22]))})

					case opcode == 1691:
						retString = input[0:11] + " " + input[11:16] + " " + input[16:22] + " " + input[22:27] + " " + input[27:32] + "\t" + strconv.Itoa(count) + "\tLSL\tR" + strconv.Itoa(binToDec(input[27:32])) + ", R" + strconv.Itoa(binToDec(input[22:27])) + ", #" + strconv.Itoa(binToDec(input[16:22])) + "\n"
						*asmCode = append(*asmCode, Line{1691, binToDec(input[27:32]), binToDec(input[22:27]), 0, binToDec(input[16:22]), 0, 0, "LSL\tR" + strconv.Itoa(binToDec(input[27:32])) + ", R" + strconv.Itoa(binToDec(input[22:27])) + ", #" + strconv.Itoa(binToDec(input[16:22]))})

					case opcode == 1984:
						retString = input[0:11] + " " + input[11:20] + " " + input[20:22] + " " + input[22:27] + " " + input[27:32] + "\t" + strconv.Itoa(count) + "\tSTUR\tR" + strconv.Itoa(binToDec(input[27:32])) + ", [R" + strconv.Itoa(binToDec(input[22:27])) + ", #" + strconv.Itoa(binToDec(input[11:20])) + "]\n"
						*asmCode = append(*asmCode, Line{1984, binToDec(input[27:32]), binToDec(input[22:27]), 0, binToDec(input[11:20]), 0, 0, "STUR\tR" + strconv.Itoa(binToDec(input[27:32])) + ", [R" + strconv.Itoa(binToDec(input[22:27])) + ", #" + strconv.Itoa(binToDec(input[11:20])) + "]"})

					case opcode == 1986:
						retString = input[0:11] + " " + input[11:20] + " " + input[20:22] + " " + input[22:27] + " " + input[27:32] + "\t" + strconv.Itoa(count) + "\tLDUR\tR" + strconv.Itoa(binToDec(input[27:32])) + ", [R" + strconv.Itoa(binToDec(input[22:27])) + ", #" + strconv.Itoa(binToDec(input[11:20])) + "]\n"
						*asmCode = append(*asmCode, Line{1986, binToDec(input[27:32]), binToDec(input[22:27]), 0, binToDec(input[11:20]), 0, 0, "LDUR\tR" + strconv.Itoa(binToDec(input[27:32])) + ", [R" + strconv.Itoa(binToDec(input[22:27])) + ", #" + strconv.Itoa(binToDec(input[11:20])) + "]"})

					case opcode == 1692:
						retString = input[0:11] + " " + input[11:16] + " " + input[16:22] + " " + input[22:27] + " " + input[27:32] + "\t" + strconv.Itoa(count) + "\tASR\tR" + strconv.Itoa(binToDec(input[27:32])) + ", R" + strconv.Itoa(binToDec(input[22:27])) + ", R" + strconv.Itoa(binToDec(input[11:16])) + "\n"
						*asmCode = append(*asmCode, Line{1692, binToDec(input[27:32]), binToDec(input[22:27]), binToDec(input[11:16]), 0, 0, 0, "ASR\tR" + strconv.Itoa(binToDec(input[27:32])) + ", R" + strconv.Itoa(binToDec(input[22:27])) + ", R" + strconv.Itoa(binToDec(input[11:16]))})

					case opcode == 1872:
						retString = input[0:11] + " " + input[11:16] + " " + input[16:22] + " " + input[22:27] + " " + input[27:32] + "\t" + strconv.Itoa(count) + "\tEOR\tR" + strconv.Itoa(binToDec(input[27:32])) + ", R" + strconv.Itoa(binToDec(input[22:27])) + ", R" + strconv.Itoa(binToDec(input[11:16])) + "\n"
						*asmCode = append(*asmCode, Line{1872, binToDec(input[27:32]), binToDec(input[22:27]), binToDec(input[11:16]), 0, 0, 0, "EOR\tR" + strconv.Itoa(binToDec(input[27:32])) + ", R" + strconv.Itoa(binToDec(input[22:27])) + ", R" + strconv.Itoa(binToDec(input[11:16]))})

					case opcode == 2038:
						*readingData = true
						retString = input + "\t" + strconv.Itoa(count) + "\tBREAK\n"
						retString = input[0:8] + " " + input[8:11] + " " + input[11:16] + " " + input[16:21] + " " + input[21:26] + " " + input[26:32] + "\t" + strconv.Itoa(count) + "\tBREAK\n"
						*asmCode = append(*asmCode, Line{2038, 0, 0, 0, 0, 0, 0, "BREAK"})

					case opcode == 0:
						retString = input + "\t" + strconv.Itoa(count) + "\tNOP\n"
						*asmCode = append(*asmCode, Line{0, 0, 0, 0, 0, 0, 0, "NOP"})

					default:
						retString = "Unknown Instruction\n"
					}
				}
			}
		}
	}
	return retString
}

// core of the simulator. reads instruction and modifies registers/PC/data
func readInstruction(instruction Line, pc *int, registers *[32]int, cycle *int, simFile *os.File, dataset *[]Data, startingAdd *int, readStart *bool) {
	switch instruction.opCode {
	//ADD
	case 1112:
		registers[instruction.destination] = registers[instruction.r1] + registers[instruction.r2]
		//SUB
	case 1624:
		registers[instruction.destination] = registers[instruction.r1] - registers[instruction.r2]
		//ADDI
	case 580:
		registers[instruction.destination] = registers[instruction.r1] + instruction.immediate
		//SUBI
	case 836:
		registers[instruction.destination] = registers[instruction.r1] - instruction.immediate
		//AND
	case 1104:
		registers[instruction.destination] = registers[instruction.r1] & registers[instruction.r2]
		//ORR
	case 1360:
		registers[instruction.destination] = registers[instruction.r1] | registers[instruction.r2]
		//EOR
	case 1872:
		registers[instruction.destination] = registers[instruction.r1] ^ registers[instruction.r2]
		//B
	case 5:
		printSim(instruction, *registers, *cycle, simFile, *dataset, *pc)
		*pc += instruction.immediate
		//CBZ
	case 180:
		if registers[instruction.r1] == 0 {
			printSim(instruction, *registers, *cycle, simFile, *dataset, *pc)
			*pc += instruction.immediate
		} else {
			printSim(instruction, *registers, *cycle, simFile, *dataset, *pc)
			*pc++
		}
		//CBNZ
	case 181:
		if registers[instruction.r1] != 0 {
			printSim(instruction, *registers, *cycle, simFile, *dataset, *pc)
			*pc += instruction.immediate
		} else {
			printSim(instruction, *registers, *cycle, simFile, *dataset, *pc)
			*pc++
		}
		//LSR
	case 1690:
		if registers[instruction.r1] < 0 {
			binNum := int32(registers[instruction.r1])
			binNum = (binNum ^ 0xFFFF) + 1
			tcNum := strconv.FormatInt(int64(binNum), 2)[1:]
			tcNum = tcNum[:len(tcNum)-1]
			newNum := ""
			for i := 0; i < 32; i++ {
				if i == 0 {
					newNum = newNum + "0"
				} else if i >= 32-len(tcNum) {
					newNum = newNum + string(tcNum[i-(32-len(tcNum))])
				} else {
					newNum = newNum + "1"
				}
			}
			registers[instruction.destination] = twosComplementToDec(newNum)
		} else {
			registers[instruction.destination] = registers[instruction.r1] >> instruction.immediate
		}

		//LSL
	case 1691:
		registers[instruction.destination] = registers[instruction.r1] << instruction.immediate
		//ASR
	case 1692:
		registers[instruction.destination] = registers[instruction.r1] >> registers[instruction.r2]
		//STUR
	case 1984:
		if findIndex(registers[instruction.r1]+(instruction.immediate*4), *dataset) != -1 {
			newData := *dataset
			newData[findIndex(registers[instruction.r1]+(instruction.immediate*4), *dataset)].value = registers[instruction.destination]
			*dataset = newData
		} else {
			if *readStart {
				*startingAdd = registers[instruction.r1] + (instruction.immediate * 4)
				*readStart = false
			}
			address := registers[instruction.r1] + (instruction.immediate * 4)
			noOffsetAddress := 0
			for i := address; i > address-32; i -= 4 {
				temp := i - *startingAdd
				if temp%32 == 0 {
					noOffsetAddress = temp
				}
			}
			noOffsetAddress += *startingAdd
			for i := 0; i < 8; i++ {
				*dataset = append(*dataset, Data{noOffsetAddress + (i * 4), 0})
			}
			newData := *dataset
			newData[findIndex(registers[instruction.r1]+(instruction.immediate*4), *dataset)].value = registers[instruction.destination]
			*dataset = newData
		}
		//LDUR
	case 1986:
		newData := *dataset
		indexForLoad := findIndex(registers[instruction.r1]+(instruction.immediate*4), *dataset)
		if indexForLoad != -1 {
			registers[instruction.destination] = newData[indexForLoad].value
		} else {
			registers[instruction.destination] = 0
		}

		//BREAK
	case 2038:
		printSim(instruction, *registers, *cycle, simFile, *dataset, *pc)
		*pc = -1
	default:
		*pc++
		break
	}
	//this order of execution -> print -> PC will guarantee that the sim will run in the correct order
	if instruction.opCode != 5 && instruction.opCode != 180 && instruction.opCode != 181 && instruction.opCode != 2038 {
		printSim(instruction, *registers, *cycle, simFile, *dataset, *pc)
		*pc++
	}
	*cycle++
}

func findIndex(target int, dataset []Data) int {
	retInt := -1
	for i := range dataset {
		if dataset[i].address == target {
			retInt = i
		}
	}
	return retInt
}

// outputs simulator to file
func printSim(instruction Line, registers [32]int, cycle int, simFile *os.File, dataset []Data, pc int) {
	//line divider
	_, err := io.WriteString(simFile, "====================\n")
	if err != nil {
		return
	}
	//prints cycle num, instruction num, and instruction string
	_, err = io.WriteString(simFile, "Cycle: "+strconv.Itoa(cycle)+"\t"+strconv.Itoa(96+(pc*4))+"\t "+instruction.opString+"\n")
	if err != nil {
		return
	}
	_, err = io.WriteString(simFile, "\n")
	if err != nil {
		return
	}
	//prints registers
	_, err = io.WriteString(simFile, "Registers:\n")
	if err != nil {
		return
	}
	for i := 1; i <= 32; i++ {
		if (i-1)%8 == 0 {
			_, err = io.WriteString(simFile, "R"+strconv.Itoa(i-1)+":")
			if err != nil {
				return
			}
		}
		_, err = io.WriteString(simFile, "\t"+strconv.Itoa(registers[i-1]))
		if err != nil {
			return
		}
		if i%8 == 0 {
			_, err = io.WriteString(simFile, "\n")
			if err != nil {
				return
			}
		}
	}
	_, err = io.WriteString(simFile, "\n")
	if err != nil {
		return
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

}

type Line struct {
	opCode      int
	destination int
	r1          int
	r2          int
	immediate   int
	shamt       int
	offset      int
	opString    string
}

type Data struct {
	address int
	value   int
}

func main() {
	//grab file name pointers
	var InputFileName *string = flag.String("i", "", "Gets the input file name")
	var OutputFileName *string = flag.String("o", "", "Gets the output file name")
	flag.Parse()

	if flag.NArg() != 0 {
		log.Fatalf("Inappropriate number of arguments")
	}

	//scan by opening file from pointer
	file, err := os.Open(*InputFileName)
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}

	//scan
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var txtlines []string
	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}

	//opening output file
	outFile, err := os.Create(*OutputFileName + "_dis.txt")
	if err != nil {
		log.Fatalf("Failed to open display output file.")
	}

	//opening output file
	simFile, err := os.Create(*OutputFileName + "_sim.txt")
	if err != nil {
		log.Fatalf("Failed to open simulator output file.")
	}

	//boolean keeping track of opCode versus data read
	readingData := false
	//dynamic array keeping track of opcodes, registers, and immediate values
	asmCode := make([]Line, 0)
	readStart := true
	startingAdd := -1
	//dynamic array keeping track of data
	dataset := make([]Data, 0)
	//static 2D array keeping track of all 32 registers
	registers := [32]int{
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
	}

	//loop through list and output into file
	for i := range txtlines {
		iString := cleanString(txtlines[i])
		if iString != "e" {
			if readingData {
				if readStart {
					startingAdd = 96 + (i * 4)
					dataset = append(dataset, Data{startingAdd, twosComplementToDec(iString)})
					for j := 1; j < 8; j++ {
						dataset = append(dataset, Data{startingAdd + (j * 4), 0})
					}
					readStart = false
				} else {
					currentAdd := 96 + (i * 4)
					if (currentAdd-startingAdd)%32 == 0 {
						dataset = append(dataset, Data{currentAdd, twosComplementToDec(iString)})
						for j := 1; j < 8; j++ {
							dataset = append(dataset, Data{96 + ((i + j) * 4), 0})
						}
					} else {
						dataset[findIndex(96+(i*4), dataset)] = Data{96 + (i * 4), twosComplementToDec(iString)}
					}
				}
				dataString := iString + "\t" + strconv.Itoa(96+(i*4)) + "\t" + strconv.Itoa(twosComplementToDec(iString)) + "\n"
				_, writeErr := io.WriteString(outFile, dataString)
				if writeErr != nil {
					log.Fatalf("Failed to write into file.")
				}
			} else {
				_, writeErr := io.WriteString(outFile, getOpCode(iString, 96+(i*4), &readingData, &asmCode))
				if writeErr != nil {
					log.Fatalf("Failed to write into file.")
				}
			}
		}
	}

	//SIMULATOR CODE
	pc := 0
	cycle := 1
	for {
		if pc != -1 {
			readInstruction(asmCode[pc], &pc, &registers, &cycle, simFile, &dataset, &startingAdd, &readStart)
		} else {
			break
		}
	}

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
}
