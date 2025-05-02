package app

func deleteFromSlice[K comparable](s []K, index int) []K {
	return append(s[:index], s[index+1:]...)
}
