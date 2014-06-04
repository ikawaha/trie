package trie

import (
	"github.com/ikawaha/trie/da"

	"bufio"
	"os"
	"sort"
)

func NewDoubleArrayTrieFromKeywords(a_keywords []string) Trie {
	sort.Strings(a_keywords)
	da := da.NewDoubleArray()
	for _, keyword := range a_keywords {
		da.Append(keyword)
	}
	return da
}

func NewDoubleArrayTrieFromFile(a_file *os.File) (Trie, error) {
	da := da.NewDoubleArray()
	scanner := bufio.NewScanner(a_file)
	for scanner.Scan() {
		da.Append(scanner.Text())
	}
	return da, scanner.Err()
}
