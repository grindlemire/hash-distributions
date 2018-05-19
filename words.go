package main

// Word represents a hashable english word
type Word string

// ParseWords parses the words from a newline delmited file
func ParseWords(lines []string) (words []Word, err error) {
	for _, line := range lines {
		words = append(words, Word(line))
	}
	return words, nil
}

// Hash hashes the word using the specified hash
func (w Word) Hash(hashName string) (hashVal string) {
	return allHashes[hashName](string(w))
}
