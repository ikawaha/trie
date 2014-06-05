package trie

import (
	"github.com/ikawaha/trie/da"

	"bufio"
	"errors"
	"os"
	"sort"
)

func NewDoubleArrayTrie(a_src interface{}) (Trie, error) {
	switch a_src.(type) {
	case []string:
		return newDoubleArrayTrieKeywords(a_src.([]string)), nil
	case *os.File:
		return newDoubleArrayTrieFile(a_src.(*os.File))
	default:
		return da.NewDoubleArray(), errors.New("cannot open unknown type src, '[]string' or '*os.File' can be specified.")
	}
}

func newDoubleArrayTrieKeywords(a_keywords []string) Trie {
	sort.Strings(a_keywords)
	da := da.NewDoubleArray()
	for _, keyword := range a_keywords {
		da.Append(keyword)
	}
	return da
}

func newDoubleArrayTrieFile(a_file *os.File) (Trie, error) {
	da := da.NewDoubleArray()
	scanner := bufio.NewScanner(a_file)
	for scanner.Scan() {
		da.Append(scanner.Text())
	}
	return da, scanner.Err()
}
