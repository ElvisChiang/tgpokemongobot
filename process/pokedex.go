package process

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

const filePrefix = "./resources/"
const fileSufix = "_vectorized.png"

// LoadGameData load data from pokedex file
func LoadGameData(nameFile string) (gameDataArray []GameData, ok bool) {
	gameDataArray = make([]GameData, 0)

	f, err := os.Open(nameFile)
	if err != nil {
		log.Print(err)
		return nil, false
	}

	r := csv.NewReader(bufio.NewReader(f))
	for {
		record, err := r.Read()
		// Stop at EOF.
		if err == io.EOF {
			break
		}
		if len(record) < 6 {
			continue
		}

		number, _ := strconv.Atoi(record[0])
		name := record[1]
		type1 := record[2]
		type2 := record[3]
		nickname := strings.Split(record[4], "_")
		evolve := strings.Split(record[5], "_")
		fmt.Printf("evolve %v", evolve)
		if evolve[0] == "-" {
			evolve = evolve[:0]
		}
		avatarFile := fmt.Sprintf("%s%03d%s", filePrefix, number, fileSufix)

		data := GameData{number, name, type1, type2, nickname, evolve, avatarFile}
		gameDataArray = append(gameDataArray, data)

		ok = true
	}

	return
}

// FindPokemon compare data in array
func FindPokemon(gameData []GameData, msg string) (pokemon GameData, ok bool) {
	ok = false
	if len(msg) == 0 {
		return
	}
	num, _ := strconv.Atoi(msg)
	if num > 0 {
		fmt.Printf("finding for pokemon #%d\n", num)
		for _, data := range gameData {
			if data.Number == num {
				pokemon = data
				ok = true
				return
			}
		}
		return
	}
	fmt.Printf("finding for pokemon name `%s`\n", msg)
	for _, data := range gameData {
		if strings.Contains(strings.ToLower(data.Name), msg) {
			pokemon = data
			ok = true
			return
		}
	}
	fmt.Printf("finding for pokemon nickname `%s`\n", msg)
	for _, data := range gameData {
		for _, nickname := range data.Nickname {
			if strings.Contains(nickname, msg) {
				pokemon = data
				ok = true
				return
			}
		}
	}

	return
}
