// Curso: CO.N1.15
// Disciplina: Sistemas Operacionais
// Exercício "b" da questão "2"
// Grupo (
//     Carlos Gabriel
//     Isaque Almeida
//     Luis Fernando
// )
//
// Testes comparativos:
//
// RAM: 16GB DDR4 2800Mhz
// CPU: i5 6600K @4.6Ghz (4 Threads)
// SO : Fedora 34 (Kernel 5.11.17-Liquorix)
//
// 10 testes efetuados com N = 20, média de 728.66ms usando o binário
//
//
// Instalar o Golang:
// 		[Debian, Ubuntu e derivados]: sudo apt install golang
// 		[Fedora]: sudo dnf install golang ou sudo dnf install golang-bin
// 		[Outras distros]: 'https://golang.org/doc/install'
// 		[Windows]: Instalação padrão de um .exe
//
// 		Rodar o programa diretamente: go run nome_do_arquivo.go
// 		Gerar um binario: go build nome_do_arquivo.go
//
// Enunciado:
// "Dado um vetor com N números inteiros aleatórios (intervalo [0 .. 1000]),
// gerar o maior número de combinações de M números (M < (N/2)) em um arquivo em disco;"

package main

import (
	"bufio"     // Buffers do sistema.
	"context"   // Contextos.
	"fmt"       // Prints e scans.
	"math/rand" // Números aleatórios.
	"os"        // I/O do sistema.
	"strconv"   // Conversão de tipos.
	"strings"   // Manipulação de strings.
	"sync"      // Sincrônia de threads.
	"time"      // Marcação de tempo.
)

// Variáveis Globais
var (
	wg sync.WaitGroup

	n    int // Quantidade de itens no vetor a ser permutado.
	qtde int // Quantidade de combinações geradas.
)

// Constantes Globais.
const (
	// Maior numero gerado pela função rand.
	maxRand int = 1000

	// Formato da data e hora
	// dd-mm-yyyy hh-mm-ss.ms
	layout string = "02-01-2006 T 15-04-05.000000"
)

// Cria o tipo 'arr', que neste caso é apenas um array to tipo string
// Criado para organizar o código e permitir trocar o tipo de array
// de forma mais simples entre todas as funções que o usa, apenas
// trocando o tipo dele na definição.
type arr []string

// Atribui um valor ao vetor que chamou a função.
func (array *arr) add(i int, value string) {
	(*array)[i] = value
}

// Retorna o tamanho do vetor(arr) que chamou a função.
func (array *arr) len() (length int) {
	length = len(*array)
	return
}

// Retorna o valor do índice do vetor(arr) que chamou a função.
func (array *arr) get(i int) (get string) {
	get = (*array)[i]
	return
}

// Gera e insere os números aleatórios no array
func (array *arr) gen() {
	var r string
	for i := 0; i < n; i++ {
		r = strconv.Itoa(rand.Intn(maxRand + 1))
		array.add(i, r)
	}
}

// Inicia a geração das combinações
func genCombinations(array *arr, m int, results chan<- arr, ctx context.Context) {
	combinations(array, m, 0, "", results, ctx)

	// Fecha o canal results, como a função 'combinations' não está em uma thread separada,
	// ela usa a mesma da 'genCombinations', o canal só irá ser fechado ao fim da função anterior
	// seja ele por finalizar as combinações ou por um sinal
	defer close(results)
}

// Função que gera as combinações.
//
// Recebe '6' argumentos onde:
// array: Recebe o array a ser combinado.
// m: Tamanho do conjunto a ser gerado.
// start: Id de ínicio do array (Para uso especifico da recursividade)
// prefix: Grupo de elementos para se juntar a permutação dos demais.
// results: Canal de comunicação, onde recebe cada combinação gerada, nesse caso
// a função recebe o canal no modo de apenas escrita.
// ctx: Contexto de parada, usado interromper a função caso receba um sinal.
func combinations(array *arr, m, start int, prefix string, results chan<- arr, ctx context.Context) {
	// Para a função caso receba o sinal
	select {
	case <-ctx.Done():
		return
	default:
	}

	// Para a função caso já tenha terminado todas as combinações
	if m == 0 {
		return
	} else {
		// Itera pelo array, usando o indice incial como 'start' para que combine diferentes
		// elementos
		for i := start; i <= array.len()-m; i++ {
			// Chamada recursiva, subtraindo 1 ao m, e somando 1 ao i(start), para que diminua a
			// quantidade de elementos a serem "permutados".
			combinations(array, m-1, i+1, prefix+" "+array.get(i), results, ctx)
			// Verifica se a combinação gerada não está repetida
			if len([]string{prefix + " " + array.get(i)}) == m {
				// Remove os espaços a esquerda do prefixo
				woS := strings.TrimLeft(prefix, "\t \n")
				var rp = []string{woS + " " + array.get(i)}
				// Envia para o canal a combinação
				results <- rp
				// Soma 1 a quantidade de combinações geradas, na variável global 'qtde'
				qtde++
			}
		}
	}
}

