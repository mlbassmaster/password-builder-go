package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func main() {
	usage := "pwd-build [passwords to generate] [password length] [path to file] [-l]\n-l language\tspecify the language to use for special rules (optional)"

	var language string
	flag.StringVar(&language, "l", "en", "Specify the language to apply any special rules. Defaults to en. (English)")
	flag.Parse()

	if len(flag.Args()) != 3 {
		log.Fatal(usage)
	}

	generationCount, err := strconv.ParseInt(flag.Arg(0), 10, 64)
	if err != nil || generationCount < 1 {
		log.Fatal(usage)
	}

	passwordLength, err := strconv.ParseInt(flag.Arg(1), 10, 64)
	if err != nil || passwordLength < 1 {
		log.Fatal(usage)
	}

	fmt.Println("Reading file...")

	startingLetterCounts := make(map[string]int)

	followingLetterCounts := make(map[string]map[string]int)

	root, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error encountered running in path %s\n", root)
		log.Fatal(err)
	}

	relativePath := filepath.Join(root, flag.Arg(2))

	file, err := os.Open(relativePath)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	wordCount := 0

	for scanner.Scan() {
		wordCount++
		word := scanner.Text()
		firstLetter := strings.ToLower(string([]rune(word)[0]))

		for position, character := range word {
			letter := strings.ToLower(string(character))
			if position == 0 {
				startingLetterCounts[letter]++
			} else {
				innerMap, ok := followingLetterCounts[firstLetter]
				if !ok {
					innerMap = make(map[string]int)
					followingLetterCounts[firstLetter] = innerMap
				}
				followingLetterCounts[firstLetter][letter]++
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("generating %d %s passwords of length %d...s\n", generationCount, language, passwordLength)
	rand.Seed(time.Now().UnixNano())
	for i := int64(0); i < generationCount; i++ {
		password := ""

		//first letter
		roll := rand.Intn(wordCount)
		currentLetter := ""
		cumulativeWeightSum := 0
		for letter, weight := range startingLetterCounts {
			cumulativeWeightSum += weight
			if roll < cumulativeWeightSum {
				password += letter
				currentLetter = letter
				break
			}
		}

		//all other letters
		for j := int64(1); j < passwordLength; j++ {
			totalFollowerCount := 0
			for _, count := range followingLetterCounts[currentLetter] {
				totalFollowerCount += count
			}
			roll := rand.Intn(totalFollowerCount)

			cumulativeWeightSum = 0
			for letter, weight := range followingLetterCounts[currentLetter] {
				cumulativeWeightSum += weight
				if roll < cumulativeWeightSum {
					password += letter
					currentLetter = letter
					break
				}
			}
		}
		fmt.Println(password)
	}
}
