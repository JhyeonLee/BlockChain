package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"io/fs"
	"math/big"
	"os"

	"github.com/JhyeonLee/BlockChain/utils"
)

const (
	fileName string = "blockchainCoin.wallet"
)

// for Interface
type fileLayer interface {
	hasWalletFile() bool
	writeFile(name string, data []byte, perm fs.FileMode) error
	readFile(name string) ([]byte, error)
}

type layer struct{}

func (layer) hasWalletFile() bool {
	_, err := os.Stat(fileName)
	return !os.IsNotExist(err) // wallet is existed
}

func (layer) writeFile(name string, data []byte, perm fs.FileMode) error {
	return os.WriteFile(name, data, perm)
}

func (layer) readFile(name string) ([]byte, error) {
	return os.ReadFile(fileName)
}

var files fileLayer = layer{}

// end, for interface

type wallet struct {
	privateKey *ecdsa.PrivateKey
	Address    string // public key (hexadecimal string)
}

// private key ~> sign
// public key ~> verify

var w *wallet

func createPrivKey() *ecdsa.PrivateKey {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.HandleErr(err)
	return privKey
}

func persistkey(key *ecdsa.PrivateKey) {
	bytes, err := x509.MarshalECPrivateKey(key)
	utils.HandleErr(err)
	err = files.writeFile(fileName, bytes, 0644) // 0644: permission to read and write
	utils.HandleErr(err)
}

func restoreKey() (key *ecdsa.PrivateKey) {
	keyAsBytes, err := files.readFile(fileName)
	utils.HandleErr(err)
	key, err = x509.ParseECPrivateKey(keyAsBytes)
	utils.HandleErr(err)
	return // key
}

func encodeBigInts(a, b []byte) string {
	z := append(a, b...)
	return fmt.Sprintf("%x", z)
}

func aFromK(key *ecdsa.PrivateKey) string {
	return encodeBigInts(key.X.Bytes(), key.Y.Bytes())
}

func Sign(payload string, w *wallet) string {
	payloadAsB, err := hex.DecodeString(payload)
	utils.HandleErr(err)
	r, s, err := ecdsa.Sign(rand.Reader, w.privateKey, payloadAsB)
	utils.HandleErr(err)
	return encodeBigInts(r.Bytes(), s.Bytes())
}

func restoreBigInts(payload string) (*big.Int, *big.Int, error) {
	bytes, err := hex.DecodeString(payload)
	if err != nil {
		return nil, nil, err
	}
	firstHalfBytes := bytes[:len(bytes)/2]
	secondtHalfBytes := bytes[len(bytes)/2:]
	bigA, bigB := big.Int{}, big.Int{}
	bigA.SetBytes(firstHalfBytes)
	bigB.SetBytes(secondtHalfBytes)
	return &bigA, &bigB, nil
}

func Verify(signature, payload, address string) bool {
	// restore signature into big.Int
	r, s, err := restoreBigInts(signature)
	utils.HandleErr(err)
	// restore public Key into big.Int
	x, y, err := restoreBigInts(address)
	utils.HandleErr(err)
	publicKey := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}
	// payload into hexadecimal bytes
	payloadBytes, err := hex.DecodeString(payload)
	utils.HandleErr(err)
	// verify
	ok := ecdsa.Verify(&publicKey, payloadBytes, r, s)
	return ok
}

func Wallet() *wallet {
	if w == nil {
		w = &wallet{}
		// if has a wallet already, restore from file
		if files.hasWalletFile() {
			w.privateKey = restoreKey()
		} else {
			// if not, create private key and sace to file
			key := createPrivKey()
			persistkey(key)
			w.privateKey = key
		}
		w.Address = aFromK(w.privateKey)
	}

	return w
}

// 1) hash a message
// "Hello World" -> hash(x) -> "hashed message"

// 2) generate key pair
// KeyPair (privateKey, publicKey)
// will save private key to a file ~> wallet

// 3) sign the hash
// ("hashed message" + privateKey) -> "signature"

// private key -> you signed
// public key -> you verify the private key

// 4) verify
// ("hashed message" + "signature" + publicKey) -> true or false
// check this signature was created by using the privateKey
// when verifying the signature, they do not need private key

/* const (
	signature     string = "3ab8ad04f3a9e5b74836219c258262d68e8034b0ecf9308c0292549ad4e8f67a189539ed9487345b49505dc53df99ad7515094e8d802b25be11ea7c16e022f39"
	privateKey    string = "30770201010420135d879aa0934f8492fb7cf2a1a7d7acc42d913ae5a7c1a7310b33a5b5d62fc8a00a06082a8648ce3d030107a144034200047efa7c330cde8a861cb395413464a5893275cd4e5e541d461349ff2094046b6076ead8269ec88bff63af8a0bde090f61a83c3ef6d920acace45c58334428a58c"
	hashedMessage string = "a591a6d40bf420404a011733cfb7b190d62c65bf0bcda32b57b277d9ad9f146e"
) */

// func Start() {
// algorithm elliptic.P256() for GeneratKey() is not same as algorithm for Bitcoiin
// because Go Standard Libray does not have Generate Key algorithm of Bitcoin, but the algorithm is similar
// privatekey struct has private key and public key struct
/*
	// Verygind Message: Genrate Keys, hashedMessage, Signature
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.HandleErr(err)
	keyAsBytes, err := x509.MarshalECPrivateKey(privateKey)
	utils.HandleErr(err)
	fmt.Printf("%x\n\n\n", keyAsBytes)
	message := "Hello World"
	hashedMessage := utils.Hash(message)
	hashAsByte, err := hex.DecodeString(hashedMessage)
	utils.HandleErr(err)
	fmt.Println(hashedMessage)
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hashAsByte)
	utils.HandleErr(err)
	signature := append(r.Bytes(), s.Bytes()...)
	fmt.Printf("\n\n\n\n%x\n", signature) */
// ok := ecdsa.Verify(&privateKey.PublicKey, hashAsByte, r, s)

/*
	// Restoring
	privByte, err := hex.DecodeString(privateKey)
	utils.HandleErr(err)
	private, err := x509.ParseECPrivateKey(privByte)
	utils.HandleErr(err)
	sigBytes, err := hex.DecodeString(signature)
	utils.HandleErr(err)
	rBytes := sigBytes[:len(sigBytes)/2]
	sBytes := sigBytes[len(sigBytes)/2:]
	var bigR, bigS = big.Int{}, big.Int{}
	bigR.SetBytes(rBytes)
	bigS.SetBytes(sBytes)
	hashBytes, err := hex.DecodeString(hashedMessage)
	utils.HandleErr(err)
	ok := ecdsa.Verify(&private.PublicKey, hashBytes, &bigR, &bigS)
	fmt.Println(ok) */
// }
