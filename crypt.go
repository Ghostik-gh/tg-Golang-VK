package main

import (
	"fmt"
)

func Crypt() {
	for i := 33; i <= 122; i++ {
		fmt.Printf("%c", i)
	}
	fmt.Printf("%v\n", Decrypt(Encrypt("!\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz")))
	fmt.Println("!\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz")
}

const offset = 1198

func Encrypt(s string) string {
	src := []byte(s)
	dst := []byte("")
	fmt.Printf("src: %v\n", src)
	for _, v := range src {
		dst = append(dst, byte(((int(v)+offset)%90)+33))
	}
	fmt.Printf("dst: %v\n", dst)
	return string(dst[:])
}

func Decrypt(s string) string {
	src := []byte(s)
	dst := []byte("")
	fmt.Printf("src: %v\n", src)
	for _, v := range src {
		dst = append(dst, byte((int(v)-32+offset)%90+33))
	}
	fmt.Printf("dst: %v\n", dst)
	return string(dst[:])
}
