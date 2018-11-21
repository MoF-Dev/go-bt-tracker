package main

import (
	"fmt"
	"github.com/MoF-Dev/go-bt-tracker/pkg/bencode"
)

func main() {
	fmt.Printf("Hello, %s\n", "World")
	xd := bencode.Integer(56)
	fmt.Println(xd.Encode())
	fmt.Println(bencode.ByteString("huehuehue").Encode())
	dx := bencode.List(make([]bencode.BValue, 3))
	dx[0] = xd
	dx[1] = bencode.ByteString("ggnore")
	dx[2] = bencode.List(make([]bencode.BValue, 0))
	fmt.Println(dx.Encode())
	gg := bencode.Dictionary(make(map[string]bencode.BValue))
	gg["hue"] = bencode.ByteString("lmao")
	gg["wew"] = bencode.Integer(69)
	gg["list"] = dx
	fmt.Println(gg.Encode())
}
