package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

// error check function

// object for holding all parameters for all instructions
type instruction struct {
	binary     string
	opcode     string
	parameters []string
}

// converts a binary string to a decimal string
func binaryToDecimal(binaryNumber string) string {
	var registerNumber float64 = 0
	for i := 0; i < len(binaryNumber); i++ {
		if string(binaryNumber[i]) == "1" {
			registerNumber = registerNumber + math.Pow(2, float64(len(binaryNumber)-(i+1)))
		}
	}
	registerNumberString := strconv.FormatFloat(registerNumber, 'f', 0, 64)
	return registerNumberString
}

// swaps the bits in the binary string, converts to decimal, adds one, returns number with negative sign

func twosComplement(binaryNumber string) string {
	if binaryNumber[0:1] == "1" {
		var newBinaryString = ""
		for i := 0; i < len(binaryNumber); i++ {
			if binaryNumber[i] == '0' {
				newBinaryString += "1"
			} else if binaryNumber[i] == '1' {
				newBinaryString += "0"
			}
		}
		//fmt.Println(newBinaryString)
		x, err := strconv.Atoi(binaryToDecimal(newBinaryString))
		checkError(err)
		x++
		var newBinaryString2 = "-" + strconv.Itoa(x)
		//fmt.Println(newBinaryString2)
		return newBinaryString2
	} else {
		return binaryToDecimal(binaryNumber)
	}

}

// find opcode given line input
func findOpcode(currentLine string, instructionsArray []instruction) []instruction {
	var instructionName string
	if currentLine[0:6] == "000101" {
		//B format instruction
		instructionName = "B"
		instructionsArray = bFormatdisassembler(instructionName, currentLine, instructionsArray)
	} else if currentLine[0:8] == "10110100" {
		//CBZ instruction
		instructionName = "CBZ"
		instructionsArray = cbFormatdisassembler(instructionName, currentLine, instructionsArray)
	} else if currentLine[0:8] == "10110101" {
		//CBNZ instruction
		instructionName = "CBNZ"
		instructionsArray = cbFormatdisassembler(instructionName, currentLine, instructionsArray)
	} else if currentLine[0:9] == "110100101" {
		//MOVZ instruction
		instructionName = "MOVZ"
		instructionsArray = imFormatdisassembler(instructionName, currentLine, instructionsArray)
	} else if currentLine[0:9] == "111100101" {
		//MOVK instruction
		instructionName = "MOVK"
		instructionsArray = imFormatdisassembler(instructionName, currentLine, instructionsArray)
	} else if currentLine[0:10] == "1001000100" {
		//ADDI instruction
		instructionName = "ADDI"
		instructionsArray = iFormatdisassembler(instructionName, currentLine, instructionsArray)
	} else if currentLine[0:10] == "1101000100" {
		//SUBI instruction
		instructionName = "SUBI"
		instructionsArray = iFormatdisassembler(instructionName, currentLine, instructionsArray)
	} else if currentLine[0:11] == "11111110110" {
		//BREAK instruction
		instructionName = "BREAK"
		instructionsArray = breakDisassembler(instructionName, currentLine, instructionsArray)
	}
	switch currentLine[0:11] {
	case "10001010000":
		//AND
		instructionName = "AND"
		instructionsArray = rFormatdisassembler(instructionName, currentLine, instructionsArray)
	case "10001011000":
		//ADD
		instructionName = "ADD"
		instructionsArray = rFormatdisassembler(instructionName, currentLine, instructionsArray)
	case "10101010000":
		//ORR
		instructionName = "ORR"
		instructionsArray = rFormatdisassembler(instructionName, currentLine, instructionsArray)
	case "11001011000":
		//SUB
		instructionName = "SUB"
		instructionsArray = rFormatdisassembler(instructionName, currentLine, instructionsArray)
	case "11010011010":
		//LSR
		instructionName = "LSR"
		instructionsArray = rFormatdisassembler(instructionName, currentLine, instructionsArray)
	case "11010011011":
		//LSL
		instructionName = "LSL"
		instructionsArray = rFormatdisassembler(instructionName, currentLine, instructionsArray)
	case "11111000000":
		//STUR
		instructionName = "STUR"
		instructionsArray = dFormatdisassembler(instructionName, currentLine, instructionsArray)
	case "11111000010":
		//LDUR
		instructionName = "LDUR"
		instructionsArray = dFormatdisassembler(instructionName, currentLine, instructionsArray)
	case "11010011100":
		//ASR
		instructionName = "ASR"
		instructionsArray = rFormatdisassembler(instructionName, currentLine, instructionsArray)
	case "11101010000":
		//EOR
		instructionName = "EOR"
		instructionsArray = rFormatdisassembler(instructionName, currentLine, instructionsArray)
	case "00000000000":
		//NOP
		instructionName = "NOP"
		var i = instruction{binary: currentLine, opcode: instructionName}
		instructionsArray = append(instructionsArray, i)
	}
	return instructionsArray
}

