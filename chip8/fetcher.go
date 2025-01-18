package chip8

const (
	programStartMemoryAddress = 0x200
	instructionBytes          = 2
)

type Fetcher struct {
	counter uint16
}

func NewFetcher() *Fetcher {
	return &Fetcher{
		counter: programStartMemoryAddress,
	}

}

func (f *Fetcher) Fetch(c *Chip8) uint16 {
	instruction := c.memory.ReadWord(f.counter)
	f.incrementCounter()
	return instruction
}

func (f *Fetcher) Skip() {
	f.incrementCounter()
}

func (f *Fetcher) incrementCounter() {
	// TODO: handle PC > 4096 / 0x1000 (12 bits). 2 options: PC overflow (error) or wrap (modulo)
	// TODO: also check PC < 521 / 0x200 (program start)
	f.counter += instructionBytes
}

func (f *Fetcher) GetCounter() uint16 {
	return f.counter
}

func (f *Fetcher) SetCounter(counter uint16) {
	f.counter = counter
}
