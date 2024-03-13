package utils

func CutTextFromLastIndex(text string, index int) string {
	return text[:len(text)-index]
}
