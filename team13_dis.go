package main

import (
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
func getOpCode(input string, count int, readingData *bool, asmCode *[]Instruction) string {

	retString := ""

	//6-digit opCode
	//two print lines per case; one for binary and one for opcode + registers
	opcode := binToDec(input[0:6])
	switch {
	case opcode == 5:
		retString = input[0:6] + " " + input[6:32] + "\t" + strconv.Itoa(count) + "\tB\t#" + strconv.Itoa(twosComplementToDec(input[6:32])) + "\n"
		*asmCode = append(*asmCode, Instruction{5, -1, -1, -1, twosComplementToDec(input[6:32]), 0, 0, "B\t#" + strconv.Itoa(twosComplementToDec(input[6:32])), 0})

	default:
		opcode = binToDec(input[0:8])
		switch {
		case opcode == 180:
			retString = input[0:8] + " " + input[8:27] + " " + input[27:32] + "\t" + strconv.Itoa(count) + "\tCBZ\tR" + strconv.Itoa(binToDec(input[27:32])) + ", #" + strconv.Itoa(twosComplementToDec(input[8:27])) + "\n"
			*asmCode = append(*asmCode, Instruction{180, -1, binToDec(input[27:32]), -1, twosComplementToDec(input[8:27]), 0, 0, "CBZ\tR" + strconv.Itoa(binToDec(input[27:32])) + ", #" + strconv.Itoa(twosComplementToDec(input[8:27])), 0})

		case opcode == 181:
			retString = input[0:8] + " " + input[8:27] + " " + input[27:32] + "\t" + strconv.Itoa(count) + "\tCBNZ\tR" + strconv.Itoa(binToDec(input[27:32])) + ", #" + strconv.Itoa(twosComplementToDec(input[8:27])) + "\n"
			*asmCode = append(*asmCode, Instruction{181, -1, binToDec(input[27:32]), -1, twosComplementToDec(input[8:27]), 0, 0, "CBNZ\tR" + strconv.Itoa(binToDec(input[27:32])) + ", #" + strconv.Itoa(twosComplementToDec(input[8:27])), 0})

		default:
			opcode = binToDec(input[0:9])
			switch {
			case opcode == 421:
				retString = input[0:9] + " " + input[9:11] + " " + input[11:27] + " " + input[27:32] + "\t" + strconv.Itoa(count) + "\tMOVZ\tR" + strconv.Itoa(binToDec(input[27:32])) + ", " + strconv.Itoa(binToDec(input[11:27])) + ", LSL " + strconv.Itoa(binToDec(input[9:11])*16) + "\n"
				*asmCode = append(*asmCode, Instruction{421, binToDec(input[27:32]), -1, -1, binToDec(input[11:27]), 0, binToDec(input[9:11]) * 16, "MOVZ\tR" + strconv.Itoa(binToDec(input[27:32])) + ", " + strconv.Itoa(binToDec(input[11:27])) + ", LSL " + strconv.Itoa(binToDec(input[9:11])*16), 0})

			case opcode == 485:
				retString = input[0:9] + " " + input[9:11] + " " + input[11:27] + " " + input[27:32] + "\t" + strconv.Itoa(count) + "\tMOVK\tR" + strconv.Itoa(binToDec(input[27:32])) + ", " + strconv.Itoa(binToDec(input[11:27])) + ", LSL " + strconv.Itoa(binToDec(input[9:11])*16) + "\n"
				*asmCode = append(*asmCode, Instruction{485, binToDec(input[27:32]), -1, -1, binToDec(input[11:27]), 0, binToDec(input[9:11]) * 16, "MOVZ\tR" + strconv.Itoa(binToDec(input[27:32])) + ", " + strconv.Itoa(binToDec(input[11:27])) + ", LSL " + strconv.Itoa(binToDec(input[9:11])*16), 0})

			default:
				opcode = binToDec(input[0:10])
				switch {
				case opcode == 580:
					retString = input[0:10] + " " + input[10:22] + " " + input[22:27] + " " + input[27:32] + "\t" + strconv.Itoa(count) + "\tADDI\tR" + strconv.Itoa(binToDec(input[27:32])) + ", R" + strconv.Itoa(binToDec(input[22:27])) + ", #" + strconv.Itoa(binToDec(input[10:22])) + "\n"
					*asmCode = append(*asmCode, Instruction{580, binToDec(input[27:32]), binToDec(input[22:27]), -1, binToDec(input[10:22]), 0, 0, "ADDI\tR" + strconv.Itoa(binToDec(input[27:32])) + ", R" + strconv.Itoa(binToDec(input[22:27])) + ", #" + strconv.Itoa(binToDec(input[10:22])), 0})

				case opcode == 836:
					retString = input[0:10] + " " + input[10:22] + " " + input[22:27] + " " + input[27:32] + "\t" + strconv.Itoa(count) + "\tSUBI\tR" + strconv.Itoa(binToDec(input[27:32])) + ", R" + strconv.Itoa(binToDec(input[22:27])) + ", #" + strconv.Itoa(binToDec(input[10:22])) + "\n"
					*asmCode = append(*asmCode, Instruction{836, binToDec(input[27:32]), binToDec(input[22:27]), -1, binToDec(input[10:22]), 0, 0, "SUBI\tR" + strconv.Itoa(binToDec(input[27:32])) + ", R" + strconv.Itoa(binToDec(input[22:27])) + ", #" + strconv.Itoa(binToDec(input[10:22])), 0})

				default:
					opcode = binToDec(input[0:11])
					switch {
					case opcode == 1104:
						retString = input[0:11] + " " + input[11:16] + " " + input[16:22] + " " + input[22:27] + " " + input[27:32] + "\t" + strconv.Itoa(count) + "\tAND\tR" + strconv.Itoa(binToDec(input[27:32])) + ", R" + strconv.Itoa(binToDec(input[22:27])) + ", R" + strconv.Itoa(binToDec(input[11:16])) + "\n"
						*asmCode = append(*asmCode, Instruction{1104, binToDec(input[27:32]), binToDec(input[22:27]), binToDec(input[11:16]), 0, 0, 0, "AND\tR" + strconv.Itoa(binToDec(input[27:32])) + ", R" + strconv.Itoa(binToDec(input[22:27])) + ", R" + strconv.Itoa(binToDec(input[11:16])), 0})

					case opcode == 1112:
						retString = input[0:11] + " " + input[11:16] + " " + input[16:22] + " " + input[22:27] + " " + input[27:32] + "\t" + strconv.Itoa(count) + "\tADD\tR" + strconv.Itoa(binToDec(input[27:32])) + ", R" + strconv.Itoa(binToDec(input[22:27])) + ", R" + strconv.Itoa(binToDec(input[11:16])) + "\n"
						*asmCode = append(*asmCode, Instruction{1112, binToDec(input[27:32]), binToDec(input[22:27]), binToDec(input[11:16]), 0, 0, 0, "ADD\tR" + strconv.Itoa(binToDec(input[27:32])) + ", R" + strconv.Itoa(binToDec(input[22:27])) + ", R" + strconv.Itoa(binToDec(input[11:16])), 0})

					case opcode == 1360:
						retString = input[0:11] + " " + input[11:16] + " " + input[16:22] + " " + input[22:27] + " " + input[27:32] + "\t" + strconv.Itoa(count) + "\tORR\tR" + strconv.Itoa(binToDec(input[27:32])) + ", R" + strconv.Itoa(binToDec(input[22:27])) + ", R" + strconv.Itoa(binToDec(input[11:16])) + "\n"
						*asmCode = append(*asmCode, Instruction{1360, binToDec(input[27:32]), binToDec(input[22:27]), binToDec(input[11:16]), 0, 0, 0, "ORR\tR" + strconv.Itoa(binToDec(input[27:32])) + ", R" + strconv.Itoa(binToDec(input[22:27])) + ", R" + strconv.Itoa(binToDec(input[11:16])), 0})

					case opcode == 1624:
						retString = input[0:11] + " " + input[11:16] + " " + input[16:22] + " " + input[22:27] + " " + input[27:32] + "\t" + strconv.Itoa(count) + "\tSUB\tR" + strconv.Itoa(binToDec(input[27:32])) + ", R" + strconv.Itoa(binToDec(input[22:27])) + ", R" + strconv.Itoa(binToDec(input[11:16])) + "\n"
						*asmCode = append(*asmCode, Instruction{1624, binToDec(input[27:32]), binToDec(input[22:27]), binToDec(input[11:16]), 0, 0, 0, "SUB\tR" + strconv.Itoa(binToDec(input[27:32])) + ", R" + strconv.Itoa(binToDec(input[22:27])) + ", R" + strconv.Itoa(binToDec(input[11:16])), 0})

					case opcode == 1690:
						retString = input[0:11] + " " + input[11:16] + " " + input[16:22] + " " + input[22:27] + " " + input[27:32] + "\t" + strconv.Itoa(count) + "\tLSR\tR" + strconv.Itoa(binToDec(input[27:32])) + ", R" + strconv.Itoa(binToDec(input[22:27])) + ", #" + strconv.Itoa(binToDec(input[16:22])) + "\n"
						*asmCode = append(*asmCode, Instruction{1690, binToDec(input[27:32]), binToDec(input[22:27]), -1, binToDec(input[16:22]), 0, 0, "LSR\tR" + strconv.Itoa(binToDec(input[27:32])) + ", R" + strconv.Itoa(binToDec(input[22:27])) + ", #" + strconv.Itoa(binToDec(input[16:22])), 0})

					case opcode == 1691:
						retString = input[0:11] + " " + input[11:16] + " " + input[16:22] + " " + input[22:27] + " " + input[27:32] + "\t" + strconv.Itoa(count) + "\tLSL\tR" + strconv.Itoa(binToDec(input[27:32])) + ", R" + strconv.Itoa(binToDec(input[22:27])) + ", #" + strconv.Itoa(binToDec(input[16:22])) + "\n"
						*asmCode = append(*asmCode, Instruction{1691, binToDec(input[27:32]), binToDec(input[22:27]), -1, binToDec(input[16:22]), 0, 0, "LSL\tR" + strconv.Itoa(binToDec(input[27:32])) + ", R" + strconv.Itoa(binToDec(input[22:27])) + ", #" + strconv.Itoa(binToDec(input[16:22])), 0})

					case opcode == 1984:
						retString = input[0:11] + " " + input[11:20] + " " + input[20:22] + " " + input[22:27] + " " + input[27:32] + "\t" + strconv.Itoa(count) + "\tSTUR\tR" + strconv.Itoa(binToDec(input[27:32])) + ", [R" + strconv.Itoa(binToDec(input[22:27])) + ", #" + strconv.Itoa(binToDec(input[11:20])) + "]\n"
						*asmCode = append(*asmCode, Instruction{1984, binToDec(input[27:32]), binToDec(input[22:27]), -1, binToDec(input[11:20]), 0, 0, "STUR\tR" + strconv.Itoa(binToDec(input[27:32])) + ", [R" + strconv.Itoa(binToDec(input[22:27])) + ", #" + strconv.Itoa(binToDec(input[11:20])) + "]", 0})

					case opcode == 1986:
						retString = input[0:11] + " " + input[11:20] + " " + input[20:22] + " " + input[22:27] + " " + input[27:32] + "\t" + strconv.Itoa(count) + "\tLDUR\tR" + strconv.Itoa(binToDec(input[27:32])) + ", [R" + strconv.Itoa(binToDec(input[22:27])) + ", #" + strconv.Itoa(binToDec(input[11:20])) + "]\n"
						*asmCode = append(*asmCode, Instruction{1986, binToDec(input[27:32]), binToDec(input[22:27]), -1, binToDec(input[11:20]), 0, 0, "LDUR\tR" + strconv.Itoa(binToDec(input[27:32])) + ", [R" + strconv.Itoa(binToDec(input[22:27])) + ", #" + strconv.Itoa(binToDec(input[11:20])) + "]", 0})

					case opcode == 1692:
						retString = input[0:11] + " " + input[11:16] + " " + input[16:22] + " " + input[22:27] + " " + input[27:32] + "\t" + strconv.Itoa(count) + "\tASR\tR" + strconv.Itoa(binToDec(input[27:32])) + ", R" + strconv.Itoa(binToDec(input[22:27])) + ", R" + strconv.Itoa(binToDec(input[11:16])) + "\n"
						*asmCode = append(*asmCode, Instruction{1692, binToDec(input[27:32]), binToDec(input[22:27]), binToDec(input[11:16]), 0, 0, 0, "ASR\tR" + strconv.Itoa(binToDec(input[27:32])) + ", R" + strconv.Itoa(binToDec(input[22:27])) + ", R" + strconv.Itoa(binToDec(input[11:16])), 0})

					case opcode == 1872:
						retString = input[0:11] + " " + input[11:16] + " " + input[16:22] + " " + input[22:27] + " " + input[27:32] + "\t" + strconv.Itoa(count) + "\tEOR\tR" + strconv.Itoa(binToDec(input[27:32])) + ", R" + strconv.Itoa(binToDec(input[22:27])) + ", R" + strconv.Itoa(binToDec(input[11:16])) + "\n"
						*asmCode = append(*asmCode, Instruction{1872, binToDec(input[27:32]), binToDec(input[22:27]), binToDec(input[11:16]), 0, 0, 0, "EOR\tR" + strconv.Itoa(binToDec(input[27:32])) + ", R" + strconv.Itoa(binToDec(input[22:27])) + ", R" + strconv.Itoa(binToDec(input[11:16])), 0})

					case opcode == 2038:
						*readingData = true
						retString = input + "\t" + strconv.Itoa(count) + "\tBREAK\n"
						retString = input[0:8] + " " + input[8:11] + " " + input[11:16] + " " + input[16:21] + " " + input[21:26] + " " + input[26:32] + "\t" + strconv.Itoa(count) + "\tBREAK\n"
						*asmCode = append(*asmCode, Instruction{2038, 0, 0, 0, 0, 0, 0, "BREAK", 0})

					case opcode == 0:
						retString = input + "\t" + strconv.Itoa(count) + "\tNOP\n"
						*asmCode = append(*asmCode, Instruction{0, 0, 0, 0, 0, 0, 0, "NOP", 0})

					default:
						retString = "Unknown Instruction\n"
					}
				}
			}
		}
	}
	return retString
}

func dis(outFile *os.File) {
	//loop through list and output into file
	for i := range txtlines {
		iString := cleanString(txtlines[i])
		if iString != "e" {
			if readingData {
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
}
