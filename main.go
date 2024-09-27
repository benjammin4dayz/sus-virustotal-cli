package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		splash()
		os.Exit(0)
	}

	input := os.Args[1]

	info, err := os.Stat(input)
	if err != nil {
		log.Fatalf("Cannot find '%s'", input)
	}
	if info.IsDir() {
		log.Fatalf("Expected a file, but '%s' is a directory", input)
	}

	hash, err := hashFile(input)
	if err != nil {
		log.Fatalf("Failed to hash file located at '%s': %v", input, err)
	}

	url := fmt.Sprintf("https://virustotal.com/gui/file/%s", hash)
	if err := openBrowser(url); err != nil {
		log.Fatal(err)
	}
}

func splash() {
	fmt.Println(`
                           ...:::...
                        -*#%%%%%%%%%%#*+.
                       =@@@@@@@@@@@@@@@@@+
                      -@@@@@%#**++++*%@@@@#.
                     .@@@@+---=+*##*=--%@%@#
                     *@%@= -*#%%@@@@%# +@@%@=
                    :@@@@#:..::----::::%@@@@@.
                    *@@@@@@#*++====+*#@@@@@%@+
                   .@@@@@@@@@@@@@@@@@@@@@@@@@%
                   +@@@@@@@@@@@@@@@@@@@@@@@@@@-
                   %@@@@@@@@@@@@@@@@@@@@@@@@@@+
                  =@%@@@@@@@@@@@@@@@@@@@@@@@@@*
                  #@@@@@@@@@@@@@@@@@@@@@@@@@@@%
                 :@@@@@@@@@@@@@@@@@@@@@@@@@@@@%
                 *@@@@@@@@@@@@@@@@@@@@@@@@@@@@@.
                .@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@:
                =@%@@@@@@@@@@@@@@@@@@@@@@@@@@@@:
                #@@@@@@@%**++++++*%@@@@@@@@@@@@-
               .@@@@@@@#.          +@@@@@@@@@@@-
               +@@@@@@@:     ඞ     %@@@@@@@@@@-
         ...   %@@@@@@@-    sus     #@@@@@@@@@@:
      -##%%%%**@%@@@%@@-  v0.1.0    %@@@@@@@@@@.
      *@@@@@@@@@@@@@@@+            -@@@@@@@@@@%.
       :=+****##***+-.        .=++#@@@@@@@@@%@%
                             -@@@@@%%%@@@@@@@%-
                             :#%@@@@@@@@%#*+-
                               .:-=---::.`)

	hash, _ := hashFile(os.Args[0])
	fmt.Printf("Usage: %s <file_path>\nRepo: https://github.com/benjammin4dayz/sus-virustotal-cli\nHash: %s",
		strings.TrimSuffix(
			filepath.Base(os.Args[0]),
			filepath.Ext(os.Args[0]),
		),
		hash,
	)
}

func hashFile(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

func openBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