// disassemble line input into specific formats based on opcode
func rFormatdisassembler(instructionName string, currentLine string, instructionArray []instruction) []instruction {
	var binaryString = currentLine[0:11] + " " + currentLine[11:16] + " " + currentLine[16:22] + " " + currentLine[22:27] + " " + currentLine[27:32]

	if instructionName == "LSL" || instructionName == "ASR" || instructionName == "LSR" {
		var i = instruction{binary: binaryString, opcode: instructionName, parameters: []string{"R" + binaryToDecimal(currentLine[27:32]), "R" + binaryToDecimal(currentLine[22:27]), "#" + binaryToDecimal(currentLine[16:22])}}
		instructionArray = append(instructionArray, i)

	} else {
		var i = instruction{binary: binaryString, opcode: instructionName, parameters: []string{"R" + binaryToDecimal(currentLine[27:32]), "R" + binaryToDecimal(currentLine[22:27]), "R" + binaryToDecimal(currentLine[11:16])}}
		instructionArray = append(instructionArray, i)
	}
	return instructionArray
}
func dFormatdisassembler(instructionName string, currentLine string, instructionArray []instruction) []instruction {
	var binaryString = currentLine[0:11] + " " + currentLine[11:20] + " " + currentLine[20:22] + " " + currentLine[22:27] + " " + currentLine[27:32]
	var i = instruction{binary: binaryString, opcode: instructionName, parameters: []string{"R" + binaryToDecimal(currentLine[27:32]), "[R" + binaryToDecimal(currentLine[22:27]), "#" + binaryToDecimal(currentLine[11:20]) + "]"}}
	instructionArray = append(instructionArray, i)
	return instructionArray
}
func iFormatdisassembler(instructionName string, currentLine string, instructionArray []instruction) []instruction {
	var binaryString = currentLine[0:10] + " " + currentLine[10:22] + " " + currentLine[22:27] + " " + currentLine[27:32]
	var i = instruction{binary: binaryString, opcode: instructionName, parameters: []string{"R" + binaryToDecimal(currentLine[27:32]), "R" + binaryToDecimal(currentLine[22:27]), "#" + binaryToDecimal(currentLine[10:22])}}
	instructionArray = append(instructionArray, i)
	return instructionArray
}
func bFormatdisassembler(instructionName string, currentLine string, instructionArray []instruction) []instruction {
	var binaryString = currentLine[0:6] + " " + currentLine[6:32]
	if currentLine[7] == '1' {
		var i = instruction{binary: binaryString, opcode: instructionName, parameters: []string{"#" + twosComplement(currentLine[6:32])}}
		instructionArray = append(instructionArray, i)
	} else {
		var i = instruction{binary: binaryString, opcode: instructionName, parameters: []string{"#" + binaryToDecimal(currentLine[6:32])}}
		instructionArray = append(instructionArray, i)
	}
	return instructionArray
}
func cbFormatdisassembler(instructionName string, currentLine string, instructionArray []instruction) []instruction {
	var binaryString = currentLine[0:8] + " " + currentLine[8:27] + " " + currentLine[27:32]
	if currentLine[9] == '1' {
		var i = instruction{binary: binaryString, opcode: instructionName, parameters: []string{"R" + binaryToDecimal(currentLine[27:32]), "#" + twosComplement(currentLine[8:27])}}
		instructionArray = append(instructionArray, i)
	} else {
		var i = instruction{binary: binaryString, opcode: instructionName, parameters: []string{"R" + binaryToDecimal(currentLine[27:32]), "#" + binaryToDecimal(currentLine[8:27])}}
		instructionArray = append(instructionArray, i)
	}
	return instructionArray
}
func imFormatdisassembler(instructionName string, currentLine string, instructionArray []instruction) []instruction {
	var binaryString = currentLine[0:9] + " " + currentLine[9:11] + " " + currentLine[11:27] + " " + currentLine[27:32]
	var shiftCode = currentLine[9:11]
	if shiftCode == "00" {
		shiftCode = "0"
	} else if shiftCode == "11" {
		shiftCode = "48"
	}
	var i = instruction{binary: binaryString, opcode: instructionName, parameters: []string{"R" + binaryToDecimal(currentLine[27:32]), binaryToDecimal(currentLine[11:27]), "LSL", shiftCode}}
	instructionArray = append(instructionArray, i)
	return instructionArray
}
func breakDisassembler(instructionName string, currentLine string, instructionArray []instruction) []instruction {
	var binaryString = currentLine[0:1] + " " + currentLine[1:6] + " " + currentLine[6:11] + " " + currentLine[11:16] + " " + currentLine[16:21] + " " + currentLine[21:26] + " " + currentLine[26:32]
	var i = instruction{binary: binaryString, opcode: instructionName, parameters: []string{}}
	instructionArray = append(instructionArray, i)
	return instructionArray
}

