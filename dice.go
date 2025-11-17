package main

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RollDice rolls dice in the format "2d6", "1d20", etc.
// Returns the total result
func RollDice(notation string) int {
	notation = strings.ToLower(strings.TrimSpace(notation))

	// Parse notation like "2d6" or "1d20"
	parts := strings.Split(notation, "d")
	if len(parts) != 2 {
		return 0
	}

	count, err1 := strconv.Atoi(parts[0])
	sides, err2 := strconv.Atoi(parts[1])

	if err1 != nil || err2 != nil || count < 1 || sides < 1 {
		return 0
	}

	total := 0
	for i := 0; i < count; i++ {
		total += rand.Intn(sides) + 1
	}

	return total
}

// RollDiceDetailed rolls dice and returns individual results
func RollDiceDetailed(notation string) ([]int, int) {
	notation = strings.ToLower(strings.TrimSpace(notation))

	parts := strings.Split(notation, "d")
	if len(parts) != 2 {
		return []int{}, 0
	}

	count, err1 := strconv.Atoi(parts[0])
	sides, err2 := strconv.Atoi(parts[1])

	if err1 != nil || err2 != nil || count < 1 || sides < 1 {
		return []int{}, 0
	}

	results := make([]int, count)
	total := 0

	for i := 0; i < count; i++ {
		roll := rand.Intn(sides) + 1
		results[i] = roll
		total += roll
	}

	return results, total
}

// D6 rolls a single 6-sided die
func D6() int {
	return rand.Intn(6) + 1
}

// D10 rolls a single 10-sided die
func D10() int {
	return rand.Intn(10) + 1
}

// D20 rolls a single 20-sided die
func D20() int {
	return rand.Intn(20) + 1
}

// D100 rolls a percentile die
func D100() int {
	return rand.Intn(100) + 1
}
