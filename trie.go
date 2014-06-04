package trie

type Trie interface {
	Search(string) bool
	PrefixSearch(string) (string, bool)
	CommonPrefixSearch(string) []string
}
