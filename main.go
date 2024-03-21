package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Uso: go run main.go <arquivo_pastas> <nome_branch>")
		return
	}

	arquivoPastas := os.Args[1]
	nomeBranch := os.Args[2]

	pastas, err := lerPastas(arquivoPastas)
	if err != nil {
		log.Fatal(err)
	}

	totalBranchesCriadas := 0

	for _, pasta := range pastas {
		criou, err := criarBranch(pasta, nomeBranch)
		if err != nil {
			log.Printf("Erro ao criar branch na pasta %s: %v\n", pasta, err)
			continue
		}
		if criou {
			fmt.Printf("Branch '%s' criada na pasta %s\n", nomeBranch, pasta)
			totalBranchesCriadas++
		}
	}

	fmt.Printf("Total de branches criadas: %d\n", totalBranchesCriadas)
}

func lerPastas(arquivoPastas string) ([]string, error) {
	file, err := os.Open(arquivoPastas)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var pastas []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		pasta := strings.TrimSpace(scanner.Text())
		if pasta != "" {
			pastas = append(pastas, pasta)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return pastas, nil
}

func criarBranch(pasta, nomeBranch string) (bool, error) {
	// Verifica se a pasta já é um repositório git
	_, err := os.Stat(filepath.Join(pasta, ".git"))
	if err != nil {
		if os.IsNotExist(err) {
			return false, fmt.Errorf("a pasta %s não é um repositório git", pasta)
		}
		return false, err
	}

	// Verifica se está na branch master
	cmd := exec.Command("git", "-C", pasta, "rev-parse", "--abbrev-ref", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return false, err
	}
	currentBranch := strings.TrimSpace(string(out))
	if currentBranch != "master" {
		cmd = exec.Command("git", "-C", pasta, "checkout", "master")
		err = cmd.Run()
		if err != nil {
			return false, err
		}
	}

	// Verifica se a branch já existe
	cmd = exec.Command("git", "-C", pasta, "rev-parse", "--verify", nomeBranch)
	err = cmd.Run()
	if err == nil {
		// Branch já existe
		return false, nil
	}

	// Cria a branch
	cmd = exec.Command("git", "-C", pasta, "checkout", "-b", nomeBranch)
	err = cmd.Run()
	if err != nil {
		return false, err
	}

	return true, nil
}