func checkError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func checkNumDigits(currentLine string, registerIndex int) int {
	if registerIndex+1 == len(currentLine) {
		x, err := strconv.Atoi(currentLine[registerIndex : registerIndex+1])
		checkError(err)
		return x
	} else if currentLine[registerIndex+1] == ',' || currentLine[registerIndex+1] == ' ' || currentLine[registerIndex+1] == ']' {
		x, err := strconv.Atoi(currentLine[registerIndex : registerIndex+1])
		checkError(err)
		return x
	} else {
		x, err := strconv.Atoi(currentLine[registerIndex : registerIndex+2])
		checkError(err)
		return x
	}
}

func findStoreRegister(currentLine string, opcode string) int {
	var registerIndex = strings.Index(currentLine, opcode) + len(opcode) + 2
	return checkNumDigits(currentLine, registerIndex)
}

func findRnRegister(currentLine string, opcode string) int {
	var registerIndex int
	if strings.Contains(currentLine, "[") {
		registerIndex = strings.Index(currentLine, opcode) + len(opcode) + len(strconv.Itoa(storeReg)) + 6
	} else {
		registerIndex = strings.Index(currentLine, opcode) + len(opcode) + len(strconv.Itoa(storeReg)) + 5
	}
	return checkNumDigits(currentLine, registerIndex)
}

func findRmRegister(currentLine string, opcode string) int {
	var registerIndex int
	registerIndex = strings.Index(currentLine, opcode) + len(opcode) + len(strconv.Itoa(storeReg)) + len(strconv.Itoa(rnReg)) + 8
	return checkNumDigits(currentLine, registerIndex)
}

func findImmediateValue(currentLine string) int {
	var hashtagIndex = strings.Index(currentLine, "#") + 1
	var immediateValue string
	if strings.Contains(currentLine, "[") {
		immediateValue = currentLine[hashtagIndex : len(currentLine)-1]
	} else {
		immediateValue = currentLine[hashtagIndex:len(currentLine)]
	}
	x, err := strconv.Atoi(immediateValue)
	checkError(err)
	return x
}

var registerArray [32]int

var cycleNum = 1

var dataMemoryStart = 0

// dataMemoryEnd int
var memoryLocationSim = 96
var storeReg int
var rnReg int
var rmReg int

var immedReg int
var condReg int
var shamtReg int
var breakFlag = false

var lineNumber = 0
var lineArray []string
var dataArray []string
var breakIndex int

var address int
var breakMemoryCounter int

// project 1 vars
var instructionsArray []instruction
var breakCheck = false
var memoryLocation = 96

