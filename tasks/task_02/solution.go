package main

func rotateRunes(s string, shift int) string {
	runes := []rune(s)
	if len(runes) == 0 {
		return ""
	}

	shift = shift % len(runes)
	if shift < 0 {
		shift += len(runes)
	}

	r := make([]rune, 0, len(runes))
	r = append(r, runes[shift:]...)
	r = append(r, runes[:shift]...)

	return string(r)
}
