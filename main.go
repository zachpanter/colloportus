package main

import (
	"bufio"
	"crypto/aes"
	"encoding/hex"
	"flag"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"crypto/cipher"

	"crypto/rand"
	"fmt"
)

func main() {
	//mode := StringPrompt("(E)ncrypt or (D)ecrypt?")
	//key := StringPrompt("Cypher?")
	//text := StringPrompt("Text?")

	openmode := flag.String("open-mode", "", "(F)ile or (T)ext")
	filename := flag.String("file", "", "Path to the file to read")
	cryptomode := flag.String("crypto-mode", "", "(E)ncrypt or (D)ecrypt?")
	text := flag.String("text", "", "Text to parse")
	key := flag.String("key", "", "Cryptographic secret")
	flag.Parse()

	var parseString string
	switch *openmode {
	case "F":
		if *filename != "" {
			// Read the file contents
			fileContents, err := ioutil.ReadFile(*filename)
			if err != nil {
				fmt.Println("Error reading file:", err)
				return
			}
			// Convert file contents to string
			parseString = string(fileContents)
		} else {
			fmt.Println("No filename provided")
			return
		}

	case "T":
		if *text != "" {
			parseString = *text
		} else {
			fmt.Println("No text provided")
			return
		}
	default:
		fmt.Println("Invalid open mode")
		return
	}

	if flag.CommandLine.Lookup("key") == nil || *key == "" {
		fmt.Println("No key provided.")
		flag.Usage()
		return
	}

	if flag.CommandLine.Lookup("crypto-mode") == nil || *cryptomode == "" {
		fmt.Println("No crypto-mode specified.")
		flag.Usage()
		return
	}
	// endregion errorchecks

	switch *cryptomode {
	case "E":
		encrypted := encrypt(parseString, *key)
		fmt.Printf("%s \n", encrypted)
	case "D":
		decrypted := decrypt(parseString, *key)
		fmt.Printf("%s \n", decrypted)
	default:
		fmt.Println("Invalid input received. Try again.")
	}
}

func encrypt(stringToEncrypt string, keyString string) (encryptedString string) {

	//Since the key is in string, we need to convert decode it to bytes
	key, _ := hex.DecodeString(keyString)
	plaintext := []byte(stringToEncrypt)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	//Create a new GCM - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	//https://golang.org/pkg/crypto/cipher/#NewGCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	//Create a nonce. Nonce should be from GCM
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	//Encrypt the data using aesGCM.Seal
	//Since we don't want to save the nonce somewhere else in this case, we add it as a prefix to the encrypted data. The first nonce argument in Seal is the prefix.
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return fmt.Sprintf("%x", ciphertext)
}

func decrypt(encryptedString string, keyString string) (decryptedString string) {

	key, _ := hex.DecodeString(keyString)
	enc, _ := hex.DecodeString(encryptedString)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	//Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	//Get the nonce size
	nonceSize := aesGCM.NonceSize()

	//Extract the nonce from the encrypted data
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	//Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}

	return fmt.Sprintf("%s", plaintext)
}

// StringPrompt asks for a string value using the label
func StringPrompt(label string) string {
	var s string
	r := bufio.NewReader(os.Stdin)
	for {
		_, _ = fmt.Fprint(os.Stderr, label+" ")
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}
	return strings.TrimSpace(s)
}
