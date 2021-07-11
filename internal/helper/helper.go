package helper

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/cespare/xxhash"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/sessions"
	"gitlab.com/pt-mai/maihelper"
	"gitlab.com/pt-mai/maihelper/maisession"
)

func Decrypt(data []byte, passphrase string) []byte {
	key := []byte(CreateHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return plaintext
}

func CreateHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

// SetSession func
func SetSession(w http.ResponseWriter, r *http.Request, sessionName string, sesVal string, reqBody interface{}) (err error) {
	// Getting Session Store
	store, err := maisession.Store()
	if err != nil {
		log.Println("Error When Getting Session Store: ", err.Error())

		maihelper.GrpcClient.MaiHttpResponseHandler(
			w,
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			"Failed Getting Session Store. Please contact Admin!",
			nil,
		)
		return
	}

	defer store.Close()

	// Run a background goroutine to clean up expired sessions from the database.
	defer store.StopCleanup(store.Cleanup(time.Minute * 5))

	// Getting Session
	session, err := store.Get(r, sessionName)
	if err != nil {
		log.Println("Error When Getting Session: ", err.Error())

		maihelper.GrpcClient.MaiHttpResponseHandler(
			w,
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			"Failed Getting Session. Please contact Admin!",
			nil,
		)
		return
	}

	// Add a value
	session.Values[sesVal] = reqBody

	// Save
	err = session.Save(r, w)
	if err != nil {
		log.Println("Error When Saving Session: ", err.Error())

		maihelper.GrpcClient.MaiHttpResponseHandler(
			w,
			http.StatusInternalServerError,
			"error",
			"Failed Saving Session. Please contact Admin!",
			nil,
		)
		return
	}

	return
}

// GetSession func
func GetSession(w http.ResponseWriter, r *http.Request, sessionName string) (sess *sessions.Session, err error) {
	// Getting Session Store
	store, err := maisession.Store()
	if err != nil {
		log.Println("Error When Getting Session Store: ", err.Error())
		return nil, maihelper.GrpcServer.MaiErrorDetail("session", "Failed Getting Session Store. Please contact Admin!", nil)
	}

	defer store.Close()

	// Run a background goroutine to clean up expired sessions from the database.
	defer store.StopCleanup(store.Cleanup(time.Minute * 5))

	// Getting Session
	session, err := store.Get(r, sessionName)
	if err != nil {
		log.Println("Error When Getting Session: ", err.Error())
		return nil, maihelper.GrpcServer.MaiErrorDetail("session", "Failed Getting Session. Please contact Admin!", nil)
	}

	return session, nil
}

// GenerateAuthCookie func
func GenerateAuthCookie(w http.ResponseWriter, email string) {
	authExpMinutes, _ := strconv.Atoi(os.Getenv("AUTH_EXP_MINUTE"))
	expirationTime := time.Now().Add(time.Duration(authExpMinutes) * time.Minute)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": email,
		"nbf": time.Now().Unix(),
		"exp": expirationTime.Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenStr, err := token.SignedString([]byte(os.Getenv("AUTH_SECRET")))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     os.Getenv("AUTH_COOKIE_NAME"),
		Value:    tokenStr,
		Expires:  expirationTime,
		HttpOnly: true,
		Path:     "/",
		Secure:   os.Getenv("GO_ENV") != "development",
	})
}

func GenerateDeviceCookie(rw http.ResponseWriter, phone string) string {
	// Expired cookie
	deviceExpMinutes, _ := strconv.Atoi(os.Getenv("AUTH_EXP_MINUTE_DEVICE"))
	expirationTime := time.Now().Add(time.Duration(deviceExpMinutes) * time.Minute)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": phone,
		"nbf": time.Now().Unix(),
		"exp": expirationTime.Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenStr, err := token.SignedString([]byte(os.Getenv("AUTH_SECRET_DEVICE")))
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return ""
	}

	http.SetCookie(rw, &http.Cookie{
		Name:     os.Getenv("AUTH_COOKIE_DEVICE"),
		Value:    tokenStr,
		Expires:  expirationTime,
		HttpOnly: true,
		Path:     "/",
		Secure:   os.Getenv("GO_ENV") != "development",
	})

	return tokenStr
}

func SecretKeyGen(skHash string) (ret string) {
	r1 := RandInt(1, 15)
	g1 := skHash[0:r1] + "/" + skHash[r1+1:]
	//log.Println(r1, len(g1), g1)

	r2 := RandInt(1, 15)
	g2 := g1[0:r2] + "+" + g1[r2+1:]
	//log.Println(r2, len(g2), g2)

	//r3 := RandInt(1, 15)
	//g3 := g2[0:r3] + "_" + g2[r3+1:]

	ts := strconv.FormatInt(time.Now().UnixNano(), 10)
	skRe := g2 + ts
	//log.Println(len(skRe),skRe)
	skHashRe := HashGenerate(skRe)
	//log.Println(len(skHashRe), skHashRe)

	skHashF := g2[0:9] + skHashRe + g2[9:16]
	//log.Println(len(skHashF),skHashF)
	skHashF = skHashF[16:32] + skHashF[0:16]
	return skHashF
}

func RandInt(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

func HashGenerate(str string) (ret string) {
	ret = fmt.Sprintf("%X", xxhash.Sum64String(str))
	if len(ret) < 16 {
		startIdx := RandInt(0, 8)
		ret = ret + ret[startIdx:startIdx+(16-len(ret))]
	}
	return
}
