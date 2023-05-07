package main

const offset = 1198

func Encrypt(s string) string {
	src := []byte(s)
	dst := []byte("")
	for _, v := range src {
		dst = append(dst, byte(((int(v)+offset)%90)+33))
	}
	return string(dst[:])
}

func Decrypt(s string) string {
	src := []byte(s)
	dst := []byte("")
	for _, v := range src {
		dst = append(dst, byte((int(v)-32+offset)%90+33))
	}
	return string(dst[:])
}
