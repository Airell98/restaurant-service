package helpers

import (
	"fmt"
	"math/rand"
)

func GenerateBaseSerial(prefix string) string {
	length := 6
  
    ran_str := make([]byte, length)
  
	for i := 0; i < length; i++ {
        ran_str[i] = byte(65 + rand.Intn(25))
    }
  
	serial := string(ran_str)

	result := fmt.Sprintf("%s-%s", prefix, serial)

    return  result
}