package main

import (
	"fmt"
	"os"
	"os/exec"
)

func F(X, Y, Z uint32) uint32 {
	return (X & Y) | (^X & Z)
}
func G(X, Y, Z uint32) uint32 {
	return (X & Z) | (Y & ^Z)
}
func H(X, Y, Z uint32) uint32 {
	return X ^ Y ^ Z
}
func I(X, Y, Z uint32) uint32 {
	return Y ^ (X | ^Z)
}

func rotateLeft(x, n uint32) uint32 {
	return (x << n) | (x >> (32 - n))
}

func byteToUintArr(buf []byte) []uint32 {
	words := make([]uint32, len(buf)/4)

	for i := 0; i < len(buf); i += 4 {
		words[i/4] += uint32(buf[i+0]) << 0
		words[i/4] += uint32(buf[i+1]) << 8
		words[i/4] += uint32(buf[i+2]) << 16
		words[i/4] += uint32(buf[i+3]) << 24
	}

	return words
}

func rawMD5ToHEX(value uint32) string {
	res := ""
	for i := 0; i < 4; i++ {
		res += fmt.Sprintf("%02X", value%256)
		value /= 256
	}

	return res
}

func clearScr(osName string) {
	var cmd *exec.Cmd

	switch osName {
	case "linux":
		cmd = exec.Command("clear")
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default:
		fmt.Println()
		return
	}

	cmd.Stdout = os.Stdout
	cmd.Run()
}
