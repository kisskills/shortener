package cases

// Hasher TODO
func Hasher(originalLink string) string {
	if len(originalLink) < 10 {
		return originalLink
	}
	return originalLink[0:10]
}
