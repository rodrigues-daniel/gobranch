package main

import (
	"fmt"
	"github.com/ktrysmt/go-bitbucket"
)

func main() {
	// Substitua 'seu-nome-de-usuario' e 'sua-senha' pelas credenciais reais
	client := bitbucket.NewBasicAuth("seu-nome-de-usuario", "sua-senha")

	// Faça uma solicitação para listar seus repositórios
	repos, err := client.Repositories.ListForAccount()
	if err != nil {
		fmt.Println("Erro ao listar repositórios:", err)
		return
	}

	// Imprima os nomes dos repositórios
	for _, repo := range repos {
		fmt.Println("Nome do Repositório:", repo.Name)
	}
}