func main() {
	//cmd line parameters
	var InputFileName *string
	var OutputFileName *string

	InputFileName = flag.String("i", "", "Gets the input file name")
	OutputFileName = flag.String("o", "", "Gets the output file name")
	flag.Parse()
	//testFile, err := os.Open("test.txt")
	//checkError(err)
	//defer testFile.Close()

	inputFile, err := os.Open(*InputFileName)
	checkError(err)
	defer inputFile.Close()

	outputFileDis, err := os.Create(*OutputFileName + "_dis.txt")
	checkError(err)
	defer outputFileDis.Close()

	outputFileSim, err := os.Create(*OutputFileName + "_sim.txt")
	checkError(err)
	defer outputFileSim.Close()

	//fmt.Println("Input:", *InputFileName)
	//fmt.Println("Output dis:", outputFileDis)
	//fmt.Println("Output sim:", outputFileSim)

	if flag.NArg() != 0 {
		os.Exit(200)
	}

	//PROJECT 1 START
	//read input file line by line (check for break being last opcode)
	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		currentLine := scanner.Text()
		if !breakCheck {
			instructionsArray = findOpcode(currentLine, instructionsArray)
		} else {
			var newInstruction = instruction{binary: currentLine, opcode: twosComplement(currentLine)}
			instructionsArray = append(instructionsArray, newInstruction)
		}
		if instructionsArray[len(instructionsArray)-1].opcode == "BREAK" {
			breakCheck = true
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(instructionsArray); i++ {

		_, err := io.WriteString(outputFileDis, instructionsArray[i].binary+"\t"+strconv.Itoa(memoryLocationSim)+"\t"+instructionsArray[i].opcode+"\t")
		checkError(err)
		memoryLocationSim += 4

		for k := 0; k < len(instructionsArray[i].parameters); k++ {
			if k == len(instructionsArray[i].parameters)-1 {
				_, err := io.WriteString(outputFileDis, instructionsArray[i].parameters[k])
				checkError(err)
			} else {
				_, err := outputFileDis.WriteString(instructionsArray[i].parameters[k] + ", ")
				checkError(err)
			}
		}
		_, err2 := io.WriteString(outputFileDis, "\n")
		checkError(err2)
	}

	//PROJECT 2 START
	outputFileDis2, err := os.Open(outputFileDis.Name())
	checkError(err)
	defer outputFileDis2.Close()

	for breakFlag != true {

		scanner2 := bufio.NewScanner(outputFileDis2) //here
		for scanner2.Scan() {
			breakLine := scanner2.Text()
			lineArray = append(lineArray, breakLine)
		}

		for k := 0; k < len(lineArray); k++ {
			if strings.Contains(lineArray[k], "BREAK") {
				breakIndex = k
			}
		}
		dataMemoryStart = (breakIndex+1)*4 + 96

		for v := breakIndex + 1; v < len(lineArray); v++ {
			var tabIndex = 34
			var memoryString string
			//fmt.Println(len(lineArray[v]))
			for p := 0; p < len(lineArray[v]); p++ {
				var currIndex = lineArray[v][(tabIndex + p - 1):(tabIndex + p)]
				if currIndex != "0" && currIndex != "1" && currIndex != "2" && currIndex != "3" && currIndex != "4" && currIndex != "5" && currIndex != "6" && currIndex != "7" && currIndex != "8" && currIndex != "9" {
					break
				} else {
					memoryString = lineArray[v][tabIndex-1 : tabIndex+p]
				}
			}
			var memorylength = len(memoryString)
			var dataIndex = tabIndex + memorylength
			dataArray = append(dataArray, lineArray[v][dataIndex:(len(lineArray[v])-1)])
			breakMemoryCounter++
		}

		for i := 0; i < len(lineArray); i++ {
			_, err := io.WriteString(outputFileSim, "===================="+"\n"+"cycle:"+strconv.Itoa(cycleNum)+"\t"+strconv.Itoa(memoryLocation)+"\t")
			checkError(err)
			currentLine := lineArray[i]

			lineNumber++
			if strings.Contains(currentLine, "BREAK") {
				breakFlag = true
				_, err := io.WriteString(outputFileSim, "BREAK")
				checkError(err)

			} else if strings.Contains(currentLine, "ADDI") { // done
				storeReg = findStoreRegister(currentLine, "ADDI")
				rnReg = findRnRegister(currentLine, "ADDI")
				immedReg = findImmediateValue(currentLine)

				registerArray[storeReg] = registerArray[rnReg] + immedReg
				_, err := io.WriteString(outputFileSim, "ADDI"+"\t"+"R"+strconv.Itoa(storeReg)+",\t"+"R"+strconv.Itoa(rnReg)+",\t"+"#"+strconv.Itoa(immedReg))
				checkError(err)

			} else if strings.Contains(currentLine, "SUBI") {
				storeReg = findStoreRegister(currentLine, "SUBI")
				rnReg = findRnRegister(currentLine, "SUBI")
				immedReg = findImmediateValue(currentLine)

				registerArray[storeReg] = registerArray[rnReg] - immedReg
				_, err := io.WriteString(outputFileSim, "SUBI"+"\t"+"R"+strconv.Itoa(storeReg)+",\t"+"R"+strconv.Itoa(rnReg)+",\t"+"#"+strconv.Itoa(immedReg))
				checkError(err)

			} else if strings.Contains(currentLine, "CBNZ") {
				condReg = findStoreRegister(currentLine, "CBZ")
				immedReg = findImmediateValue(currentLine)
				if registerArray[condReg] != 0 {
					i = i + immedReg - 1
					memoryLocation = memoryLocation + (4 * immedReg) - 4
				}
				_, err := io.WriteString(outputFileSim, "CBNZ"+"\t"+"R"+strconv.Itoa(condReg)+",\t"+"#"+strconv.Itoa(immedReg))
				checkError(err)

			} else if strings.Contains(currentLine, "CBZ") {
				condReg = findStoreRegister(currentLine, "CBZ")
				immedReg = findImmediateValue(currentLine)
				if registerArray[condReg] == 0 {
					i = i + immedReg - 1
					memoryLocation = memoryLocation + (4 * immedReg) - 4
				}

				_, err := io.WriteString(outputFileSim, "CBZ"+"\t"+"R"+strconv.Itoa(condReg)+",\t"+"#"+strconv.Itoa(immedReg))
				checkError(err)

			} else if strings.Contains(currentLine, "ASR") {
				storeReg = findStoreRegister(currentLine, "ASR")
				rnReg = findRnRegister(currentLine, "ASR")
				shamtReg = findImmediateValue(currentLine)

				registerArray[storeReg] = registerArray[rnReg] >> shamtReg
				_, err := io.WriteString(outputFileSim, "ASR"+"\t"+"R"+strconv.Itoa(storeReg)+",\t"+"R"+strconv.Itoa(rnReg)+",\t#"+strconv.Itoa(shamtReg))
				checkError(err)

			} else if strings.Contains(currentLine, "LSL") {
				storeReg = findStoreRegister(currentLine, "LSL")
				rnReg = findRnRegister(currentLine, "LSL")
				shamtReg = findImmediateValue(currentLine)

				registerArray[storeReg] = registerArray[rnReg] << shamtReg
				_, err := io.WriteString(outputFileSim, "LSL"+"\t"+"R"+strconv.Itoa(storeReg)+",\t"+"R"+strconv.Itoa(rnReg)+",\t#"+strconv.Itoa(shamtReg))
				checkError(err)

			} else if strings.Contains(currentLine, "LSR") {
				storeReg = findStoreRegister(currentLine, "LSR")
				rnReg = findRnRegister(currentLine, "LSR")
				shamtReg = findImmediateValue(currentLine)

				var unsignedReg1 uint32 = uint32(registerArray[rnReg])

				var shiftedReg uint32 = uint32(unsignedReg1) >> shamtReg
				shiftedRegInt := int(shiftedReg)

				registerArray[storeReg] = shiftedRegInt
				_, err := io.WriteString(outputFileSim, "LSR"+"\t"+"R"+strconv.Itoa(storeReg)+",\t"+"R"+strconv.Itoa(rnReg)+",\t#"+strconv.Itoa(shamtReg))
				checkError(err)

			} else if strings.Contains(currentLine, "SUB") {
				storeReg = findStoreRegister(currentLine, "SUB")
				rnReg = findRnRegister(currentLine, "SUB")
				rmReg = findRmRegister(currentLine, "SUB")

				registerArray[storeReg] = registerArray[rnReg] - registerArray[rmReg]
				_, err := io.WriteString(outputFileSim, "SUB"+"\t"+"R"+strconv.Itoa(storeReg)+",\t"+"R"+strconv.Itoa(rnReg)+",\tR"+strconv.Itoa(rmReg))
				checkError(err)

			} else if strings.Contains(currentLine, "ADD") {
				storeReg = findStoreRegister(currentLine, "ADD")
				rnReg = findRnRegister(currentLine, "ADD")
				rmReg = findRmRegister(currentLine, "ADD")

				registerArray[storeReg] = registerArray[rmReg] + registerArray[rnReg]
				_, err := io.WriteString(outputFileSim, "ADD"+"\t"+"R"+strconv.Itoa(storeReg)+",\t"+"R"+strconv.Itoa(rnReg)+",\t"+"R"+strconv.Itoa(rmReg))
				checkError(err)

			} else if strings.Contains(currentLine, "B") {
				immedReg = findImmediateValue(currentLine)
				if immedReg != 0 {
					i = i + immedReg - 1
				}
				_, err := io.WriteString(outputFileSim, "B"+"\t"+"#"+strconv.Itoa(immedReg))
				memoryLocation = memoryLocation + (4 * immedReg) - 4
				checkError(err)
			} else if strings.Contains(currentLine, "AND") {
				storeReg = findStoreRegister(currentLine, "AND")
				rnReg = findRnRegister(currentLine, "AND")
				rmReg = findRmRegister(currentLine, "AND")
				registerArray[storeReg] = registerArray[rmReg] & registerArray[rnReg]
				_, err := io.WriteString(outputFileSim, "AND"+"\t"+"R"+strconv.Itoa(storeReg)+",\t"+"R"+strconv.Itoa(rnReg)+",\t"+"R"+strconv.Itoa(rmReg))
				checkError(err)

			} else if strings.Contains(currentLine, "ORR") {
				storeReg = findStoreRegister(currentLine, "ORR")
				rnReg = findRnRegister(currentLine, "ORR")
				rmReg = findRmRegister(currentLine, "ORR")
				registerArray[storeReg] = registerArray[rmReg] | registerArray[rnReg]
				_, err := io.WriteString(outputFileSim, "ORR"+"\t"+"R"+strconv.Itoa(storeReg)+",\t"+"R"+strconv.Itoa(rnReg)+",\t"+"R"+strconv.Itoa(rmReg))
				checkError(err)

			} else if strings.Contains(currentLine, "EOR") {
				storeReg = findStoreRegister(currentLine, "EOR")
				rnReg = findRnRegister(currentLine, "EOR")
				rmReg = findRmRegister(currentLine, "EOR")
				registerArray[storeReg] = registerArray[rmReg] ^ registerArray[rnReg]
				_, err := io.WriteString(outputFileSim, "EOR"+"\t"+"R"+strconv.Itoa(storeReg)+",\t"+"R"+strconv.Itoa(rnReg)+",\t"+"R"+strconv.Itoa(rmReg))
				checkError(err)
			} else if strings.Contains(currentLine, "LDUR") {
				storeReg = findStoreRegister(currentLine, "LDUR")
				rnReg = findRnRegister(currentLine, "LDUR")
				immedReg = findImmediateValue(currentLine)

				address = registerArray[rnReg] + (immedReg * 4)
				for address > dataMemoryStart+4*(len(dataArray)-1) {
					dataArray = append(dataArray, "0")
				}
				registerArray[storeReg], err = strconv.Atoi(dataArray[(address-dataMemoryStart)/4])

				_, err := io.WriteString(outputFileSim, "LDUR"+"\t"+"R"+strconv.Itoa(storeReg)+",\t"+"[R"+strconv.Itoa(rnReg)+",\t"+"#"+strconv.Itoa(immedReg)+"]")
				checkError(err)

			} else if strings.Contains(currentLine, "STUR") {

				storeReg = findStoreRegister(currentLine, "STUR")
				rnReg = findRnRegister(currentLine, "STUR")
				immedReg = findImmediateValue(currentLine)
				address = registerArray[rnReg] + (immedReg * 4)

				for address > dataMemoryStart+4*(len(dataArray)-1) {
					dataArray = append(dataArray, "0")
				}
				dataArray[(address-dataMemoryStart)/4] = strconv.Itoa(registerArray[storeReg])

				_, err := io.WriteString(outputFileSim, "STUR"+"\t"+"R"+strconv.Itoa(storeReg)+"\t"+"[R"+strconv.Itoa(rnReg)+",\t#"+strconv.Itoa(immedReg)+"]")
				checkError(err)

			}
			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}
			_, err2 := io.WriteString(outputFileSim, "\n"+"\n")
			checkError(err2)

			_, err3 := io.WriteString(outputFileSim, "registers:"+"\n")
			checkError(err3)

			_, err4 := io.WriteString(outputFileSim, "r00:"+"\t"+strconv.Itoa(registerArray[0])+"\t"+strconv.Itoa(registerArray[1])+"\t"+strconv.Itoa(registerArray[2])+"\t"+strconv.Itoa(registerArray[3])+"\t"+strconv.Itoa(registerArray[4])+"\t"+strconv.Itoa(registerArray[5])+"\t"+strconv.Itoa(registerArray[6])+"\t"+strconv.Itoa(registerArray[7])+"\n")
			checkError(err4)
			_, err5 := io.WriteString(outputFileSim, "r08:"+"\t"+strconv.Itoa(registerArray[8])+"\t"+strconv.Itoa(registerArray[9])+"\t"+strconv.Itoa(registerArray[10])+"\t"+strconv.Itoa(registerArray[11])+"\t"+strconv.Itoa(registerArray[12])+"\t"+strconv.Itoa(registerArray[13])+"\t"+strconv.Itoa(registerArray[14])+"\t"+strconv.Itoa(registerArray[15])+"\n")
			checkError(err5)
			_, err6 := io.WriteString(outputFileSim, "r16:"+"\t"+strconv.Itoa(registerArray[16])+"\t"+strconv.Itoa(registerArray[17])+"\t"+strconv.Itoa(registerArray[18])+"\t"+strconv.Itoa(registerArray[19])+"\t"+strconv.Itoa(registerArray[20])+"\t"+strconv.Itoa(registerArray[21])+"\t"+strconv.Itoa(registerArray[22])+"\t"+strconv.Itoa(registerArray[23])+"\n")
			checkError(err6)
			_, err7 := io.WriteString(outputFileSim, "r24:"+"\t"+strconv.Itoa(registerArray[24])+"\t"+strconv.Itoa(registerArray[25])+"\t"+strconv.Itoa(registerArray[26])+"\t"+strconv.Itoa(registerArray[27])+"\t"+strconv.Itoa(registerArray[28])+"\t"+strconv.Itoa(registerArray[29])+"\t"+strconv.Itoa(registerArray[30])+"\t"+strconv.Itoa(registerArray[31])+"\n"+"\n")
			checkError(err7)

			_, err8 := io.WriteString(outputFileSim, "data:"+"\n")
			checkError(err8)

			//print data
			var holder int
			if len(dataArray) > 0 {
				for t := 0; t < len(dataArray); t++ {
					if t%8 == 0 {
						_, err := io.WriteString(outputFileSim, strconv.Itoa(dataMemoryStart+(t*4))+":"+"\t"+dataArray[t])
						checkError(err)
					} else if t%7 == 0 {
						_, err := io.WriteString(outputFileSim, "\t"+dataArray[t]+"\n")
						checkError(err)
					} else {
						_, err2 := io.WriteString(outputFileSim, "\t"+dataArray[t])
						checkError(err2)
					}
					holder = t
				}
				for f := holder % 8; f < 7; f++ {
					_, err2 := io.WriteString(outputFileSim, "\t"+"0")
					checkError(err2)
				}
			}

			_, err10 := io.WriteString(outputFileSim, "\n")
			checkError(err10)

			cycleNum++
			memoryLocation += 4
			if breakFlag {
				break
			}
		}
	}
}
