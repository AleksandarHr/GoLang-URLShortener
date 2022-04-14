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

// ShortenURL takes as input a long URL and performs a shortening algorithm.
// First, it appends to the end of the original URL stringified unix timestamp.
// This creates additional entropy which should be fairly large to avoid collision
// in short link generation.
// Then, it hashes the resulting string by computings its sha256. The hash is used
// to derive a large integer.
// Lastly, it computes the base58 of the said integer and returns the first 5 chars.
//
// TODO: Incrementally shorten URLs to 1, 2, 3, 4, and 5 characters as we use up all combinations
// of the corresponding length.
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
	base58Encoding := base58.BitcoinEncoding
	base58Encoded, err := base58Encoding.Encode(numberBytes)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Return the first 5
	shortUrl := string(base58Encoded)[:5]
	return shortUrl
}
