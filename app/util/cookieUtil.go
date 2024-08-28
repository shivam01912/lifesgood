package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

var ErrInvalidValue = errors.New("invalid admin cookie value")

func SetCookieHandler(w http.ResponseWriter, r *http.Request, username string, value string) {
	// Initialize a new cookie containing the string "Hello world!" and some
	// non-default attributes.
	cookie := http.Cookie{
		Name:     "__session",
		Value:    username + "#" + value,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}

	writeEncrypted(w, &cookie, []byte(os.Getenv("AES_Key")))

	// Use the http.SetCookie() function to send the cookie to the client.
	// Behind the scenes this adds a `Set-Cookie` header to the response
	// containing the necessary cookie data.
	http.SetCookie(w, &cookie)

	w.Header().Set("Cache-Control", "private")

	r.AddCookie(&cookie)
}

func writeEncrypted(w http.ResponseWriter, cookie *http.Cookie, secretKey []byte) {
	// Create a new AES cipher block from the secret key.
	block, err := aes.NewCipher(secretKey)
	if err != nil {
		log.Println("Error creating cipher : ", err)
		return
	}

	// Wrap the cipher block in Galois Counter Mode.
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		log.Println("Error creating block : ", err)
		return
	}

	// Create a unique nonce containing 12 random bytes.
	nonce := make([]byte, aesGCM.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		log.Println("Error creating nonce : ", err)
		return
	}

	// Prepare the plaintext input for encryption. Because we want to
	// authenticate the cookie name as well as the value, we make this plaintext
	// in the format "{cookie name}:{cookie value}". We use the : character as a
	// separator because it is an invalid character for cookie names and
	// therefore shouldn't appear in them.
	plaintext := fmt.Sprintf("%s:%s", cookie.Name, cookie.Value)

	// Encrypt the data using aesGCM.Seal(). By passing the nonce as the first
	// parameter, the encrypted data will be appended to the nonce — meaning
	// that the returned encryptedValue variable will be in the format
	// "{nonce}{encrypted plaintext data}".
	encryptedValue := base64.URLEncoding.EncodeToString(aesGCM.Seal(nonce, nonce, []byte(plaintext), nil))

	// Set the cookie value to the encryptedValue.
	cookie.Value = encryptedValue
}

func ValidateCookie(r *http.Request, cookieName string) bool {
	secretKey := []byte(os.Getenv("AES_Key"))
	value, err := readEncrypted(r, cookieName, secretKey)
	if err != nil {
		log.Fatal("Error : ", err)
		return false
	}

	expectedUsername, passHash, ok := strings.Cut(value, "#")
	if !ok {
		return false
	}

	credentials := FetchAdminCredentials(expectedUsername)
	if credentials.Password != passHash {
		return false
	}

	return true
}

func readEncrypted(r *http.Request, name string, secretKey []byte) (string, error) {
	// Read the encrypted value from the cookie as normal.
	log.Println("Reading cookie with name : ", name)
	log.Println("Request : ", r)
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}

	encryptedValue, err := base64.URLEncoding.DecodeString(cookie.Value)
	if err != nil {
		return "", err
	}

	// Create a new AES cipher block from the secret key.
	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", err
	}

	// Wrap the cipher block in Galois Counter Mode.
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Get the nonce size.
	nonceSize := aesGCM.NonceSize()

	// To avoid a potential 'index out of range' panic in the next step, we
	// check that the length of the encrypted value is at least the nonce
	// size.
	if len(encryptedValue) < nonceSize {
		return "", ErrInvalidValue
	}

	// Split apart the nonce from the actual encrypted data.
	nonce := encryptedValue[:nonceSize]
	ciphertext := encryptedValue[nonceSize:]

	// Use aesGCM.Open() to decrypt and authenticate the data. If this fails,
	// return a ErrInvalidValue error.
	plaintext, err := aesGCM.Open(nil, []byte(nonce), []byte(ciphertext), nil)
	if err != nil {
		return "", ErrInvalidValue
	}

	// The plaintext value is in the format "{cookie name}:{cookie value}". We
	// use strings.Cut() to split it on the first ":" character.
	expectedName, value, ok := strings.Cut(string(plaintext), ":")
	if !ok {
		return "", ErrInvalidValue
	}

	// Check that the cookie name is the expected one and hasn't been changed.
	if expectedName != name {
		return "", ErrInvalidValue
	}

	// Return the plaintext cookie value.
	return value, nil
}
