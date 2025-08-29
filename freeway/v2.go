package freeway

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"net"
	"strings"
	"time"
)

func PortKnockV2(addr string) (string, error) {
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

		// Phase base64
		if strings.Contains(msg, "Phase: b64-encode") {
			out := Base64EncodeV2(msg)
			fmt.Printf("Phase base64-encode: %v\n", out)
			_, _ = conn.Write([]byte(fmt.Sprintf("%v\n", out)))
			continue
		}

		// Phase octal length
		if strings.Contains(msg, "Phase: len-octal") {
			out := OctalLenV2(msg)
			fmt.Printf("Phase octal-len: %v\n", out)
			_, _ = conn.Write([]byte(fmt.Sprintf("%v\n", out)))
			continue
		}

		// Phase reverse
		if strings.Contains(msg, "Phase: reverse") {
			out := ReverseV2(msg)
			fmt.Printf("Phase reverse: %v\n", out)
			_, _ = conn.Write([]byte(fmt.Sprintf("%v\n", out)))
			continue
		}

		// Phase: sum-gopher
		if strings.Contains(msg, "Phase: sum-gopher") {
			out := CheckSumV2(msg)
			fmt.Printf("Phase sum-gopher: %v\n", out)
			_, _ = conn.Write([]byte(fmt.Sprintf("%v\n", out)))
			continue
		}

		// Phase: consonant-count
		if strings.Contains(msg, "Phase: consonant-count") {
			out := CountConsonantsV2(msg)
			fmt.Printf("Phase consonant-count: %v\n", out)
			_, _ = conn.Write([]byte(fmt.Sprintf("%v\n", out)))
			continue
		}

		// Phase: len-binary
		if strings.Contains(msg, "Phase: len-binary") {
			out := BinaryLenV2(msg)
			fmt.Printf("Phase len-binary: %v\n", out)
			_, _ = conn.Write([]byte(fmt.Sprintf("%v\n", out)))
			continue
		}

		// Phase: unique-letters
		if strings.Contains(msg, "Phase: unique-letters") {
			out := UniqueLettersCountV2(msg)
			fmt.Printf("Phase unique-letters: %v\n", out)
			_, _ = conn.Write([]byte(fmt.Sprintf("%v\n", out)))
			continue
		}

		// Phase: vowel-count
		if strings.Contains(msg, "Phase: vowel-count") {
			out := VowelCountV2(msg)
			fmt.Printf("Phase vowel-count: %v\n", out)
			_, _ = conn.Write([]byte(fmt.Sprintf("%v\n", out)))
			continue
		}

		// Phase: xor-bytes.
		if strings.Contains(msg, "Phase: xor-bytes") {
			out := XORV2(msg)
			fmt.Printf("Phase xor-bytes: %v\n", out)
			_, _ = conn.Write([]byte(fmt.Sprintf("%v\n", out)))
			continue
		}

		// Phase: rot13-confs.
		if strings.Contains(msg, "Phase: rot13-confs") {
			out := ROT13V2(msg)
			fmt.Printf("Phase rot13-confs: %v\n", out)
			_, _ = conn.Write([]byte(fmt.Sprintf("%v\n", out)))
			continue
		}

		// Phase: sum-hex.
		if strings.Contains(msg, "Phase: sum-hex") {
			out := SumHexV2(msg)
			fmt.Printf("Phase: sum-hex: %v\n", out)
			_, _ = conn.Write([]byte(fmt.Sprintf("%v\n", out)))
			continue
		}

		// Phase: hex-pseudo.
		if strings.Contains(msg, "Phase: hex-pseudo") {
			out := PseudoHashV2(msg)
			fmt.Printf("Phase: hex-pseudo: %v\n", out)
			_, _ = conn.Write([]byte(fmt.Sprintf("%v\n", out)))
			continue
		}

		// Phase: unlock
		if strings.Contains(msg, "Ghost Twin solidified") {
			out := GetFlagV2(msg)
			fmt.Printf("Phase: get flag: %v\n", out)
			_, _ = conn.Write([]byte(fmt.Sprintf("%v\n", out)))
			continue
		}

		return "", fmt.Errorf("unexpected message: %s", msg)
	}
}

