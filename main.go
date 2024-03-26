package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

	for i := 2; i <= 3; i++ {
		destino := fmt.Sprintf("%s-%d", origem, i)
		err := copiarPasta(origem, destino)
		if err != nil {
			fmt.Printf("Erro ao copiar pasta: %v\n", err)
			break
		}
		fmt.Printf("Pasta copiada e renomeada para %s\n", destino)

		err = substituirNomes(destino, origem, destino)
		if err != nil {
			fmt.Printf("Erro ao substituir nomes nos arquivos: %v\n", err)
		}
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

func substituirNomes(pasta, antigoNome, novoNome string) error {
	// Lista de arquivos onde faremos as substituições
	arquivos := []string{"pom.xml", "sonar-project.properties", "Jenkinsfile"}

	// Percorre os arquivos na pasta
	for _, nomeArquivo := range arquivos {
		caminhoArquivo := filepath.Join(pasta, nomeArquivo)

		// Verifica se o arquivo existe
		_, err := os.Stat(caminhoArquivo)
		if err != nil {
			if os.IsNotExist(err) {
				// Arquivo não existe, passa para o próximo
				continue
			}
			return err
		}

		// Abre o arquivo para leitura e cria um temporário para escrita
		arquivoOrigem, err := os.Open(caminhoArquivo)
		if err != nil {
			return err
		}
		defer arquivoOrigem.Close()

		arquivoTemp, err := os.Create(caminhoArquivo + ".temp")
		if err != nil {
			return err
		}
		defer arquivoTemp.Close()

		// Faz a substituição linha por linha
		scanner := bufio.NewScanner(arquivoOrigem)
		for scanner.Scan() {
			linha := scanner.Text()
			linha = strings.ReplaceAll(linha, antigoNome, novoNome)
			fmt.Fprintln(arquivoTemp, linha)
		}

		// Verifica se houve algum erro durante a leitura do arquivo original
		if err := scanner.Err(); err != nil {
			return err
		}

		// Fecha os arquivos
		arquivoOrigem.Close()
		arquivoTemp.Close()

		// Remove o arquivo original e renomeia o temporário para o nome original
		err = os.Remove(caminhoArquivo)
		if err != nil {
			return err
		}
		err = os.Rename(caminhoArquivo+".temp", caminhoArquivo)
		if err != nil {
			return err
		}
	}

	return nil
}
