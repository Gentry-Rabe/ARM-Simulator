package main

import (
	"io"
	"os"
	"strconv"
)

/*=================================================================================
TODO:
- handle noOP instruction in instruction fetch
- begin cache system
- build system to take ASM as well as binary (not strictly necessary)
=================================================================================*/

// Instruction Fetch pipeline phase. Grab instruction and store into ID queue.
func instructionFetch() {
	if !breakCalled {
		numInstructionsFilled := 0
		line := asmCode[lineCounter]
		for i := 0; i < 4; i++ {
			if checkDataHazards(line.r1) || checkDataHazards(line.r2) {
				break
			}
			if numInstructionsFilled > 1 {
				break
			}
			if preIssue[i].opCode == -1 {
				//store instruction in preIssue
				if line.opCode != 5 && line.opCode != 180 && line.opCode != 181 {
					//if a break instruction is recieved, begin shutdown process
					//this is accomplished by filling pre-issue, pre-alu/mem, and post-alu/mem with blank instructions
					if line.opCode == 2038 || lineCounter == len(asmCode)-1 {
						breakCalled = true
						break
					}
					preIssue[i] = line
					numInstructionsFilled++
					lineCounter++
					line = asmCode[lineCounter]
				} else {
					//note that B instruction uses immediate values and can always be performed immediately
					if line.opCode == 5 {
						numInstructionsFilled = 2
						lineCounter += line.immediate
						break
					}
					//conditional branches
					//CBZ
					if line.opCode == 180 {
						check := true
						for _, element := range postALU {
							if element.destination == line.r1 {
								check = false
							}
						}
						for _, element := range postMEM {
							if element.destination == line.r1 {
								check = false
							}
						}
						if check {
							if registers[line.r1] == 0 {
								lineCounter += line.immediate
							} else {
								lineCounter++
								line = asmCode[lineCounter]
							}
						}
						break
					}
					//CBNZ
					if line.opCode == 181 {
						check := true
						for _, element := range postALU {
							if element.destination == line.r1 {
								check = false
							}
						}
						for _, element := range postMEM {
							if element.destination == line.r1 {
								check = false
							}
						}
						if check {
							if registers[line.r1] != 0 {
								lineCounter += line.immediate
							} else {
								lineCounter++
								line = asmCode[lineCounter]
							}
						}
						break
					}
				}
			}
		}
	}
}

// Instruction issuing pipeline phase. Either branch if B instruction, or sort into EX/MEM phase queue.
func issueUnit() {
	numInstructionsIssued := 0
	for i := 0; i < 4; i++ {
		if numInstructionsIssued > 1 {
			break
		}
		line := preIssue[i]
		if line.opCode == 1984 || line.opCode == 1986 {
			//memory operations, handled by MEM
			if preMEM[0].opCode == -1 {
				preMEM[0] = line
				preIssue[i] = Instruction{-1, -1, -1, -1, 0, 0, 0, "", 0}
				numInstructionsIssued++
			} else if preMEM[1].opCode == -1 {
				preMEM[1] = line
				preIssue[i] = Instruction{-1, -1, -1, -1, 0, 0, 0, "", 0}
				numInstructionsIssued++
			}
		} else if line.opCode != -1 {
			//arithmetic or logical operations, handled by ALU
			if preALU[0].opCode == -1 {
				preALU[0] = line
				preIssue[i] = Instruction{-1, -1, -1, -1, 0, 0, 0, "", 0}
				numInstructionsIssued++
			} else if preALU[1].opCode == -1 {
				preALU[1] = line
				preIssue[i] = Instruction{-1, -1, -1, -1, 0, 0, 0, "", 0}
				numInstructionsIssued++
			}
		}
	}
	preIssue = shiftArray(preIssue)
}

// Execute pipeline phase. Handles all non-MEM operations and pushes to WB queue if needed.
func aluPhase() {
	instruction := preALU[0]
	if instruction.opCode != -1 {
		switch instruction.opCode {
		//ADD
		case 1112:
			instruction.result = registers[instruction.r1] + registers[instruction.r2]
		//SUB
		case 1624:
			instruction.result = registers[instruction.r1] - registers[instruction.r2]
		//ADDI
		case 580:
			instruction.result = registers[instruction.r1] + instruction.immediate
		//SUBI
		case 836:
			instruction.result = registers[instruction.r1] - instruction.immediate
		//AND
		case 1104:
			instruction.result = registers[instruction.r1] & registers[instruction.r2]
		//ORR
		case 1360:
			instruction.result = registers[instruction.r1] | registers[instruction.r2]
		//EOR
		case 1872:
			instruction.result = registers[instruction.r1] ^ registers[instruction.r2]
		//LSR
		case 1690:
			instruction.result = registers[instruction.r1] >> instruction.immediate
		//LSL
		case 1691:
			instruction.result = registers[instruction.r1] << instruction.immediate
		//ASR
		case 1692:
			instruction.result = registers[instruction.r1] >> registers[instruction.r2]
		}
		postALU[0] = instruction
		preALU[0] = Instruction{-1, -1, -1, -1, 0, 0, 0, "", 0}
		preALU = shiftArray(preALU)
	}
}

