package main

import (
	"bufio"
	"crypto/aes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
)

func GetUserInput() {
	text := StringPrompt("Cypher")
	encoded := base64.StdEncoding.EncodeToString([]byte(text))
	fmt.Printf("%s\n", encoded)
}

func main() {
	mode := StringPrompt("(E)ncrypt or (D)ecrypt?")
	key := StringPrompt("Cypher?")
	text := StringPrompt("Text?")
	switch mode {
	case "E":
		EncryptAES([]byte(key), text)
	case "D":
		DecryptAES([]byte(key), text)
	default:
		fmt.Println("Invalid input received. Try again.")
	}
}

// EncryptAES encrypts a plaintext-string
func EncryptAES(key []byte, plaintext string) {

	c, err := aes.NewCipher(key)
	CheckError(err)

	out := make([]byte, len(plaintext))

	c.Encrypt(out, []byte(plaintext))

	fmt.Println(hex.EncodeToString(out))
}

// DecryptAES decrypts an encrypted string
func DecryptAES(key []byte, ct string) {
	ciphertext, _ := hex.DecodeString(ct)

	c, err := aes.NewCipher(key)
	CheckError(err)

	pt := make([]byte, len(ciphertext))
	c.Decrypt(pt, ciphertext)

	s := string(pt[:])
	fmt.Println(s)
}

// CheckError is basic error handling
func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

// StringPrompt asks for a string value using the label
func StringPrompt(label string) string {
	var s string
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stderr, label+" ")
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}
	return strings.TrimSpace(s)
}
