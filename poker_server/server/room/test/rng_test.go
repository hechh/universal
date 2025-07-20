package test

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"math/big"
	"testing"
)

func TestUUid(t *testing.T) {

	uuid2 := uuid.New()
	uuid3 := uuid.New()
	t.Logf("uuid2: %v, uuid3: %v", uuid2, uuid3)
}

func TestSecureRNG(t *testing.T) {
	// soft rng
	rng := NewSecureRNG()
	if err := rng.VerifyEntropySource(); err != nil {
		panic("RNG Init err: " + err.Error())
	}

	auditedRNG := &AuditedRNG{
		rng:    rng,
		secret: []byte("your-hmac-secret"),
	}
	key := make([]byte, 8)
	auditedRNG.Read(key)

	randRange, err := rng.Range(1, 9999)
	if err != nil {
		fmt.Printf("%v", err)
	}
	fmt.Printf("made [1, 9999) rand number : %d \n", randRange)

}

type SecureRNG struct{}

func NewSecureRNG() *SecureRNG {
	return &SecureRNG{}
}

func (r *SecureRNG) Read(buf []byte) (int, error) {
	if _, err := rand.Read(buf); err != nil {
		return 0, err
	}
	return len(buf), nil
}

// Uint64 made [0, 2^64-1) rand number
func (r *SecureRNG) Uint64() (uint64, error) {
	buf := make([]byte, 8)
	if _, err := rand.Read(buf); err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint64(buf), nil
}

// Uint64 made []byte
func (r *SecureRNG) Seed() ([]byte, error) {
	buf := make([]byte, 8)
	if _, err := rand.Read(buf); err != nil {
		return nil, err
	}
	return buf, nil
}

// Range made [min, max) rand number
func (r *SecureRNG) Range(min, max int64) (int64, error) {
	if min >= max {
		return 0, errors.New("invalid range")
	}
	bigMax := big.NewInt(max - min)
	n, err := rand.Int(rand.Reader, bigMax)
	if err != nil {
		return 0, err
	}
	return min + n.Int64(), nil
}

func (r *SecureRNG) VerifyEntropySource() error {
	// 尝试读取1字节验证熵源
	buf := make([]byte, 1)
	if _, err := rand.Read(buf); err != nil {
		return errors.New("System Entropy Source Cant use: " + err.Error())
	}
	return nil
}

type AuditedRNG struct {
	rng    io.Reader
	secret []byte // HMAC密钥
}

func (a *AuditedRNG) Read(p []byte) (n int, err error) {
	n, err = a.rng.Read(p)
	if err == nil {
		h := hmac.New(sha256.New, a.secret)
		h.Write(p)
		log.Printf("RNG Made: size=%d, HMAC=%x", n, h.Sum(nil))
	}
	return
}
