package main

import (
	"crypto/md5"
	"encoding"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"hash"
	"slices"
	"strings"
)

var (
	secret_slider    = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	dorchadas_slider = "0,328,33297,6,0,B|48:323|61:274|61:274|45:207|45:207|63:169|103:169|103:169|249:199|249:199|215:214|205:254,1,450.000017166138,6|6,1:1|2:1,0:0:0:0:"
)

func sign(beatmap string) string {
	hsh := md5.Sum([]byte(secret_slider + beatmap))
	return hex.EncodeToString(hsh[:])
}

func verify(beatmap, signature string) bool {
	hsh := md5.Sum([]byte(secret_slider + beatmap))
	return hex.EncodeToString(hsh[:]) == signature
}

func has_dorchadas(beatmap string) bool {
	return strings.Contains(beatmap, dorchadas_slider)
}

func hashFromSum(hashStr string, ll int) hash.Hash {
	hash := md5.New()

	md5Bytes, err := hex.DecodeString(hashStr)
	if err != nil {
		panic(err)
	}

	if len(md5Bytes) != hash.Size() {
		panic("invalid md5 len")
	}

	state, err := hash.(encoding.BinaryMarshaler).MarshalBinary()
	if err != nil {
		panic(err)
	}

	slices.Reverse(md5Bytes[0 : 0+4])
	slices.Reverse(md5Bytes[4 : 4+4])
	slices.Reverse(md5Bytes[8 : 8+4])
	slices.Reverse(md5Bytes[12 : 12+4])

	copy(state[4:], md5Bytes)

	copy(state[len(state)-8:], binary.BigEndian.AppendUint64(nil, uint64(ll)))

	err = hash.(encoding.BinaryUnmarshaler).UnmarshalBinary(state)
	if err != nil {
		panic(err)
	}

	return hash
}

func main() {
	secretPadding := strings.Repeat("_", 64-len(secret_slider)%64)
	dorchadasWithPadding := dorchadas_slider + strings.Repeat("_", 64-len(dorchadas_slider)%64)

	getPaddingBuffer := func(ll int) []byte {
		tmp := [1 + 63 + 8]byte{0x80}
		pad := (55 - uint64(ll)) % 64
		binary.LittleEndian.PutUint64(tmp[1+pad:], (uint64(ll))<<3)
		fmt.Println("padding", ll, hex.EncodeToString(tmp[:1+pad+8]))
		return tmp[:1+pad+8]
	}

	fmt.Printf("secret_slider: %v %q\n", len(secret_slider), secret_slider)
	fmt.Println("")

	beatmap0 := secretPadding + strings.Repeat(".", 64)
	signature0 := sign(beatmap0)
	fmt.Printf("beatmap0 str: %q\n", beatmap0)
	fmt.Printf("beatmap0 base64: %q\n", base64.StdEncoding.EncodeToString([]byte(beatmap0)))
	fmt.Printf("has_dorchadas: %v\n", has_dorchadas(beatmap0))
	fmt.Printf("signature0: %v\n", signature0)
	fmt.Print("\n")

	{
		padding := string(getPaddingBuffer(len(secret_slider) + len(beatmap0)))
		beatmap1 := beatmap0 + padding + dorchadasWithPadding

		h := hashFromSum(signature0, len(secret_slider)+len(beatmap0)+len(padding))
		h.Write([]byte(dorchadasWithPadding))
		signature1 := hex.EncodeToString(h.Sum(nil)[:])

		fmt.Printf("beatmap1 str: %q\n", beatmap1)
		fmt.Printf("beatmap1 base64: %q\n", base64.StdEncoding.EncodeToString([]byte(beatmap1)))
		fmt.Printf("has_dorchadas: %v\n", has_dorchadas(beatmap1))
		fmt.Printf("signature1: %v\n", signature1)
		fmt.Printf("sign(): %v\n", sign(beatmap1))
		fmt.Printf("verify(): %v\n", verify(beatmap1, signature1))
	}
}
