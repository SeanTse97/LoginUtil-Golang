package main

import (
	"crypto/rsa"
	"encoding/hex"
	"math/big"
)

func encryptRSA(pub *rsa.PublicKey, data []byte) []byte {
	encrypted := new(big.Int)
	e := big.NewInt(int64(pub.E))
	payload := new(big.Int).SetBytes(data)
	encrypted.Exp(payload, e, pub.N)
	return encrypted.Bytes()
}

func getEncrypPassword(password string) string {
	var publicKey = &rsa.PublicKey{
		N: new(big.Int),
	}
	var plainText, encrypted []byte

	bytes := []rune(password)
	for from, to := 0, len(bytes)-1; from < to; from, to = from+1, to-1 {
		bytes[from], bytes[to] = bytes[to], bytes[from]
	}
	str := string(bytes) // Reverse string of user password

	plainText = []byte(str)

	var bignum, _ = new(big.Int).SetString("9c2899b8ceddf9beafad2db8e431884a79fd9b9c881e459c0e1963984779d6612222cee814593cc458845bbba42b2d3474c10b9d31ed84f256c6e3a1c795e68e18585b84650076f122e763289a4bcb0de08762c3ceb591ec44d764a69817318fbce09d6ecb0364111f6f38e90dc44ca89745395a17483a778f1cc8dc990d87c3", 16)
	publicKey.N = bignum
	publicKey.E = 0x10001

	encrypted = encryptRSA(publicKey, plainText)
	encodedStr := hex.EncodeToString(encrypted)
	return encodedStr
}
