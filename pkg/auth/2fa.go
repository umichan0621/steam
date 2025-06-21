package auth

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
)

const (
	CHAR_SET     = "23456789BCDFGHJKMNPQRTVWXY"
	CHAR_SET_LEN = uint32(len(CHAR_SET))
)

func GenerateTwoFactorCode(sharedSecret string, current int64) (string, error) {
	data, err := base64.StdEncoding.DecodeString(sharedSecret)
	if err != nil {
		return "", err
	}

	ful := make([]byte, 8)
	binary.BigEndian.PutUint32(ful[4:], uint32(current/30))

	hash := hmac.New(sha1.New, data)
	_, err = hash.Write(ful)
	if err != nil {
		return "", err
	}

	sum := hash.Sum(nil)
	start := sum[19] & 0x0F
	slice := binary.BigEndian.Uint32(sum[start:start+4]) & 0x7FFFFFFF

	buf := make([]byte, 5)
	for i := 0; i < 5; i++ {
		buf[i] = CHAR_SET[slice%CHAR_SET_LEN]
		slice /= CHAR_SET_LEN
	}
	return string(buf), nil
}
