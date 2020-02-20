package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"runtime"
)

var (
	A, B, C, D uint32
	T          [64]uint32
)

func main() {
	var md5 string
	in := bufio.NewScanner(os.Stdin)

	for {
		if len(os.Args) != 2 {
			fmt.Print("String to hash: ")
			in.Scan()

			md5 = calcMD5([]byte(in.Text()))
		} else {
			file, err := ioutil.ReadFile(os.Args[1])
			if err != nil {
				log.Fatalln(err)
			}

			md5 = calcMD5(file)

			fmt.Printf("File: %s\n", os.Args[1])
		}

		fmt.Printf("md5:  %s\n", md5)

		in.Scan()

		clearScr(runtime.GOOS)
	}
}

func calcMD5(buf []byte) string {
	dataLen := uint64(len(buf) * 8)

	appendPaddingBytes(&buf)
	appendLength(&buf, dataLen)
	initMDbuf()
	processMsgIn16WordBlocks(byteToUintArr(buf))

	return rawMD5ToHEX(A) + rawMD5ToHEX(B) + rawMD5ToHEX(C) + rawMD5ToHEX(D)
}

func appendPaddingBytes(buf *[]byte) {
	*buf = append(*buf, 128)

	for len(*buf)%64 != 56 {
		*buf = append(*buf, 0)
	}
}

func appendLength(buf *[]byte, length uint64) {
	buffer := bytes.NewBuffer(nil)

	if err := binary.Write(buffer, binary.LittleEndian, length); err != nil {
		log.Fatalln(err)
	}

	for _, b := range buffer.Bytes() {
		*buf = append(*buf, b)
	}
}

func initMDbuf() {
	A = 0x67452301
	B = 0xEFCDAB89
	C = 0x98BADCFE
	D = 0x10325476

	for i := 0; i < 64; i++ {
		T[i] = uint32(math.Pow(2, 32) * math.Abs(math.Sin(float64(i+1))))
	}
}

