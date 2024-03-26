package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {

	origem := "dpsp-lmp-cornershop-carga-item"

	// Cria o diretório
	err := os.Mkdir("repositorios", 0755)
	if err != nil {
		fmt.Println("Erro ao criar o diretório:", err)
		return
	}

	errcd := os.Chdir("repositorios")
	if errcd != nil {
		fmt.Println("Erro ao mudar de diretório:", errcd)
		return
	}

	// Comando para criar um arquivo
	cmd3 := exec.Command("git", "clone", "http://daniel.crodrigues@bitbucket:7990/scm/lmp/"+origem+".git")
	if err := cmd3.Run(); err != nil {
		fmt.Println("Erro ao tentar clonar:", err)
		return
	}

	for i := 2; i <= 5; i++ {
		destino := fmt.Sprintf("%s-%d", origem, i)
		err := copiarPasta(origem, destino)
		if err != nil {
			fmt.Printf("Erro ao copiar pasta: %v\n", err)
			break
		}
		fmt.Printf("Pasta copiada e renomeada para %s\n", destino)
	}

	// Obter o diretório atual
	diretorioAtual, err := os.Getwd()
	if err != nil {
		fmt.Println("Erro ao obter o diretório atual:", err)
		return
	}

	fmt.Println("O diretório atual é:", diretorioAtual)
}

func copiarPasta(origem, destino string) error {
	// Cria a pasta de destino
	err := os.MkdirAll(destino, os.ModePerm)
	if err != nil {
		return err
	}

	// Percorre os arquivos da pasta de origem
	err = filepath.Walk(origem, func(caminhoOrigem string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Calcula o caminho de destino para o arquivo atual
		caminhoDestino := filepath.Join(destino, caminhoOrigem[len(origem):])

		if info.IsDir() {
			// Cria a pasta no destino
			err = os.MkdirAll(caminhoDestino, os.ModePerm)
			if err != nil {
				return err
			}
		} else {
			// Copia o arquivo para o destino
			arquivoOrigem, err := os.Open(caminhoOrigem)
			if err != nil {
				return err
			}
			defer arquivoOrigem.Close()

			arquivoDestino, err := os.Create(caminhoDestino)
			if err != nil {
				return err
			}
			defer arquivoDestino.Close()

			_, err = io.Copy(arquivoDestino, arquivoOrigem)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}