// Memory pipeline phase. Handles all MEM operations and pushes to WB queue if needed.
func memPhase() {
	instruction := preMEM[0]
	if instruction.opCode != -1 {
		switch instruction.opCode {
		//STUR
		case 1984:
			if findIndex(registers[instruction.r2]+(instruction.immediate*4), pipelineDataset) != -1 {
				newData := pipelineDataset
				newData[findIndex(registers[instruction.r2]+(instruction.immediate*4), pipelineDataset)].value = registers[instruction.r1]
				pipelineDataset = newData
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
					pipelineDataset = append(pipelineDataset, Data{noOffsetAddress + (i * 4), 0})
				}
				newData := pipelineDataset
				newData[findIndex(registers[instruction.r2]+(instruction.immediate*4), pipelineDataset)].value = registers[instruction.r1]
				pipelineDataset = newData
			}
			//LDUR
		case 1986:
			newData := pipelineDataset
			indexForLoad := findIndex(registers[instruction.r1]+(instruction.immediate*4), pipelineDataset)
			if indexForLoad != -1 {
				instruction.result = newData[indexForLoad].value
			} else {
				instruction.result = 0
			}
		}
		postMEM[0] = instruction
		preMEM[0] = Instruction{-1, -1, -1, -1, 0, 0, 0, "", 0}
		preMEM = shiftArray(preMEM)
	}
}

// writes next in post-ALU/post-MEM queue into registers
func writeBack() {
	if postALU[0].opCode != -1 {
		registers[postALU[0].destination] = postALU[0].result
		postALU[0] = Instruction{-1, -1, -1, -1, 0, 0, 0, "", 0}
		postALU = shiftArray(postALU)
	}
	if postMEM[0].opCode != -1 && postMEM[0].opCode != 1984 {
		registers[postMEM[0].destination] = postMEM[0].result
		postMEM[0] = Instruction{-1, -1, -1, -1, 0, 0, 0, "", 0}
		postMEM = shiftArray(postMEM)
	} else if postMEM[0].opCode == 1984 {
		postMEM[0] = Instruction{-1, -1, -1, -1, 0, 0, 0, "", 0}
		postMEM = shiftArray(postMEM)
	}
}

// prints the current state of the pipeline as well as registers, (eventually) cache, and data
func printPipeline(pipelineFile *os.File) {
	_, err := io.WriteString(pipelineFile, "--------------------\n")
	if err != nil {
		return
	}
	_, err = io.WriteString(pipelineFile, "Cycle:"+strconv.Itoa(cycle)+"\n\n")
	if err != nil {
		return
	}
	//=======Issue buffer=======
	_, err = io.WriteString(pipelineFile, "Pre-Issue Buffer:\n")
	if err != nil {
		return
	}
	for i := 0; i < 4; i++ {
		_, err := io.WriteString(pipelineFile, "\tEntry "+strconv.Itoa(i)+": "+preIssue[i].opString+"\n")
		if err != nil {
			return
		}
	}

	//=======ALU=======
	_, err = io.WriteString(pipelineFile, "Pre-ALU Queue:\n")
	if err != nil {
		return
	}
	_, err = io.WriteString(pipelineFile, "\tEntry 0: "+preALU[0].opString+"\n")
	if err != nil {
		return
	}
	_, err = io.WriteString(pipelineFile, "\tEntry 1: "+preALU[1].opString+"\n")
	if err != nil {
		return
	}
	_, err = io.WriteString(pipelineFile, "Post-ALU Queue:\n")
	if err != nil {
		return
	}
	_, err = io.WriteString(pipelineFile, "\tEntry 0: "+postALU[0].opString+"\n")
	if err != nil {
		return
	}

	//=======MEM=======
	_, err = io.WriteString(pipelineFile, "Pre-MEM Queue:\n")
	if err != nil {
		return
	}
	_, err = io.WriteString(pipelineFile, "\tEntry 0: "+preMEM[0].opString+"\n")
	if err != nil {
		return
	}
	_, err = io.WriteString(pipelineFile, "\tEntry 1: "+preMEM[1].opString+"\n")
	if err != nil {
		return
	}
	_, err = io.WriteString(pipelineFile, "Post-MEM Queue:\n")
	if err != nil {
		return
	}
	_, err = io.WriteString(pipelineFile, "\tEntry 0: "+postMEM[0].opString+"\n\n")
	if err != nil {
		return
	}

	//=======Registers=======
	_, err = io.WriteString(pipelineFile, "Registers:\n")
	if err != nil {
		return
	}
	for i := 1; i <= 32; i++ {
		if (i-1)%8 == 0 {
			_, err = io.WriteString(pipelineFile, "R"+strconv.Itoa(i-1)+":")
			if err != nil {
				return
			}
		}
		_, err = io.WriteString(pipelineFile, "\t"+strconv.Itoa(registers[i-1]))
		if err != nil {
			return
		}
		if i%8 == 0 {
			_, err = io.WriteString(pipelineFile, "\n")
			if err != nil {
				return
			}
		}
	}
	_, err = io.WriteString(pipelineFile, "\n")
	if err != nil {
		return
	}
	/*
	   Cache: to be written


	   Cache
	   Set 0: LRU=0
	   Entry 0:[(0,0,0)<0,0>]
	   Entry 1:[(0,0,0)<0,0>]
	   Set 1: LRU=0
	   Entry 0:[(0,0,0)<0,0>]
	   Entry 1:[(0,0,0)<0,0>]
	   Set 2: LRU=0
	   Entry 0:[(0,0,0)<0,0>]
	   Entry 1:[(0,0,0)<0,0>]
	   Set 3: LRU=0
	   Entry 0:[(0,0,0)<0,0>]
	   Entry 1:[(0,0,0)<0,0>]
	*/
	//=======Data=======
	_, err = io.WriteString(pipelineFile, "Data:\n")
	if err != nil {
		return
	}
	for i := range pipelineDataset {
		if i%8 == 0 {
			_, err = io.WriteString(pipelineFile, strconv.Itoa(pipelineDataset[i].address)+": ")
		}
		_, err = io.WriteString(pipelineFile, "\t"+strconv.Itoa(pipelineDataset[i].value))
		if i%8 == 7 {
			_, err = io.WriteString(pipelineFile, "\n")
		}
		if err != nil {
			return
		}
	}
}