// Gera o nome do arquivo usando a data e hora.
func genFileName(t *time.Time) string {
	filename := "2B-" + t.Format(layout)
	return filename + ".txt"
}

// Recebe uma variável do tipo 'error' e verifica se possui algum.
// Caso haja algum erro o programa é interrompido.
// Função criada para aumentar a legibilidade do código, já que esse processo
// poderia ser efetuado na própria função que a chamou.
func errCheck(err error) {
	if err != nil {
		panic(err)
	}
}

// Cria e escreve no arquivo as combinações geradas conforme são geradas pela função 'comb'.
// Recebe um ponteiro para a hora de início do processamento.
// Recebe os resultados das combinações pelo canal 'results' em modo apenas leitura.
func write(t *time.Time, array *arr, results <-chan arr, notWait context.CancelFunc) {
	// Cria o arquivo atribuindo ele à variável 'fo' (File Output), e a saída de algum erro
	// a variável 'err'.
	fo, err := os.Create(genFileName(t))

	// Verifica se possui algum erro
	errCheck(err)

	// Cria uma promessa, para que se feche o arquivo quando a função for finalizada.
	// Caso o 'defer' não seja usado o arquivo iria fechar antes que as informações
	// fossem inseridas.
	defer fo.Close()

	// Insere a hora de inicio.
	fo.WriteString(fmt.Sprintf("ínicio = %s\n\n", t))

	// Insere o vetor base.
	fo.WriteString(fmt.Sprintf("Vetor base = %s\n\n", (*array)))

	// Insere todos as combinações geradas.
	fo.WriteString(fmt.Sprintf("Combinações de %d para %d elementos:\n", n/2, n))
	for c := range results {
		fo.WriteString(fmt.Sprintln(c))
	}

	// Marca o tempo do fim da escrita.
	te := time.Now()

	// Escreve a hora do fim, quantidade de permutações e o tempo levado.
	fo.WriteString(fmt.Sprintf("\ntermino = %s\n", te))
	fo.WriteString(fmt.Sprintf("qtde = %d\n", qtde))
	fo.WriteString(fmt.Sprintf("tempo = %s", te.Sub(*t)))

	// Dispara um sinal indicando que o programa já terminou
	notWait()

	// Finaliza a thread(goroutine).
	defer wg.Done()
}

// Inicia um array
func newArr() *arr {
	var array arr
	array = make([]string, n)
	return &array
}

// Limpa o buffer de teclado
func flushStdin() {
	stdin := bufio.NewReader(os.Stdin)
	stdin.ReadString('\n')
}

func main() {
	// Armazena na variável global 'n' o tamanho do vetor a ser combinado.
	fmt.Scanf("%d", &n)

	// Marca o tempo de ínicio do programa (Após a primeira entrada do usuário).
	t := time.Now()

	// Gera uma seed para a geração de números aleatórios.
	rand.Seed(time.Now().UnixNano())

	// Tamanho das combinações
	m := n / 2

	// Cria as variáveis de contexto, que serão usadas para iterromper a função que gera
	// as combinações.
	ctx, cancel := context.WithCancel(context.Background())
	writeCtx, notWait := context.WithCancel(context.Background())

	// Cria o canal para comunicação dos resultados entre as funções.
	results := make(chan arr)

	// Cria o vetor a ser combinado.
	var array = newArr()
	// Gera os números aleatórios no vetor.
	array.gen()

	// Inicia uma thread(goroutine) para a função de combinação.
	// Ela não entra no WaitGroup por ser uma função recursiva.
	go genCombinations(array, m, results, ctx)

	wg.Add(1)
	// Inicia uma thread(goroutine) para a função que escreve os resultados no arquivo.
	go write(&t, array, results, notWait)

	// Cria uma thread(goroutine), para finalizar o programa em diferentes situações
	// a thread finaliza caso receba o sinal da função 'write' informando que já acabou de
	// escrever tudo no arquivo ou se o usuário pressionar ENTER.
	go func() {
		// O select é uma junção de 'do while' e 'switch case', ele fica em um loop
		// até que um dos casos finalize o loop.
		select {
		case <-writeCtx.Done(): // Verifica o sinal da função 'write'
			break
		default:
			flushStdin() // Limpa o buffer do teclado
			fmt.Scanln() // Espera o usuário pressionar ENTER
			// Aciona o sinal de parada para a função 'combinations', porém o programa só irá terminar
			// após a função 'write' terminar de escrever no arquivo tudo que existe no canal 'results'
			cancel()
		}
	}()

	// Faz a main aguardar a função write() até que todos os resultados sejam escritos
	// no arquivo.
	wg.Wait()
}
