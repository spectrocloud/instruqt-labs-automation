package util

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"

	log "github.com/spectrocloud/gomi/pkg/logger"
)

func RandomInt64(max int64) int64 {
	if max <= 0 {
		return 0
	}

	var randomInt int64
	err := binary.Read(rand.Reader, binary.BigEndian, &randomInt)
	if err != nil {
		log.Error("Error in generating random int : %v", err)
		return 0
	}

	randomInt = int64(math.Abs(float64(randomInt)))
	randomInt = randomInt % max

	return randomInt
}

func GenerateRandomUint64() uint64 {
	var b [8]byte
	// Try generating a random uint64 up to 3 times
	for i := 0; i < 3; i++ {
		_, err := rand.Read(b[:])
		if err != nil {
			continue
		}
		return binary.LittleEndian.Uint64(b[:])
	}
	//error handling
	return 0
}

func GenerateRandomUInt64(max *big.Int) uint64 {

	randomNum := GenerateRandom(max)
	if randomNum == nil {
		return 0 // Default value when errors occur
	}
	// Successfully generated a random number, return it
	return randomNum.Uint64()
}

func GenerateRandomInt64(max *big.Int) int64 {

	randomNum := GenerateRandom(max)
	if randomNum == nil {
		return 0 // Default value when errors occur
	}
	// Successfully generated a random number, return it
	return randomNum.Int64()
}

func GenerateRandom(max *big.Int) *big.Int {

	if max.Int64() == 0 {
		return big.NewInt(0)
	}
	// Try generating a random uint64 up to 3 times
	for i := 0; i < 3; i++ {
		randomNum, err := rand.Int(rand.Reader, max)
		if err != nil {
			// Log the error but continue trying (retrying up to 3 times)
			log.Error(fmt.Sprintf("Error generating random number (attempt %d): %v", i+1, err))
			continue
		}

		// Successfully generated a random number, return it
		return randomNum
	}

	// If all retries fail, return a default value (for example, 0)
	log.Info("Failed to generate a random number after 3 attempts, returning default value.")
	return nil // Default value when errors occur
}
