package gamble

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"regexp"
	"strconv"
	"strings"
)

const (
	diceRollRegex = `^(\d*)d(\d+)(([+-]\d+)+)?$`
	modifierRegex = `[+-]?\d+`
)

// Roll parses standard dice notation (for example "2d6", "d20" or "3d8+2-1") and rolls the dice.
// A bare number is treated as the number of sides on a single die.
// It returns a human-readable breakdown of the individual rolls and modifiers, the final total, and an error if the input is malformed.
func Roll(input string) (string, int, error) {
	// Format input string
	formattedInput := strings.ToLower(strings.ReplaceAll(input, " ", ""))
	if num, err := strconv.Atoi(formattedInput); err == nil {
		formattedInput = fmt.Sprintf("d%d", num)
	}

	// Check if the input string matches the dice roll regex
	matches := regexp.MustCompile(diceRollRegex).FindStringSubmatch(formattedInput)
	if len(matches) < 1 {
		return "", 0, errors.New("gamble: invalid input format")
	}

	// Parse the number of rolls (default to 1 if not present)
	rollCount := 1
	if matches[1] != "" {
		var err error
		rollCount, err = strconv.Atoi(matches[1])
		if err != nil {
			return "", 0, fmt.Errorf("gamble: invalid number of rolls: %w", err)
		}
	}

	// Parse the number of sides on the dice
	sides, err := strconv.Atoi(matches[2])
	if err != nil {
		return "", 0, fmt.Errorf("gamble: invalid number of sides: %w", err)
	}

	// Parse the modifiers and evaluate it
	modifiers := matches[3]
	modifier, err := evaluateModifiers(modifiers)
	if err != nil {
		return "", 0, fmt.Errorf("gamble: invalid modifier: %w", err)
	}

	// Roll the dice and sum the results
	totalRoll := 0
	var rolls []int
	for range rollCount {
		roll := rand.IntN(sides) + 1
		rolls = append(rolls, roll)
		totalRoll += roll
	}
	result := totalRoll + modifier

	// Construct the output string showing individual rolls and modifiers
	var rollStrings []string
	for _, roll := range rolls {
		rollStrings = append(rollStrings, "("+strconv.Itoa(roll)+")")
	}
	rollOutput := strings.Join(rollStrings, " + ")
	if modifiers != "" {
		// Split the modifiers into individual terms with spaces
		modifierTerms := regexp.MustCompile(modifierRegex).FindAllString(modifiers, -1)
		for _, modifierTerm := range modifierTerms {
			if string(modifierTerm[0]) == "+" {
				modifierTerm = modifierTerm[1:]
			}
			rollOutput += " + (" + modifierTerm + ")"
		}
	}

	return rollOutput, result, nil
}

// evaluateModifiers evaluates the arithmetic expression in the modifiers.
func evaluateModifiers(modifiers string) (int, error) {
	if modifiers == "" {
		return 0, nil
	}

	// Split the modifiers into individual terms
	terms := regexp.MustCompile(modifierRegex).FindAllString(modifiers, -1)
	total := 0

	// Sum up all the terms
	for _, term := range terms {
		value, err := strconv.Atoi(term)
		if err != nil {
			return 0, err
		}
		total += value
	}

	return total, nil
}
