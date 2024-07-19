package utils

import "math/rand"

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func GetExtensionByMimeType(mimeType string) string {
	switch mimeType {
	case "image/png":
		return "png"
	case "image/jpeg":
		return "jpg"
	case "text/plain":
		return "txt"
	default:
		return ""
	}
}

func GetMimeTypeByExtension(extention string) string {
	switch extention {
	case "png":
		return "image/png"
	case "jpeg":
		return "image/jpg"
	case "txt":
		return "text/plain"
	default:
		return ""
	}
}

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
