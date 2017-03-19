package main

import (
	"fmt"
	"log"

	"github.com/jbert/zog"
)

func main() {
	memSize := uint16(16 * 1024)
	z := zog.New(memSize)
	instructions, err := z.Assembler.Assemble(`
	; test program for assembler
			LD HL, 0x0100		; and support end of line comments
			LD (HL), 0xFF
			LD B, (HL)
			LD (HL), 0x02
			LD A, (HL)
			ADD A, B
			HALT
		`)
	if err != nil {
		log.Fatalf("Failed to assemble: %s", err)
	}
	for _, i := range instructions {
		fmt.Printf("A: %v\n", i)
	}
	err = z.Encode(0, instructions)
	if err != nil {
		log.Fatalf("Failed to encode: %s", err)
	}

	/*
		z.Poke(0, 0x0E) // LD C, imm
		z.Poke(1, 0x11) // 0x11
		z.Poke(2, 0x3E) // LD A, imm
		z.Poke(3, 0x22) // 0x22
		z.Poke(4, 0x81) // ADD A, C
		z.Poke(5, 0x76) // HALT
	*/
	a, err := z.Run()
	if err != nil {
		log.Fatalf("Terminated: %s\n", err)
	}
	fmt.Printf("A is 0x%02X\n", a)
}
