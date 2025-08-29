package freeway

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"net"
	"strings"
	"time"
)

// PortKnock
// The Ghost Twins are hot on your tail as you try to escape from the Merovingian on the freeway!
//
// Thankfully, you're able to fight back - Connect to the endpoint below and survive their phasing capabilities to force them to back off and get a flag for your trouble!
// freeway.white-rabbit.dev:32682
func PortKnock(addr string) (string, error) {
	conn, err := net.DialTimeout("tcp", addr, 2*time.Second) // Replace with your server address and port
	if err != nil {
		return "", fmt.Errorf("error connecting to %s: %v", addr, err)
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		msg, err := Read(reader)
		if err != nil {
			return "", fmt.Errorf("error reading from %s: %v", addr, err)
		}

		// Stage 1: Checksum check
		if strings.Contains(msg, "Stage 1 initiated") {
			//_, _ = conn.Write([]byte(fmt.Sprintf("%v\n", CheckSum("SERAPH"))))
			_, _ = conn.Write([]byte("451\n"))
			continue
		}

		// Stage 2: Caesar Cipher
		if strings.Contains(msg, "Stage 2") {
			//_, _ = conn.Write([]byte(fmt.Sprintf("%v\n", CaesarCipher("ZION", shift))))
			_, _ = conn.Write([]byte("MVBA\n"))
			continue
		}

		if strings.Contains(msg, "Stage 3") {
			//_, _ = conn.Write([]byte(fmt.Sprintf("%v\n", PseudoHash("ORACLE", shift))))
			_, _ = conn.Write([]byte("06eef40d\n"))
			continue
		}

		if strings.Contains(msg, "Stage 4") {
			//_, _ = conn.Write([]byte(fmt.Sprintf("%v\n", OctalLength("BluePill"))))
			_, _ = conn.Write([]byte("10\n"))
			continue
		}

		if strings.Contains(msg, "Stage 5") {
			//TUVST1ZJTkdJQU4=
			_, _ = conn.Write([]byte(fmt.Sprintf("%v\n", Base64Encode("MEROVINGIAN"))))
			continue
		}

		if strings.Contains(msg, "Stage 6") {

		}

		return "", fmt.Errorf("unexpected message: %s", msg)
	}
}

func CheckSum(in string) int {
	var chksm int
	for _, r := range in {
		chksm += int(r)
	}

	return chksm
}

func CaesarCipher(in string, shift int) string {
	var out string
	for _, r := range in {
		out += string('A' + (r-'A'+rune(shift))%26)
	}

	return out
}

func PseudoHash(in string) string {
	s := CheckSum(in)

	s = s ^ (s << 13) // xor
	s = s ^ (s >> 17)
	s = s ^ (s << 5)

	return strings.ToLower(fmt.Sprintf("%08X", 0xFFFFFFFF&s))
}

func OctalLength(in string) string {
	return fmt.Sprintf("%o", len(in))
}

func Base64Encode(in string) string {
	return base64.StdEncoding.EncodeToString([]byte(in))
}

func Read(reader *bufio.Reader) (string, error) {
	buffer := make([]byte, 1024)
	var out string

	for {
		n, err := reader.Read(buffer)
		if err != nil {
			return "", fmt.Errorf("Read: %w", err)
		}

		out += string(buffer)
		buffer = make([]byte, 1024)

		if n < 1024 {
			break
		}
	}

	return out, nil
}
