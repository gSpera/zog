package speccy

import (
	"fmt"
	"time"

	"github.com/jbert/zog"
)

type Machine struct {
	printState speccyPrintState
	screen     *Screen
	z          *zog.Zog

	done chan struct{}
}

func NewMachine(z *zog.Zog) *Machine {
	screen, err := NewScreen(z.Mem)
	if err != nil {
		panic(fmt.Sprintf("Can't create screen: %s", err))
	}

	return &Machine{
		z:      z,
		screen: screen,
		done:   make(chan struct{}),
	}
}

func (m Machine) LoadAddr() uint16 {
	return 0x8000
}

func (m Machine) RunAddr() uint16 {
	return 0x0000
}

func (m Machine) Name() string {
	return "speccy"
}

func (m *Machine) Start() error {
	err := m.loadROMs()
	if err != nil {
		return err
	}
	InstallKeyboardInputPorts(m.z)
	every := time.Second / 50
	go func() {
		tick := time.Tick(every)
		for {
			select {
			case <-m.done:
				break
			case <-tick:
				m.screen.Draw()
				m.z.DoInterrupt()
			}
		}
	}()

	return nil
}

func (m *Machine) Stop() {
	close(m.done)
}

const romFileName = "/usr/share/spectrum-roms/48.rom"

func (m *Machine) loadROMs() error {
	loadRealROM := true
	if loadRealROM {
		return m.z.LoadROMFile(0x0000, romFileName)
	} else {
		return m.loadConsolePrintROMs()
	}
}

func (m *Machine) loadConsolePrintROMs() error {
	m.z.RegisterOutputHandler(0xffff, m.printState.speccyPrintByte)

	// We only use RST 16
	zeroPageAssembly, err := zog.Assemble(`
	ORG 0000h
	HALT
	NOP
	NOP
	NOP
	NOP
	NOP
	NOP
	NOP
	NOP
	NOP
	NOP
	NOP
	NOP
	NOP
	NOP
	NOP
	; One entry point at 10h (RST 16), to print char in A
	PUSH DE
	LD E, A
	call printchar
	POP DE
	RET
` + printAssembly)
	if err != nil {
		return fmt.Errorf("Failed to assemble prelude: %s", err)
	}
	err = m.z.Load(zeroPageAssembly)
	if err != nil {
		return fmt.Errorf("Load zero page assembly: %s", err)
	}

	chanOpenAssembly, err := zog.Assemble(`
	ORG 1601h
	RET
`)
	if err != nil {
		return fmt.Errorf("Failed to assemble chan-open: %s", err)
	}
	err = m.z.Load(chanOpenAssembly)
	if err != nil {
		return fmt.Errorf("Load chan open assembly: %s", err)
	}

	return nil
}
