package algorithms

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"time"

	"github.com/itchyny/base58-go"
)

func ShortenURL(originalUrl string) string {
	// Get a unique url -- url + timestamp
	uniqueId := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	uniqueUrl := originalUrl + uniqueId

	// Get sha256 hash of the unique url
	hashAlgorithm := sha256.New()
	hashAlgorithm.Write([]byte(uniqueUrl))
	urlHash := hashAlgorithm.Sum(nil)

	// Generate big integer from the hash bytes
	urlHashNumber := new(big.Int).SetBytes(urlHash).Uint64()
	numberBytes := []byte(fmt.Sprintf("%d", urlHashNumber))

	// Encode the integer to base58 (easy to read/remember output)
	// TODO: add reasoning and resources for choosing base 58
	base58Encoding := base58.BitcoinEncoding
	base58Encoded, err := base58Encoding.Encode(numberBytes)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Return the first 5
	return string(base58Encoded)[:5]
}
