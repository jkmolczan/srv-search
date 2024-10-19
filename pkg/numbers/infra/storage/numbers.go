package storage

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type NumbersStorage struct {
	data []int
}

func NewNumbersStorage(filepath string) (*NumbersStorage, error) {
	ns := &NumbersStorage{}
	err := ns.loadData(filepath)
	if err != nil {
		return nil, err
	}
	return ns, nil
}

func (ns *NumbersStorage) GetSortedCollection() ([]int, error) {
	return ns.data, nil
}

func (ns *NumbersStorage) loadData(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue // Skip empty lines
		}
		num, err := strconv.Atoi(line)
		if err != nil {
			return fmt.Errorf("failed to parse number on line %d: %v", lineNumber, err)
		}
		ns.data = append(ns.data, num)
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}
	return nil
}
