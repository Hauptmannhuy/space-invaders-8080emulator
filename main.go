package main

func main() {
	b := loadHex()
	pc := 0
	length := len(b)
	for pc < length {
		n := disassebmle(b, pc)
		pc += n
	}

}
