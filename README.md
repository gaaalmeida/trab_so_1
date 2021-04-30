# trab_so_1
Enunciados estão contidos em cada código-fonte.

### Como testar
1. Instale o Golang:
  - [Debian, Ubuntu e derivados]: `sudo apt install golang`
  - [Fedora]: `sudo dnf install golang` ou `sudo dnf install golang-bin`
  - [Outras distros]: [Golang docs](https://golang.org/doc/install)
  - [Windows]: Instalação padrão de um .exe
2. Iniciar diretamente: `go run nome_do_arquivo.go`
3. Gerar um binario: `go build nome_do_arquivo.go` e depois `./nome_do_arquivo`

### Grupo
- Carlos Gabriel
- Isaque Almeida
- Luis Fernando

### Testes
Média dos testes em um computador diferente, para comparação em diferentes hardwares.

Config:
- RAM: 16GB DDR4 2800Mhz
- CPU: i5 6600K @4.6Ghz (4 Threads)
- SO : Fedora 34 (Kernel 5.11.17-Liquorix)

Ambos os testes foram feitos com o binário!

- Ex 2a: 10 testes efeutados com N = 500, média de 869.19ms
- Ex 2b: 10 testes efetuados com N = 20, média de 728.66ms
