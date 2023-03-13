package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"testing"
)

type MenuItem struct {
	FoodMenuID int
	Count      int
}

func TopMenuItems(logFile string) ([]MenuItem, error) {
	menuItems := make(map[int]int)

	// Open the log file
	file, err := os.Open(logFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")

		// Parse the eater ID and food menu ID
		eaterID, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("error parsing eater ID: %v", err)
		}
		foodMenuID, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("error parsing food menu ID: %v", err)
		}

		// Check for duplicate entries
		if count, ok := menuItems[foodMenuID]; ok {
			if eaterID == count {
				return nil, fmt.Errorf("error: duplicate entry for eater ID %d and food menu ID %d", eaterID, foodMenuID)
			}
		}

		// Increment the count for the food menu ID
		menuItems[foodMenuID]++
	}

	// Convert the map to a slice of MenuItem values
	menuItemList := make([]MenuItem, 0, len(menuItems))
	for foodMenuID, count := range menuItems {
		menuItemList = append(menuItemList, MenuItem{FoodMenuID: foodMenuID, Count: count})
	}

	// Sort the slice by count in descending order and food menu ID in ascending order
	sort.Slice(menuItemList, func(i, j int) bool {
		if menuItemList[i].Count == menuItemList[j].Count {
			return menuItemList[i].FoodMenuID < menuItemList[j].FoodMenuID
		}
		return menuItemList[i].Count > menuItemList[j].Count
	})

	// Return the top 3 menu items
	if len(menuItemList) > 3 {
		menuItemList = menuItemList[:3]
	}
	return menuItemList, nil
}

func TestTopMenuItems(t *testing.T) {
	testCases := []struct {
		logFile    string
		expected   []MenuItem
		expectFail bool
	}{
		{
			logFile: "testdata/log1.txt",
			expected: []MenuItem{
				{FoodMenuID: 1, Count: 3},
				{FoodMenuID: 2, Count: 2},
				{FoodMenuID: 3, Count: 1},
			},
			expectFail: false,
		},
		{
			logFile: "testdata/log2.txt",
			expected: []MenuItem{
				{FoodMenuID: 1, Count: 2},
				{FoodMenuID: 2, Count: 1},
				{FoodMenuID: 3, Count: 1},
			},
			expectFail: false,
		},
		{
			logFile:    "testdata/log3.txt",
			expected:   []MenuItem{},
			expectFail: false,
		},
		{
			logFile:    "testdata/log4.txt",
			expected:   []MenuItem{},
			expectFail: true,
		},
	}

	for _, tc := range testCases {
		actual, err := TopMenuItems(tc.logFile)
		if err != nil && !tc.expectFail {
			t.Errorf("unexpected error: %v", err)
		}
		if err == nil && tc.expectFail {
			t.Errorf("expected an error, but got nil")
		}
		if err != nil && tc.expectFail {
			continue
		}
		if !menuItemsEqual(actual, tc.expected) {
			t.Errorf("incorrect top menu items: expected %v, but got %v", tc.expected, actual)
		}
	}
}

func menuItemsEqual(actual, expected []MenuItem) bool {
	if len(actual) != len(expected) {
		return false
	}
	for i := range actual {
		if actual[i].FoodMenuID != expected[i].FoodMenuID || actual[i].Count != expected[i].Count {
			return false
		}
	}
	return true
}

func main() {
	menuItems, err := TopMenuItems("data/log1.txt")
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	fmt.Println(menuItems)
}