// Base64EncodeV2
//
//	Matrix shift. Phase: b64-encode.
//	Base64 encode 'MORPHEUS' (standard RFC 4648). Include '=' padding if required.
func Base64EncodeV2(msg string) string {
	return base64.StdEncoding.EncodeToString([]byte(GetQuoted(msg)))
}

// OctalLenV2
//
//	Matrix shift. Phase: len-octal.
//	Reply with the length of 'AGENTSMITH' in octal (base-8, no prefix).
func OctalLenV2(msg string) string {
	return fmt.Sprintf("%o", len(GetQuoted(msg)))
}

// ReverseV2
// Matrix shift. Phase: reverse.
// Reverse the string 'MORPHEUS' exactly.
func ReverseV2(msg string) string {
	var out []string
	for _, c := range GetQuoted(msg) {
		out = append([]string{string(c)}, out...)
	}

	return strings.Join(out, "")
}

// CountConsonantsV2
// Matrix shift. Phase: consonant-count.
// Count consonants in 'GHOST' (letters A–Z excluding vowels).
func CountConsonantsV2(msg string) string {
	var out string
	for _, c := range GetQuoted(msg) {
		if !strings.Contains("AEIOU", string(c)) {
			out += string(c)
		}
	}
	return fmt.Sprintf("%v", len(out))
}

// CheckSumV2
// Matrix shift. Phase: sum-gopher.
// Compute the ASCII checksum (decimal) for 'NEO' and return the number.
func CheckSumV2(msg string) int {
	var chksm int
	for _, r := range GetQuoted(msg) {
		chksm += int(r)
	}

	return chksm
}

// BinaryLenV2
// Matrix shift. Phase: len-binary.
// Reply with the length of 'KEYMAKER' in binary (base-2, no prefix).
func BinaryLenV2(msg string) string {
	return fmt.Sprintf("%b", len(GetQuoted(msg)))
}

// UniqueLettersCountV2
// Matrix shift. Phase: unique-letters.
// Count unique letters A–Z in 'NEO' (case-insensitive).
func UniqueLettersCountV2(msg string) string {
	var out string
	for _, c := range GetQuoted(msg) {
		if !strings.Contains(out, string(c)) {
			out += string(c)
		}
	}
	return fmt.Sprintf("%v", len(out))
}

// VowelCountV2
// Matrix shift. Phase: vowel-count.
// Count vowels in 'AGENTSMITH' (A,E,I,O,U).
func VowelCountV2(msg string) string {
	var out string
	for _, c := range GetQuoted(msg) {
		if strings.Contains("AEIOU", string(c)) {
			out += string(c)
		}
	}
	return fmt.Sprintf("%v", len(out))
}

// XORV2
// Matrix shift. Phase: xor-bytes.
// XOR all bytes of 'ORACLE' together (byte-wise) and return the decimal value (0-255).
func XORV2(msg string) int {
	var out int
	for _, c := range GetQuoted(msg) {
		out = out ^ int(c)
	}

	return out % 255
}

// ROT13V2
// Matrix shift. Phase: rot13-confs.
// Apply ROT13 to 'PERSEPHONE' and send the result.
func ROT13V2(msg string) string {
	return CaesarCipher(GetQuoted(msg), 13)
}

// SumHexV2
// Matrix shift. Phase: sum-hex.
// Compute the ASCII checksum for 'OSIRIS' and return it in lowercase hexadecimal (no prefix).
func SumHexV2(msg string) string {
	return fmt.Sprintf("%X", CheckSumV2(msg))
}

// PseudoHashV2
// Matrix shift. Phase: hex-pseudo.
// Pseudo-hash 'SERAPH': take its ASCII checksum, then mix with XOR shifts (v^=v<<13; v^=v>>17; v^=v<<5). Reply with the first 8 lowercase hex digits.
func PseudoHashV2(msg string) string {
	return PseudoHash(GetQuoted(msg))
}

// GetFlagV2
// Stabilized. Next phase...
// Ghost Twin solidified. Say 'UNLOCK' to claim the flag.
func GetFlagV2(msg string) string {
	return GetQuoted(msg)
}

func GetQuoted(msg string) string {
	strs := strings.Split(msg, "'")
	if len(strs) < 2 {
		return ""
	}

	return strs[1]
}
