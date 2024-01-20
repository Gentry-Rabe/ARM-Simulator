package main

import (
	"io"
	"os"
	"strconv"
)

func readInstruction(instruction Instruction, pc *int, registers *[32]int, cycle *int, simFile *os.File, dataset *[]Data) {
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
		printSim(instruction, *registers, *cycle, simFile, *dataset)
		*pc += instruction.immediate
		//CBZ
	case 180:
		if registers[instruction.r1] == 0 {
			printSim(instruction, *registers, *cycle, simFile, *dataset)
			*pc += instruction.immediate
		} else {
			printSim(instruction, *registers, *cycle, simFile, *dataset)
			*pc++
		}
		//CBNZ
	case 181:
		if registers[instruction.r1] != 0 {
			printSim(instruction, *registers, *cycle, simFile, *dataset)
			*pc += instruction.immediate
		} else {
			printSim(instruction, *registers, *cycle, simFile, *dataset)
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
		dataPush := [8]int{registers[instruction.destination], 0, 0, 0, 0, 0, 0, 0}
		*dataset = append(*dataset, Data{registers[instruction.r1] + (instruction.immediate * 4), dataPush})
		//LDUR
	case 1986:
		//loop through dataset to find address
		targetAddress := registers[instruction.r1] + (instruction.immediate * 4)
		datasetGrab := *dataset
		for i := range datasetGrab {
			if datasetGrab[i].address == targetAddress {
				registers[instruction.destination] = datasetGrab[i].info[0]
			}
		}
		//BREAK
	case 2038:
		printSim(instruction, *registers, *cycle, simFile, *dataset)
		*pc = -1
	default:
		*pc++
		break
	}
	//this order of execution -> print -> PC will guarantee that the sim will run in the correct order
	if instruction.opCode != 5 && instruction.opCode != 180 && instruction.opCode != 181 && instruction.opCode != 2038 {
		printSim(instruction, *registers, *cycle, simFile, *dataset)
		*pc++
	}
	*cycle++
}

// outputs simulator to file
func printSim(instruction Instruction, registers [32]int, cycle int, simFile *os.File, dataset []Data) {
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
	for i := range dataset {
		_, err = io.WriteString(simFile, strconv.Itoa(dataset[i].address)+":")
		if err != nil {
			return
		}
		for j := 0; j < 8; j++ {
			_, err = io.WriteString(simFile, "\t"+strconv.Itoa(dataset[i].info[j]))
			if err != nil {
				return
			}
		}
		_, err = io.WriteString(simFile, "\n")
		if err != nil {
			return
		}
	}

}

func sim(simFile *os.File) {
	pc := 0
	cycle := 1
	for {
		if pc != -1 {
			readInstruction(asmCode[pc], &pc, &registers, &cycle, simFile, &dataset)
		} else {
			break
		}
	}
}