// checks to see if a register is currently being used by another function, in order to prevent data hazards
func checkDataHazards(register int) bool {
	ret := false
	for _, element := range preIssue {
		if register == element.destination {
			ret = true
		}
	}
	for _, element := range preALU {
		if register == element.destination {
			ret = true
		}
	}
	for _, element := range preMEM {
		if register == element.destination {
			ret = true
		}
	}
	if register == -1 {
		ret = false
	}
	//no need to check post-ALU and post-MEM, they can be written back and accessed in the same cycle
	return ret
}

// shifts data in arrays up to maintain FIFO procedure - length: 4 arrays only
func shiftArray(array []Instruction) []Instruction {
	nArray := make([]Instruction, 0)
	valueIndex := 0
	for _, element := range array {
		if element.opCode != -1 {
			nArray = append(nArray, element)
		} else {
			valueIndex++
		}
	}
	for i := 0; i < valueIndex; i++ {
		nArray = append(nArray, Instruction{-1, -1, -1, -1, 0, 0, 0, "", 0})
	}
	return nArray
}

// global pipeline variables
var running = true
var breakCalled = false
var cycle = 1
var lineCounter = 0
var pipelineDataset = make([]Data, 0)

// pipeline queues
var preIssue = []Instruction{
	{-1, -1, -1, -1, 0, 0, 0, "", 0},
	{-1, -1, -1, -1, 0, 0, 0, "", 0},
	{-1, -1, -1, -1, 0, 0, 0, "", 0},
	{-1, -1, -1, -1, 0, 0, 0, "", 0},
}
var preALU = []Instruction{
	{-1, -1, -1, -1, 0, 0, 0, "", 0},
	{-1, -1, -1, -1, 0, 0, 0, "", 0},
}
var preMEM = []Instruction{
	{-1, -1, -1, -1, 0, 0, 0, "", 0},
	{-1, -1, -1, -1, 0, 0, 0, "", 0},
}
var postALU = []Instruction{
	{-1, -1, -1, -1, 0, 0, 0, "", 0},
}
var postMEM = []Instruction{
	{-1, -1, -1, -1, 0, 0, 0, "", 0},
}

// Calls components of pipeline in REVERSE ORDER. This is done to prevent data hazards.
func pipeline(pipelineFile *os.File) {
	//reset registers and data for accurate simulation
	registers = [32]int{
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
	}
	readStart = true
	//running code for pipeline simulator. loops until end of code, then breaks
	for {
		if running {
			writeBack()
			memPhase()
			aluPhase()
			issueUnit()
			instructionFetch()
			printPipeline(pipelineFile)
			cycle++
			if breakCalled && preIssue[0].opCode == -1 && postALU[0].opCode == -1 && postMEM[0].opCode == -1 {
				running = false
			}
		} else {
			break
		}
	}
	printPipeline(pipelineFile)
}
