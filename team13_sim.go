package main

import (
	"io"
	"os"
	"strconv"
)

func readInstruction(instruction Instruction, pc *int, registers *[32]int, cycle *int, simFile *os.File) {
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
		printSim(instruction, *registers, *cycle, simFile)
		*pc += instruction.immediate
		//CBZ
	case 180:
		if registers[instruction.r1] == 0 {
			printSim(instruction, *registers, *cycle, simFile)
			*pc += instruction.immediate
		} else {
			printSim(instruction, *registers, *cycle, simFile)
			*pc++
		}
		//CBNZ
	case 181:
		if registers[instruction.r1] != 0 {
			printSim(instruction, *registers, *cycle, simFile)
			*pc += instruction.immediate
		} else {
			printSim(instruction, *registers, *cycle, simFile)
			*pc++
		}
		//LSR
	case 1690:
		registers[instruction.destination] = registers[instruction.r1] >> instruction.immediate
		//LSL
	case 1691:
		registers[instruction.destination] = registers[instruction.r1] << instruction.immediate
		//ASR
	case 1692:
		registers[instruction.destination] = registers[instruction.r1] >> registers[instruction.r2]
		//STUR
	case 1984:
		if findIndex(registers[instruction.r2]+(instruction.immediate*4), dataset) != -1 {
			newData := dataset
			newData[findIndex(registers[instruction.r2]+(instruction.immediate*4), dataset)].value = registers[instruction.r1]
			dataset = newData
		} else {
			if readStart {
				startingAdd = registers[instruction.r2] + (instruction.immediate * 4)
				readStart = false
			}
			address := registers[instruction.r2] + (instruction.immediate * 4)
			noOffsetAddress := 0
			for i := address; i > address-32; i -= 4 {
				temp := i - startingAdd
				if temp%32 == 0 {
					noOffsetAddress = temp
				}
			}
			noOffsetAddress += startingAdd
			for i := 0; i < 8; i++ {
				dataset = append(dataset, Data{noOffsetAddress + (i * 4), 0})
			}
			newData := dataset
			newData[findIndex(registers[instruction.r2]+(instruction.immediate*4), dataset)].value = registers[instruction.r1]
			dataset = newData
		}
		//LDUR
	case 1986:
		newData := dataset
		indexForLoad := findIndex(registers[instruction.r1]+(instruction.immediate*4), dataset)
		if indexForLoad != -1 {
			registers[instruction.destination] = newData[indexForLoad].value
		} else {
			registers[instruction.destination] = 0
		}
		//BREAK
	case 2038:
		printSim(instruction, *registers, *cycle, simFile)
		*pc = -1
	default:
		*pc++
		break
	}
	//this order of execution -> print -> PC will guarantee that the sim will run in the correct order
	if instruction.opCode != 5 && instruction.opCode != 180 && instruction.opCode != 181 && instruction.opCode != 2038 {
		printSim(instruction, *registers, *cycle, simFile)
		*pc++
	}
	*cycle++
}

// outputs simulator to file
func printSim(instruction Instruction, registers [32]int, cycle int, simFile *os.File) {
	//line divider
	_, err := io.WriteString(simFile, "====================\n")
	if err != nil {
		return
	}
	//prints cycle num, instruction num, and instruction string
	_, err = io.WriteString(simFile, "Cycle: "+strconv.Itoa(cycle)+"\t"+strconv.Itoa(instruction.opCode)+"\t "+instruction.opString+"\n")
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
	for i, element := range dataset {
		if i%8 == 0 {
			_, err = io.WriteString(simFile, strconv.Itoa(element.address)+": ")
		}
		_, err = io.WriteString(simFile, "\t"+strconv.Itoa(element.value))
		if i%8 == 7 {
			_, err = io.WriteString(simFile, "\n")
		}
		if err != nil {
			return
		}
	}

}

func findIndex(target int, funcDataset []Data) int {
	retInt := -1
	for i := range funcDataset {
		if funcDataset[i].address == target {
			retInt = i
			break
		}
	}
	return retInt
}

func sim(simFile *os.File) {
	pc := 0
	cycle := 1
	for {
		if pc != -1 {
			readInstruction(asmCode[pc], &pc, &registers, &cycle, simFile)
		} else {
			break
		}
	}
}