func processMsgIn16WordBlocks(buf []uint32) {
	for n := 0; n < len(buf); n += 16 {
		AA, BB, CC, DD := A, B, C, D

		// Round 1

		A = B + rotateLeft((A+F(B, C, D)+buf[n+0]+T[0]), 7)
		D = A + rotateLeft((D+F(A, B, C)+buf[n+1]+T[1]), 12)
		C = D + rotateLeft((C+F(D, A, B)+buf[n+2]+T[2]), 17)
		B = C + rotateLeft((B+F(C, D, A)+buf[n+3]+T[3]), 22)

		A = B + rotateLeft((A+F(B, C, D)+buf[n+4]+T[4]), 7)
		D = A + rotateLeft((D+F(A, B, C)+buf[n+5]+T[5]), 12)
		C = D + rotateLeft((C+F(D, A, B)+buf[n+6]+T[6]), 17)
		B = C + rotateLeft((B+F(C, D, A)+buf[n+7]+T[7]), 22)

		A = B + rotateLeft((A+F(B, C, D)+buf[n+8]+T[8]), 7)
		D = A + rotateLeft((D+F(A, B, C)+buf[n+9]+T[9]), 12)
		C = D + rotateLeft((C+F(D, A, B)+buf[n+10]+T[10]), 17)
		B = C + rotateLeft((B+F(C, D, A)+buf[n+11]+T[11]), 22)

		A = B + rotateLeft((A+F(B, C, D)+buf[n+12]+T[12]), 7)
		D = A + rotateLeft((D+F(A, B, C)+buf[n+13]+T[13]), 12)
		C = D + rotateLeft((C+F(D, A, B)+buf[n+14]+T[14]), 17)
		B = C + rotateLeft((B+F(C, D, A)+buf[n+15]+T[15]), 22)

		// Round 2

		A = B + rotateLeft((A+G(B, C, D)+buf[n+1]+T[16]), 5)
		D = A + rotateLeft((D+G(A, B, C)+buf[n+6]+T[17]), 9)
		C = D + rotateLeft((C+G(D, A, B)+buf[n+11]+T[18]), 14)
		B = C + rotateLeft((B+G(C, D, A)+buf[n+0]+T[19]), 20)

		A = B + rotateLeft((A+G(B, C, D)+buf[n+5]+T[20]), 5)
		D = A + rotateLeft((D+G(A, B, C)+buf[n+10]+T[21]), 9)
		C = D + rotateLeft((C+G(D, A, B)+buf[n+15]+T[22]), 14)
		B = C + rotateLeft((B+G(C, D, A)+buf[n+4]+T[23]), 20)

		A = B + rotateLeft((A+G(B, C, D)+buf[n+9]+T[24]), 5)
		D = A + rotateLeft((D+G(A, B, C)+buf[n+14]+T[25]), 9)
		C = D + rotateLeft((C+G(D, A, B)+buf[n+3]+T[26]), 14)
		B = C + rotateLeft((B+G(C, D, A)+buf[n+8]+T[27]), 20)

		A = B + rotateLeft((A+G(B, C, D)+buf[n+13]+T[28]), 5)
		D = A + rotateLeft((D+G(A, B, C)+buf[n+2]+T[29]), 9)
		C = D + rotateLeft((C+G(D, A, B)+buf[n+7]+T[30]), 14)
		B = C + rotateLeft((B+G(C, D, A)+buf[n+12]+T[31]), 20)

		// Round 3

		A = B + rotateLeft((A+H(B, C, D)+buf[n+5]+T[32]), 4)
		D = A + rotateLeft((D+H(A, B, C)+buf[n+8]+T[33]), 11)
		C = D + rotateLeft((C+H(D, A, B)+buf[n+11]+T[34]), 16)
		B = C + rotateLeft((B+H(C, D, A)+buf[n+14]+T[35]), 23)

		A = B + rotateLeft((A+H(B, C, D)+buf[n+1]+T[36]), 4)
		D = A + rotateLeft((D+H(A, B, C)+buf[n+4]+T[37]), 11)
		C = D + rotateLeft((C+H(D, A, B)+buf[n+7]+T[38]), 16)
		B = C + rotateLeft((B+H(C, D, A)+buf[n+10]+T[39]), 23)

		A = B + rotateLeft((A+H(B, C, D)+buf[n+13]+T[40]), 4)
		D = A + rotateLeft((D+H(A, B, C)+buf[n+0]+T[41]), 11)
		C = D + rotateLeft((C+H(D, A, B)+buf[n+3]+T[42]), 16)
		B = C + rotateLeft((B+H(C, D, A)+buf[n+6]+T[43]), 23)

		A = B + rotateLeft((A+H(B, C, D)+buf[n+9]+T[44]), 4)
		D = A + rotateLeft((D+H(A, B, C)+buf[n+12]+T[45]), 11)
		C = D + rotateLeft((C+H(D, A, B)+buf[n+15]+T[46]), 16)
		B = C + rotateLeft((B+H(C, D, A)+buf[n+2]+T[47]), 23)

		// Round 4

		A = B + rotateLeft((A+I(B, C, D)+buf[n+0]+T[48]), 6)
		D = A + rotateLeft((D+I(A, B, C)+buf[n+7]+T[49]), 10)
		C = D + rotateLeft((C+I(D, A, B)+buf[n+14]+T[50]), 15)
		B = C + rotateLeft((B+I(C, D, A)+buf[n+5]+T[51]), 21)

		A = B + rotateLeft((A+I(B, C, D)+buf[n+12]+T[52]), 6)
		D = A + rotateLeft((D+I(A, B, C)+buf[n+3]+T[53]), 10)
		C = D + rotateLeft((C+I(D, A, B)+buf[n+10]+T[54]), 15)
		B = C + rotateLeft((B+I(C, D, A)+buf[n+1]+T[55]), 21)

		A = B + rotateLeft((A+I(B, C, D)+buf[n+8]+T[56]), 6)
		D = A + rotateLeft((D+I(A, B, C)+buf[n+15]+T[57]), 10)
		C = D + rotateLeft((C+I(D, A, B)+buf[n+6]+T[58]), 15)
		B = C + rotateLeft((B+I(C, D, A)+buf[n+13]+T[59]), 21)

		A = B + rotateLeft((A+I(B, C, D)+buf[n+4]+T[60]), 6)
		D = A + rotateLeft((D+I(A, B, C)+buf[n+11]+T[61]), 10)
		C = D + rotateLeft((C+I(D, A, B)+buf[n+2]+T[62]), 15)
		B = C + rotateLeft((B+I(C, D, A)+buf[n+9]+T[63]), 21)

		A += AA
		B += BB
		C += CC
		D += DD
	}
}
