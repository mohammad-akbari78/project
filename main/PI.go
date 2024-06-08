package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	file, oErr := os.OpenFile("A.pas", os.O_CREATE|os.O_RDWR, 0644)
	if oErr != nil {
		fmt.Println("can't open")

		return
	}
	var b = make([]byte, 256)
	_, rErr := file.Read(b)
	if rErr != nil {
		fmt.Println("can't read the file", rErr)
	}

	strB := string(b)
	strB = strings.Trim(strB, " ")

	var lineSlice []string
	var w string
	for _, x := range strB {
		if x == 10 {
			w = strings.Trim(w, " ")
			lineSlice = append(lineSlice, w)
			w = ""

			continue
		} else if x == 0 {
			lineSlice = append(lineSlice, w)

			break
		}
		w += string(x)
	}
	strB = strings.Replace(strB, "\n", " ", -1)
	strB = strings.Replace(strB, "\r", " ", -1)

	var new_str, s string
	var textSlice []string

	for i := 0; i < len(strB); i++ {
		if strB[i] == 32 && strB[i+1] == 32 {
			continue
		}
		new_str += string(strB[i])
	}
	for _, v := range new_str {
		if string(v) == " " {
			textSlice = append(textSlice, s)
			s = ""
			continue
		}

		s += string(v)
	}

	var newVarMapStr = make(map[string]string)
	var newVarMapReal = make(map[string]string)
	var newVarMapInt = make(map[string]string)
	var r int

	// *** read lines ***

	for x := 0; x <= len(lineSlice)-1; x++ {
		lineSlice[x] = strings.Trim(lineSlice[x], " ")
		lineSlice[x] = strings.Trim(lineSlice[x], " \r")
		if lineSlice[x] == "var" {
			newVarMapStr, newVarMapReal, newVarMapInt, r = varProcess(lineSlice)
			x = r

			continue

		} else if lineSlice[x] == "end." {

		} else {
			for i := 0; i < len(lineSlice[x]); i++ {
				if string(lineSlice[x][i]) == "i" && string(lineSlice[x][i+1]) == "f" {
					newVarMapInt, newVarMapReal, newVarMapStr, x = ifProcess(newVarMapInt, newVarMapReal, newVarMapStr, lineSlice, x)

					break
				} else if string(lineSlice[x][i]) == "f" && string(lineSlice[x][i+1]) == "o" && string(lineSlice[x][i+2]) == "r" {
					newVarMapInt, newVarMapReal, newVarMapStr, x = forProcess(newVarMapInt, newVarMapReal, newVarMapStr, lineSlice, x)
				} else if string(lineSlice[x][i]) == "w" {
					if string(lineSlice[x][i+1]) == "r" {
						if string(lineSlice[x][i+2]) == "i" {
							if string(lineSlice[x][i+3]) == "t" {
								if string(lineSlice[x][i+4]) == "e" {
									if string(lineSlice[x][i+5]) == "l" {
										if string(lineSlice[x][i+6]) == "n" {
											writelnProcess(newVarMapInt, newVarMapReal, newVarMapStr, lineSlice, x)

											break
										}
									}
								}
							}
						}
					}

				} else if string(lineSlice[x][i]) == "r" {
					if string(lineSlice[x][i+1]) == "e" {
						if string(lineSlice[x][i+2]) == "a" {
							if string(lineSlice[x][i+3]) == "d" {
								if string(lineSlice[x][i+4]) == "l" {
									if string(lineSlice[x][i+5]) == "n" {
										newVarMapInt, newVarMapReal, newVarMapStr = readlnProcess(newVarMapInt, newVarMapReal, newVarMapStr, lineSlice, x)

										break
									}
								}
							}
						}
					}

				} else if string(lineSlice[x][i]) == ":" && string(lineSlice[x][i+1]) == "=" {
					newVarMapInt, newVarMapReal, newVarMapStr = initialVariableProcess(i, newVarMapInt, newVarMapReal, newVarMapStr, lineSlice, x)

					break
				}

			}
		}

	}
}
func varProcess(s []string) (map[string]string, map[string]string, map[string]string, int) {
	var var_index, begin_index, r int
	var variableNames, secondSectionLine2 string
	var firstSplitLine2, variableNamesSplit []string
	var varMapInt = make(map[string]string)
	var varMapReal = make(map[string]string)
	var varMapString = make(map[string]string)
	var varDefMap = make(map[string]string)

	for i, v := range s {
		v = strings.Trim(v, "\r")
		if v == "var" {
			var_index = i
			for j, x := range s {
				x = strings.Trim(x, "\r")
				if x == "begin" {
					begin_index = j

					break
				}
			}

		}
	}
	r = begin_index

	for i := var_index + 1; i < begin_index; i++ {
		for q, p := range s[i] {
			var t = -2
			if string(p) == ":" {
				t = q
			}
			if q == t+1 && string(p) != "=" {

			}
		}

		firstSplitLine2 = strings.Split(s[i], ":")
		if len(firstSplitLine2) < 2 {

			continue
		}

		variableNames = firstSplitLine2[0]
		variableNames = strings.Trim(variableNames, " ")

		secondSectionLine2 = firstSplitLine2[1]

		variableNamesSplit = strings.Split(variableNames, ",")
		for _, y := range variableNamesSplit {
			//*** variable define checking
			// checking first char
			if ((y[0] > 90 || y[0] < 65) && (y[0] < 97 || y[0] > 122)) && y[0] != 32 {
				fmt.Println("!!! variable name is not valid !!!", y)

				os.Exit(0)
			}
			// checking valid characters
			for _, v := range y {
				if ((v > 57 || v < 48) && (v > 90 || v < 65) && (v < 97 || v > 122)) && v != 95 {
					fmt.Println("!!!* variable name is not valid *!!!", y)

					os.Exit(0)
				}
			}

			for key := range varMapInt {
				if key == y {
					fmt.Printf("variable %s is already defined as a int type", y)

					os.Exit(0)
				}
			}
			for key := range varMapString {
				if key == y {
					fmt.Printf("variable %s is already defined as a string type", y)

					os.Exit(0)
				}
			}
			for key := range varMapReal {
				if key == y {
					fmt.Printf("variable %s is already defined as a real type", y)

					os.Exit(0)
				}
			}
		}

		secondSectionLine2 = strings.Trim(secondSectionLine2, "\r")
		secondSectionLine2 = strings.Trim(secondSectionLine2, " ")
		if secondSectionLine2[len(secondSectionLine2)-1] != 59 {
			fmt.Println("you forgot the semicolon in line", i)
			os.Exit(0)
		} else {
			secondSectionLine2 = strings.Replace(secondSectionLine2, ";", "", -1)
			secondSectionLine2 = strings.Trim(secondSectionLine2, " ")

			switch secondSectionLine2 {
			case "integer":
				for _, d := range variableNamesSplit {
					varMapInt[d] = "0"

				}
			case "real":
				for _, d := range variableNamesSplit {
					varMapReal[d] = "0"

				}

			case "string":
				for _, d := range variableNamesSplit {
					varMapString[d] = ""

				}
			default:
				fmt.Println(" !!!type of value is invalid!!! ")
				os.Exit(0)

			}

		}

	}
	for _, v := range variableNamesSplit {
		varDefMap[v] = secondSectionLine2
	}

	if begin_index < var_index {
		err := fmt.Errorf(" In Pascal, variables must be defined before any statement ")
		fmt.Println(err)
	}
	return varMapString, varMapReal, varMapInt, r
}
func ifProcess(newVarMapInt map[string]string, newVarMapReal map[string]string, newVarMapStr map[string]string, lineSlice []string, x int) (map[string]string, map[string]string, map[string]string, int) {
	var fParentheses, lParentheses, newRow int
	var resKey, resKey2, conditionalPhrase, firstComparativeValue, secondComparativeValue string

	for m, n := range lineSlice[x] {
		if string(n) == "(" {
			fParentheses = m
		}
		if string(n) == ")" {
			lParentheses = m

			break
		}
	}
	for u := 0; u < fParentheses; u++ {
		resKey += string(lineSlice[x][u])
	}
	for b := lParentheses + 1; b < len(lineSlice[x]); b++ {
		resKey2 += string(lineSlice[x][b])
	}
	for c := fParentheses + 1; c < lParentheses; c++ {
		conditionalPhrase += string(lineSlice[x][c])
	}

	resKey = strings.Trim(resKey, " ")
	resKey2 = strings.Trim(resKey2, " ")
	resKey2 = strings.Trim(resKey2, "\r")
	conditionalPhrase = strings.Trim(conditionalPhrase, " ")
	if resKey != "if" {
		fmt.Printf("%+v is not reserved key     line: %d", resKey, x+1)

		os.Exit(0)
	}
	if resKey2 != "then" {
		fmt.Printf("%+v is not reserved key     line: %d", resKey2, x+1)
	}

	var firstValueResult, secondValueResult string
	var typeNumber1, typeNumber2 int
	var firstTemp, secondTemp float64
	var sErr, sErr2 error
	isFoundSign := false
	for e := len(conditionalPhrase) - 1; e > 0; e-- {
		if string(conditionalPhrase[e]) == "=" {
			isFoundSign = true
			if string(conditionalPhrase[e-1]) == "<" {
				for a := 0; a < e-1; a++ {
					firstComparativeValue += string(conditionalPhrase[a])
				}
				for b := e + 1; b < len(conditionalPhrase); b++ {
					secondComparativeValue += string(conditionalPhrase[b])
				}
				firstValueResult, typeNumber1 = comparativeValuesCheckingInt(firstComparativeValue, newVarMapInt, newVarMapReal, newVarMapStr)
				secondValueResult, typeNumber2 = comparativeValuesCheckingInt(secondComparativeValue, newVarMapInt, newVarMapReal, newVarMapStr)

				if typeNumber1 == 1 {
					firstTemp, sErr = strconv.ParseFloat(firstValueResult, 5)
				}
				if sErr != nil {
					fmt.Println("can't convert string to float")
				}
				if typeNumber2 == 1 {
					secondTemp, sErr2 = strconv.ParseFloat(secondValueResult, 2)
				}
				if sErr2 != nil {
					fmt.Println("can't convert string to float2")
				}

				if typeNumber2 == typeNumber1 {
					if firstTemp <= secondTemp {
						newRow, newVarMapInt, newVarMapReal, newVarMapStr = readLines(lineSlice, newVarMapInt, newVarMapReal, newVarMapStr, x)
						return newVarMapInt, newVarMapReal, newVarMapStr, newRow
					} else {
						newRow = fakeReadLines(lineSlice, x)

						return newVarMapInt, newVarMapReal, newVarMapStr, newRow
					}
				} else {
					fmt.Println("can't compare two different data types")

					os.Exit(0)
				}
			} else if string(conditionalPhrase[e-1]) == ">" {
				for a := 0; a < e-1; a++ {
					firstComparativeValue += string(conditionalPhrase[a])
				}
				for b := e + 1; b < len(conditionalPhrase); b++ {
					secondComparativeValue += string(conditionalPhrase[b])
				}
				firstValueResult, typeNumber1 = comparativeValuesCheckingInt(firstComparativeValue, newVarMapInt, newVarMapReal, newVarMapStr)
				secondValueResult, typeNumber2 = comparativeValuesCheckingInt(secondComparativeValue, newVarMapInt, newVarMapReal, newVarMapStr)

				if typeNumber1 == 1 {
					firstTemp, sErr = strconv.ParseFloat(firstValueResult, 5)
				}
				if sErr != nil {
					fmt.Println("can't convert string to float")
				}
				if typeNumber2 == 1 {
					secondTemp, sErr2 = strconv.ParseFloat(secondValueResult, 2)
				}
				if sErr2 != nil {
					fmt.Println("can't convert string to float2")
				}
				if typeNumber2 == typeNumber1 {
					if firstTemp >= secondTemp {
						newRow, newVarMapInt, newVarMapReal, newVarMapStr = readLines(lineSlice, newVarMapInt, newVarMapReal, newVarMapStr, x)
						return newVarMapInt, newVarMapReal, newVarMapStr, newRow
					} else {
						newRow = fakeReadLines(lineSlice, x)

						return newVarMapInt, newVarMapReal, newVarMapStr, newRow
					}
				} else {
					fmt.Println("can't compare two different data types")

					os.Exit(0)
				}
			} else {
				for a := 0; a < e; a++ {
					firstComparativeValue += string(conditionalPhrase[a])
				}
				for b := e + 1; b < len(conditionalPhrase); b++ {
					secondComparativeValue += string(conditionalPhrase[b])
				}
				firstValueResult, typeNumber1 = comparativeValuesCheckingInt(firstComparativeValue, newVarMapInt, newVarMapReal, newVarMapStr)
				secondValueResult, typeNumber2 = comparativeValuesCheckingInt(secondComparativeValue, newVarMapInt, newVarMapReal, newVarMapStr)

				if typeNumber1 == 1 {
					firstTemp, sErr = strconv.ParseFloat(firstValueResult, 5)
				}
				if sErr != nil {
					fmt.Println("can't convert string to float")
				}
				if typeNumber2 == 1 {
					secondTemp, sErr2 = strconv.ParseFloat(secondValueResult, 2)
				}
				if sErr2 != nil {
					fmt.Println("can't convert string to float2")
				}

				if typeNumber2 == typeNumber1 {
					if firstTemp == secondTemp {
						newRow, newVarMapInt, newVarMapReal, newVarMapStr = readLines(lineSlice, newVarMapInt, newVarMapReal, newVarMapStr, x)
						return newVarMapInt, newVarMapReal, newVarMapStr, newRow
					} else {
						newRow = fakeReadLines(lineSlice, x)

						return newVarMapInt, newVarMapReal, newVarMapStr, newRow
					}
				} else {
					fmt.Println("can't compare two different data types")

					os.Exit(0)
				}
			}
		} else if string(conditionalPhrase[e]) == ">" {
			isFoundSign = true
			if string(conditionalPhrase[e-1]) == "<" {
				for a := 0; a < e-1; a++ {
					firstComparativeValue += string(conditionalPhrase[a])
				}
				for b := e + 1; b < len(conditionalPhrase); b++ {
					secondComparativeValue += string(conditionalPhrase[b])
				}
				firstValueResult, typeNumber1 = comparativeValuesCheckingInt(firstComparativeValue, newVarMapInt, newVarMapReal, newVarMapStr)
				secondValueResult, typeNumber2 = comparativeValuesCheckingInt(secondComparativeValue, newVarMapInt, newVarMapReal, newVarMapStr)

				if typeNumber1 == 1 {
					firstTemp, sErr = strconv.ParseFloat(firstValueResult, 5)
				}
				if sErr != nil {
					fmt.Println("can't convert string to float")
				}
				if typeNumber2 == 1 {
					secondTemp, sErr2 = strconv.ParseFloat(secondValueResult, 2)
				}
				if sErr2 != nil {
					fmt.Println("can't convert string to float2")
				}

				if typeNumber2 == typeNumber1 {
					if firstTemp != secondTemp {
						newRow, newVarMapInt, newVarMapReal, newVarMapStr = readLines(lineSlice, newVarMapInt, newVarMapReal, newVarMapStr, x)
						return newVarMapInt, newVarMapReal, newVarMapStr, newRow
					} else {
						newRow = fakeReadLines(lineSlice, x)

						return newVarMapInt, newVarMapReal, newVarMapStr, newRow
					}
				} else {
					fmt.Println("can't compare two different data types")

					os.Exit(0)
				}

			} else {
				for a := 0; a < e; a++ {
					firstComparativeValue += string(conditionalPhrase[a])
				}
				for b := e + 1; b < len(conditionalPhrase); b++ {
					secondComparativeValue += string(conditionalPhrase[b])
				}

				firstValueResult, typeNumber1 = comparativeValuesCheckingInt(firstComparativeValue, newVarMapInt, newVarMapReal, newVarMapStr)
				secondValueResult, typeNumber2 = comparativeValuesCheckingInt(secondComparativeValue, newVarMapInt, newVarMapReal, newVarMapStr)

				if typeNumber1 == 1 {
					firstTemp, sErr = strconv.ParseFloat(firstValueResult, 5)
				}
				if sErr != nil {
					fmt.Println("can't convert string to float")
				}
				if typeNumber2 == 1 {
					secondTemp, sErr2 = strconv.ParseFloat(secondValueResult, 2)
				}
				if sErr2 != nil {
					fmt.Println("can't convert string to float2")
				}

				if typeNumber2 == typeNumber1 {
					if firstTemp > secondTemp {
						newRow, newVarMapInt, newVarMapReal, newVarMapStr = readLines(lineSlice, newVarMapInt, newVarMapReal, newVarMapStr, x)
						return newVarMapInt, newVarMapReal, newVarMapStr, newRow
					} else {
						newRow = fakeReadLines(lineSlice, x)

						return newVarMapInt, newVarMapReal, newVarMapStr, newRow
					}
				} else {
					fmt.Println("can't compare two different data types")

					os.Exit(0)
				}
			}
		} else if string(conditionalPhrase[e]) == "<" {
			isFoundSign = true
			for a := 0; a < e; a++ {
				firstComparativeValue += string(conditionalPhrase[a])
			}
			for b := e + 1; b < len(conditionalPhrase); b++ {
				secondComparativeValue += string(conditionalPhrase[b])
			}
			firstValueResult, typeNumber1 = comparativeValuesCheckingInt(firstComparativeValue, newVarMapInt, newVarMapReal, newVarMapStr)
			secondValueResult, typeNumber2 = comparativeValuesCheckingInt(secondComparativeValue, newVarMapInt, newVarMapReal, newVarMapStr)

			if typeNumber1 == 1 {
				firstTemp, sErr = strconv.ParseFloat(firstValueResult, 5)
			}
			if sErr != nil {
				fmt.Println("can't convert string to float")
			}
			if typeNumber2 == 1 {
				secondTemp, sErr2 = strconv.ParseFloat(secondValueResult, 2)
			}
			if sErr2 != nil {
				fmt.Println("can't convert string to float2")
			}
			if typeNumber2 == typeNumber1 {
				if firstTemp < secondTemp {
					newRow, newVarMapInt, newVarMapReal, newVarMapStr = readLines(lineSlice, newVarMapInt, newVarMapReal, newVarMapStr, x)
					return newVarMapInt, newVarMapReal, newVarMapStr, newRow
				} else {
					newRow = fakeReadLines(lineSlice, x)

					return newVarMapInt, newVarMapReal, newVarMapStr, newRow
				}
			} else {
				fmt.Println("can't compare two different data types")

				os.Exit(0)
			}
		} else {

			continue
		}
	}
	if !isFoundSign {
		fmt.Println(" !!! your conditional expression doesn't have a comparative sign !!! ")

		os.Exit(0)
	}

	return newVarMapInt, newVarMapReal, newVarMapStr, newRow
}
func forProcess(newVarMapInt map[string]string, newVarMapReal map[string]string, newVarMapStr map[string]string, lineSlice []string, x int) (map[string]string, map[string]string, map[string]string, int) {
	var indexOfTo, indexOfDo, indexOfSign int
	var varDefExpr, limExpr, varName, valueOfVariable string

	lineSlice[x] = strings.Trim(lineSlice[x], " ")
	lineSlice[x] = strings.Trim(lineSlice[x], "\r")

	if string(lineSlice[x][0]) != "f" || string(lineSlice[x][1]) != "o" || string(lineSlice[x][2]) != "r" {
		fmt.Println(" !!! for loop implementation is not correct !!!")

		os.Exit(0)
	}

	for i := 3; i < len(lineSlice[x]); i++ {
		if string(lineSlice[x][i]) == "t" && string(lineSlice[x][i+1]) == "o" {
			indexOfTo = i
		}
		if string(lineSlice[x][i]) == "d" && string(lineSlice[x][i+1]) == "o" {
			indexOfDo = i
		}
	}
	for j := 3; j < indexOfTo; j++ {
		varDefExpr += string(lineSlice[x][j])
	}
	for k := indexOfTo + 2; k < indexOfDo; k++ {
		limExpr += string(lineSlice[x][k])
	}
	varDefExpr = strings.Trim(varDefExpr, " ")
	limExpr = strings.Trim(limExpr, " ")

	for u := 0; u < len(varDefExpr); u++ {
		if string(varDefExpr[u]) == ":" && string(varDefExpr[u+1]) == "=" {
			indexOfSign = u
		}
	}
	for s := 0; s < indexOfSign; s++ {
		varName += string(varDefExpr[s])
	}
	for a := indexOfSign + 2; a < len(varDefExpr); a++ {
		valueOfVariable += string(varDefExpr[a])
	}
	valueOfVariable = strings.Trim(valueOfVariable, " ")
	varName = strings.Trim(varName, " ")

	isDef := false
	for index := range newVarMapInt {
		if index == varName {
			newVarMapInt[varName] = valueOfVariable
			isDef = true
		}
	}
	if !isDef {
		fmt.Printf("%v is not define ", varName)

		os.Exit(0)
	}
	isStr := false
	for p := range limExpr {
		if p != 48 && p != 49 && p != 50 && p != 51 && p != 52 && p != 53 && p != 54 && p != 55 && p != 56 && p != 57 {
			isStr = true

			break
		}
	}
	if isStr {
		for index, value := range newVarMapInt {
			if index == limExpr {
				limExpr = value
			}
		}
	}
	if !isStr {
		fmt.Printf("%v is not exist ", limExpr)

		os.Exit(0)
	}
	startValue, cErr := strconv.Atoi(valueOfVariable)
	if cErr != nil {
		fmt.Println("can't convert string to int , line", x+1)

		os.Exit(0)
	}
	endValue, coErr := strconv.Atoi(limExpr)
	if coErr != nil {
		fmt.Println("can't convert string to int , line", x+1)

		os.Exit(0)
	}
	for t := startValue; t < endValue-1; t++ {
		readLines(lineSlice, newVarMapInt, newVarMapReal, newVarMapStr, x)
		newVarMapInt[varName] = strconv.Itoa(t + 1)

	}
	x, newVarMapInt, newVarMapReal, newVarMapStr = readLines(lineSlice, newVarMapInt, newVarMapReal, newVarMapStr, x)

	return newVarMapInt, newVarMapReal, newVarMapStr, x
}
func writelnProcess(newVarMapInt map[string]string, newVarMapReal map[string]string, newVarMapStr map[string]string, lineSlice []string, x int) {
	var firstParentheses, lastParentheses, fQuote, lQuote int
	var printTextWithQuotation, printTextWithoutQuotation, printTextVariable, newPrintTextVariable string
	var isFound1, isFound2, isFound3 bool

	lineSlice[x] = strings.Trim(lineSlice[x], " ")
	lineSlice[x] = strings.Trim(lineSlice[x], "\r")
	if lineSlice[x][len(lineSlice[x])-1] != 59 {
		fmt.Println("you forgot the semicolon in line ", x+1)

		os.Exit(0)
	}
	for m, n := range lineSlice[x] {
		if string(n) == "(" {
			firstParentheses = m
		}
		if string(n) == ")" {
			lastParentheses = m

			break
		}
	}
	for i := firstParentheses + 1; i < lastParentheses; i++ {
		printTextWithQuotation += string(lineSlice[x][i])
	}
	printTextWithQuotation = strings.Trim(printTextWithQuotation, " ")
	for index, item := range printTextWithQuotation {
		if index == 0 && item == 34 {
			fQuote = index
			isFound1 = true
		} else if index == len(printTextWithQuotation)-1 && item == 34 {
			lQuote = index
			isFound2 = true
		}
	}
	if isFound2 && isFound1 {
		for i := fQuote + 1; i < lQuote; i++ {
			printTextWithoutQuotation += string(printTextWithQuotation[i])
		}
		fmt.Println(printTextWithoutQuotation)
	} else if isFound1 && !isFound2 {
		fmt.Println("you forgot the quotation in line ", x+1)

		os.Exit(0)
	} else if !isFound1 && isFound2 {
		fmt.Println("you forgot the quotation in line ", x+1)

		os.Exit(0)
	} else {
		for _, t := range printTextWithQuotation {
			printTextVariable += string(t)
		}
		for k, v := range newVarMapInt {
			if k == printTextVariable {
				newPrintTextVariable = v
				isFound3 = true

				break
			}
		}
		for key, value := range newVarMapReal {
			if key == printTextVariable {
				newPrintTextVariable = value
				isFound3 = true

				break
			}
		}
		for key, value := range newVarMapStr {
			if key == printTextVariable {
				newPrintTextVariable = value
				isFound3 = true

				break
			}
		}
		if isFound3 {
			fmt.Println(newPrintTextVariable)
		} else {
			fmt.Println(printTextVariable, " is not declare!!!")

			os.Exit(0)
		}

	}

}
func readlnProcess(newVarMapInt map[string]string, newVarMapReal map[string]string, newVarMapStr map[string]string, lineSlice []string, x int) (map[string]string, map[string]string, map[string]string) {
	var firstParentheses, lastParentheses int
	var readText string
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	inputValue := scanner.Text()

	lineSlice[x] = strings.Trim(lineSlice[x], " ")
	lineSlice[x] = strings.Trim(lineSlice[x], "\r")
	if lineSlice[x][len(lineSlice[x])-1] != 59 {
		fmt.Println("sss", x)
		fmt.Println("you forgot the semicolon in line ", x+1)

		os.Exit(0)
	}
	for m, n := range lineSlice[x] {
		if string(n) == "(" {
			firstParentheses = m
		}
		if string(n) == ")" {
			lastParentheses = m

			break
		}
	}
	for i := firstParentheses + 1; i < lastParentheses; i++ {
		readText += string(lineSlice[x][i])
	}
	readText = strings.Trim(readText, " ")
	_, typeNum := comparativeValuesCheckingInt(readText, newVarMapInt, newVarMapReal, newVarMapStr)

	point := false
	if typeNum == 1 {
		for _, item := range inputValue {
			if item == 48 || item == 49 || item == 50 || item == 51 || item == 52 || item == 53 || item == 54 || item == 55 || item == 56 || item == 57 {
			} else if item == 46 && !point {
				point = true
			} else {
				fmt.Println("!!! The value of variable is not compatible with the type of variable !!!")

				os.Exit(0)
			}
		}
		for index := range newVarMapInt {
			if index == readText {
				newVarMapInt[readText] = inputValue
			}
		}
		for index := range newVarMapReal {
			if index == readText {
				newVarMapReal[readText] = inputValue
			}
		}
	} else {
		for index := range newVarMapStr {
			if index == readText {
				newVarMapStr[readText] = inputValue
			}
		}
	}

	return newVarMapInt, newVarMapReal, newVarMapStr
}
func initialVariableProcess(j int, newVarMapInt map[string]string, newVarMapReal map[string]string, newVarMapStr map[string]string, lineSlice []string, x int) (map[string]string, map[string]string, map[string]string) {
	var initVar, initVal, newInitVal, initValSection1, initValSection2, initValSign string
	var n int
	a := j
	for b := 0; b < a; b++ {
		initVar += string(lineSlice[x][b])
	}
	for c := a + 2; c < len(lineSlice[x]); c++ {
		initVal += string(lineSlice[x][c])
	}
	initVar = strings.Trim(initVar, " ")
	initVal = strings.Trim(initVal, " ")

	if initVal[len(initVal)-1] != 59 {
		fmt.Println("you forgot the semicolon in line", x+1)
		os.Exit(0)
	} else {
		initVal = strings.Replace(initVal, string(59), "", -1)
	}

	for a, b := range initVal {
		if b == 42 || b == 43 || b == 45 {
			n = a
			initValSign = string(b)
		}
	}
	for f := 0; f < n; f++ {
		initValSection1 += string(initVal[f])
	}
	for j := n + 1; j < len(initVal); j++ {
		initValSection2 += string(initVal[j])
	}
	exist := false
	notNum := false
	isDef := false

	if initValSection1 != "" && initValSection2 != "" {
		for _, p := range initValSection1 {
			if p != 48 && p != 49 && p != 50 && p != 51 && p != 52 && p != 53 && p != 54 && p != 55 && p != 56 && p != 57 && p != 46 {
				notNum = true

				break
			}
		}
		if notNum {
			for index, value := range newVarMapInt {
				if index == initValSection1 {
					initValSection1 = value
					isDef = true

					break
				}
			}
			for index, value := range newVarMapReal {
				if index == initValSection1 {
					initValSection1 = value
					isDef = true

					break
				}
			}
			for index := range newVarMapStr {
				if index == initValSection1 {
					fmt.Println("!!! we can't use logical operators for strings !!! ")

					os.Exit(0)
				}
			}
			if !isDef {
				fmt.Printf("%v is not define	\n", initValSection1)

				os.Exit(0)
			}
			notNum = false
			isDef = false
			initValSection1 = strings.Trim(initValSection1, " ")
		}
		for _, p := range initValSection2 {
			if p != 48 && p != 49 && p != 50 && p != 51 && p != 52 && p != 53 && p != 54 && p != 55 && p != 56 && p != 57 && p != 46 {
				notNum = true

				break
			}
		}
		if notNum {
			for index, value := range newVarMapInt {
				if index == initValSection2 {
					initValSection2 = value
					isDef = true

					break
				}
			}
			for index, value := range newVarMapReal {
				if index == initValSection2 {
					initValSection2 = value
					isDef = true

					break
				}
			}
			for index := range newVarMapStr {
				if index == initValSection2 {
					fmt.Println("!!! we can't use logical operators for strings !!! ")

					os.Exit(0)
				}
			}
			if !isDef {
				fmt.Printf("%v is not define	\n", initValSection2)

				os.Exit(0)
			}
			notNum = false
			isDef = false
		}
		if initValSign == "*" {
			initValSec1Int, fErr := strconv.Atoi(initValSection1)
			if fErr != nil {
				fmt.Println("can't convert string to int")

				os.Exit(0)
			}
			initValSec2Int, fErr := strconv.Atoi(initValSection2)
			if fErr != nil {
				fmt.Println("can't convert string to int")

				os.Exit(0)
			}
			result := initValSec1Int * initValSec2Int
			initVal = strconv.Itoa(result)
			for c := range newVarMapInt {
				if c == initVar {
					newVarMapInt[initVar] = initVal
				}
			}
			for c := range newVarMapReal {
				if c == initVar {
					newVarMapReal[initVar] = initVal
				}
			}
		} else if initValSign == "+" {
			initValSec1Int, fErr := strconv.Atoi(initValSection1)
			if fErr != nil {
				fmt.Println("can't convert string to int")

				os.Exit(0)
			}
			initValSec2Int, fErr := strconv.Atoi(initValSection2)
			if fErr != nil {
				fmt.Println("can't convert string to int")

				os.Exit(0)
			}
			result := initValSec1Int + initValSec2Int
			initVal = strconv.Itoa(result)
			for c := range newVarMapInt {
				if c == initVar {
					newVarMapInt[initVar] = initVal
				}
			}
			for c := range newVarMapReal {
				if c == initVar {
					newVarMapReal[initVar] = initVal
				}
			}
		} else if initValSign == "-" {
			initValSec1Int, fErr := strconv.Atoi(initValSection1)
			if fErr != nil {
				fmt.Println("can't convert string to int")

				os.Exit(0)
			}
			initValSec2Int, fErr := strconv.Atoi(initValSection2)
			if fErr != nil {
				fmt.Println("can't convert string to int")

				os.Exit(0)
			}
			result := initValSec1Int - initValSec2Int
			initVal = strconv.Itoa(result)
			for c := range newVarMapInt {
				if c == initVar {
					newVarMapInt[initVar] = initVal
				}
			}
			for c := range newVarMapReal {
				if c == initVar {
					newVarMapReal[initVar] = initVal
				}
			}
		}
	} else if initValSection2 == "" && initValSection1 == "" {
		for index := range newVarMapInt {
			if index == initVar {
				for _, item := range initVal {
					if item != 48 && item != 49 && item != 50 && item != 51 && item != 52 && item != 53 && item != 54 && item != 55 && item != 56 && item != 57 {
						fmt.Println("!!! The value of variable is not compatible with the type of variable ,int")

						os.Exit(0)
					}
				}
				newVarMapInt[index] = initVal
				exist = true
			}
		}
		for index := range newVarMapStr {
			if index == initVar {
				if initVal[0] == 34 && initVal[len(initVal)-1] == 34 {
					for i, a := range initVal {
						if i == 0 || i == len(initVal)-1 {

							continue
						} else {
							newInitVal += string(a)
						}
					}
					newVarMapStr[index] = newInitVal
				} else {
					fmt.Println("this type of value is not valid for string", initVal)

					os.Exit(0)
				}
				exist = true
			}
		}
		for index := range newVarMapReal {
			if index == initVar {
				point := false
				for _, item := range initVal {
					if item == 48 || item == 49 || item == 50 || item == 51 || item == 52 || item == 53 || item == 54 || item == 55 || item == 56 || item == 57 {
					} else if item == 46 && !point {
						point = true
					} else {
						fmt.Println("!!! The value of variable is not compatible with the type of variable ,real")

						os.Exit(0)
					}
				}
				newVarMapReal[index] = initVal
				exist = true
			}
		}
		if !exist {
			fmt.Printf("%+v is not define", initVar)

			os.Exit(0)
		}
	} else {
		fmt.Println("initialization is not correct")

		os.Exit(0)
	}

	exist = false
	initVar = ""
	initVal = ""
	initValSection2 = ""
	initValSection1 = ""
	initValSign = ""

	return newVarMapInt, newVarMapReal, newVarMapStr
}
func comparativeValuesCheckingInt(s string, i map[string]string, r map[string]string, st map[string]string) (string, int) {
	var q string
	var w int
	s = strings.Trim(s, " ")
	if len(s) > 1 {
		if (s[0] < 91 && s[0] > 64) || (s[0] > 96 && s[0] < 123) {
			for _, h := range s {
				if ((h > 57 || h < 48) && (h > 90 || h < 65) && (h < 97 || h > 122)) && h != 95 {
					fmt.Println("this type of expressions is not valid for variable")

					os.Exit(0)
				}
			}
			isFound := false
			for b, c := range i {
				if b == s {
					q = c
					w = 1
					isFound = true

					break
				}
			}
			if !isFound {
				for d, e := range r {
					if d == s {
						q = e
						w = 1
						isFound = true

						break
					}
				}
			}
			if !isFound {
				for f, g := range st {
					if f == s {
						q = g
						w = 2
						isFound = true

						break
					}
				}
			}
			if !isFound {
				fmt.Printf("%+v is not defined", s)

				os.Exit(0)
			}
		} else if s[0] < 58 && s[0] > 47 {
			point := false
			for _, i := range s {
				if i < 58 && i > 47 {

				} else if i == 46 && !point {
					point = true
				} else {
					fmt.Println("this value is not valid", s)

					os.Exit(0)
				}
			}
			q = s
			w = 1
		} else if s[0] == 34 && s[len(s)-1] == 34 {
			for j := 1; j < len(s)-1; j++ {
				q += string(s[j])
			}
			w = 2
		} else {
			fmt.Println("this value is not valid", s)

			os.Exit(0)
		}
	} else {
		if (s < string(91) && s > string(64)) || (s > string(96) && s < string(123)) {
			isFound := false
			for b, c := range i {
				if b == s {
					q = c
					w = 1
					isFound = true

					break
				}
			}
			if !isFound {
				for d, e := range r {
					if d == s {
						q = e
						w = 1
						isFound = true

						break
					}
				}
			}
			if !isFound {
				for f, g := range st {
					if f == s {
						q = g
						w = 2
						isFound = true

						break
					}
				}
			}
			if !isFound {
				fmt.Printf("%+v is not defined", s)

				os.Exit(0)
			}
		} else if s < string(58) && s > string(47) {
			q = s
			w = 1
		}
	}
	return q, w
}
func readLines(lineSlice []string, newVarMapInt map[string]string, newVarMapReal map[string]string, newVarMapStr map[string]string, x int) (int, map[string]string, map[string]string, map[string]string) {
	var z int
	for z = x + 1; z < len(lineSlice)-1; z++ {
		lineSlice[z] = strings.Trim(lineSlice[z], " ")
		lineSlice[z] = strings.Trim(lineSlice[z], "\r")

		if lineSlice[z] == "begin" || lineSlice[z] == "" {

			continue
		} else if lineSlice[z] == "end;" {

			return z, newVarMapInt, newVarMapReal, newVarMapStr
		} else {
			for i := 0; i < len(lineSlice[z]); i++ {
				if string(lineSlice[z][i]) == "i" && string(lineSlice[z][i+1]) == "f" {
					newVarMapInt, newVarMapReal, newVarMapStr, z = ifProcess(newVarMapInt, newVarMapReal, newVarMapStr, lineSlice, z)

					break
				} else if string(lineSlice[z][i]) == "f" && string(lineSlice[z][i+1]) == "o" && string(lineSlice[z][i+2]) == "r" {
					newVarMapInt, newVarMapReal, newVarMapStr, z = forProcess(newVarMapInt, newVarMapReal, newVarMapStr, lineSlice, z)
				} else if string(lineSlice[z][i]) == "w" && string(lineSlice[z][i+1]) == "r" && string(lineSlice[z][i+2]) == "i" && string(lineSlice[z][i+3]) == "t" && string(lineSlice[z][i+4]) == "e" && string(lineSlice[z][i+5]) == "l" && string(lineSlice[z][i+6]) == "n" {
					writelnProcess(newVarMapInt, newVarMapReal, newVarMapStr, lineSlice, z)

					continue
				} else if string(lineSlice[z][i]) == "r" && string(lineSlice[z][i+1]) == "e" && string(lineSlice[z][i+2]) == "a" && string(lineSlice[z][i+3]) == "d" && string(lineSlice[z][i+4]) == "l" && string(lineSlice[z][i+5]) == "n" {
					newVarMapInt, newVarMapReal, newVarMapStr = readlnProcess(newVarMapInt, newVarMapReal, newVarMapStr, lineSlice, z)

					continue
				} else if string(lineSlice[z][i]) == ":" && string(lineSlice[z][i+1]) == "=" {
					newVarMapInt, newVarMapReal, newVarMapStr = initialVariableProcess(i, newVarMapInt, newVarMapReal, newVarMapStr, lineSlice, z)

					break
				}

			}
		}

	}
	return z, newVarMapInt, newVarMapReal, newVarMapStr
}
func fakeReadLines(lineSlice []string, x int) int {
	var z, endRow, counter int

	for z = x + 1; z < len(lineSlice)-1; z++ {

		lineSlice[z] = strings.Trim(lineSlice[z], " ")
		lineSlice[z] = strings.Trim(lineSlice[z], "\r")

		if lineSlice[z] == "begin" {
			counter += 1

			continue
		} else if lineSlice[z] == "end;" {
			counter -= 1
			if counter == 0 {
				endRow = z

				return endRow
			}

		} else {

			continue
		}

	}
	return endRow
}
