// Curso: CO.N1.15
// Disciplina: Sistemas Operacionais
// Exercício "a" da questão "2"
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
// 10 testes efeutados com N = 500, média de 869.19ms usando o binário
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
// "Dada duas matrizes quadradas de dimensão N x N composta por números inteiros aleatórios
// (intervalo [0 .. 1000]), construir um algoritmo paralelo para multiplica-las, gerando uma matriz
// resultado em um arquivo em disco;"

package main

import (
	"fmt"       // Prints e scans
	"math/rand" // Números aleatórios
	"os"        // I/O do sistema
	"sync"      // Sincrônia de threads
	"time"      // Marcação de tempo
)

// Variáveis Globais
var (
	wg       sync.WaitGroup
	n        int // Tamanho da matriz quadrada
	mapMutex = sync.RWMutex{}
)

// Constantes Globais.
const (
	maxRand        = 1000
	layout  string = "02-01-2006 T 15-04-05.000000"
)

// Cria o tipo 'arr', que neste caso é apenas um array to tipo string
// Criado para organizar o código e permitir trocar o tipo de array
// de forma mais simples entre todas as funções que o usa, apenas
// trocando o tipo dele na definição.
type arr struct {
	content map[int][]int64
}

func newArr() *arr {
	var array arr
	array.content = make(map[int][]int64)
	return &array
}

// Retorna o valor do índice do vetor(arr) que chamou a função.
func (array *arr) get(i, j int) (get int64) {
	get = array.content[i][j]
	return
}

// Atribui um valor ao vetor que chamou a função.
func (array *arr) add(i int, value []int64) {
	array.content[i] = value
}

// Gera e insere os números aleatórios no array
func (array *arr) gen(genChan, sizeChan chan bool) {
	<-sizeChan
	var tmp = []int64{}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			r := rand.Intn(maxRand + 1)
			tmp = append(tmp, int64(r))
		}
		array.add(i, tmp)
		tmp = nil
	}
	// Envia a confirmação para o canal indicando que a função já finalizou o processo
	// de gerar os números aleatórios.
	genChan <- true
	close(genChan)
	defer wg.Done()
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

// Gera o nome do arquivo usando a data e hora.
func genFileName(t *time.Time) string {
	filename := "2A-" + t.Format(layout)
	return filename + ".txt"
}

func write(t *time.Time, arrayA, arrayB, arrayF *arr) {
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
	fo.WriteString(fmt.Sprintf("%s\n\n", *t))

	// Insere a matriz A.
	fo.WriteString(fmt.Sprintf("MATRIZ A:\n"))
	for i := 0; i < n; i++ {
		fo.WriteString(fmt.Sprintf("%d\n", arrayA.content[i]))
	}

	// Insere a matriz B.
	fo.WriteString(fmt.Sprintf("\nMATRIZ B:\n"))
	for i := 0; i < n; i++ {
		fo.WriteString(fmt.Sprintf("%d\n", arrayB.content[i]))
	}

	// Insere o resultado da multiplicação entre a matriz A e B
	fo.WriteString(fmt.Sprintf("\nRESULTADO:\n"))
	for i := 0; i < n; i++ {
		fo.WriteString(fmt.Sprintf("%d\n", arrayF.content[i]))
	}

	// Marca o tempo do fim da escrita.
	te := time.Now()

	// Escreve a hora do fim e o tempo levado.
	fo.WriteString(fmt.Sprintf("\ntermino: %s\n", te))
	fo.WriteString(fmt.Sprintf("tempo: %s\n", te.Sub(*t)))
}

// Gera todas as threads para multiplicar as matrizes
func workers(done *int, jobs <-chan int, arrayA, arrayB, arrayF *arr) {
	// Inicia uma trhead para cada linha da matriz.
	for job := range jobs {
		go partial(job, done, arrayA, arrayB, arrayF)
	}
}

// Multiplica a matriz parcialmente, uma linha apenas.
// i: Corresponde ao número contido no canal 'jobs', a posição
// de ínicio da thread é igual a linha que ela irá multiplicar.
func partial(i int, done *int, arrayA, arrayB, arrayF *arr) {
	var tmp int64
	var arrTmp []int64
	for j := 0; j < n; j++ {
		for k := 0; k < n; k++ {
			tmp += arrayA.get(i, k) * arrayB.get(k, j)
		}
		arrTmp = append(arrTmp, tmp)
		tmp = 0
	}
	// Usa Mutex para impedir que multiplas threads insiram o array
	// no canal ao mesmo tempo, gerando um erro.
	mapMutex.Lock()
	arrayF.add(i, arrTmp)
	mapMutex.Unlock()

	// Soma 1 a quantidade de threads finalizadas
	*done++
}

func multiply(t *time.Time, arrayA, arrayB, arrayF *arr, multiplyChan, genChanA, genChanB chan bool) {
	// Aguarda a confirmação das duas funções que geram números aleatórios.
	<-genChanA
	<-genChanB

	// Canal usado para gerar as threads(goroutines).
	jobs := make(chan int)
	var done int

	// Cria uma thread(goroutine), que irá criar todas as outras threads que irão
	// multiplicar as matrizes.
	go workers(&done, jobs, arrayA, arrayB, arrayF)

	// Adiciona um no canal 'jobs' para cada linha da matriz.
	for i := 0; i < n; i++ {
		jobs <- i
	}
	// Fecha o canal
	close(jobs)

	// Espera todas as threads de multiplicação terminarem.
	for n > done {
	}

	// Chama a função que escreve os resultados no arquivo.
	write(t, arrayA, arrayB, arrayF)

	// Manda um sinal indicando que a função 'multiply' já finalizou.
	multiplyChan <- true

	// Finaliza a thread
	defer wg.Done()
}

func main() {
	// Armazena na variável global 'n' o tamanho da matriz a ser multiplicada.
	fmt.Scanf("%d", &n)

	// Marca o tempo de ínicio do programa (Após a primeira entrada do usuário).
	t := time.Now()

	// Gera uma seed para a geração de números aleatórios.
	rand.Seed(time.Now().UnixNano())

	// Inicia os dois canais usados, para enviar o sinal de quando as funções
	// que geram números aleatórios.
	genChanA := make(chan bool)
	genChanB := make(chan bool)

	// Canal usado pela função multiply para enviar o sinal que a função terminou.
	multiplyChan := make(chan bool)

	// Canal usado para confirmar que existe um tamanho para a matriz.
	sizeChan := make(chan bool)

	// Gera as matrizes base.
	var (
		arrayA = newArr()
		arrayB = newArr()
		arrayF = newArr()
	)

	// Adiciona 3 threads ao WaitGroup.
	wg.Add(3)

	// Gera os números aleatórios na matriz A e B.
	go arrayA.gen(genChanA, sizeChan)
	go arrayB.gen(genChanB, sizeChan)

	// Gera a thread da função que multiplica as matrizes.
	go multiply(&t, arrayA, arrayB, arrayF, multiplyChan, genChanA, genChanB)

	// Fecha o canal que passa a confirmação do tamanho da matriz.
	sizeChan <- true
	close(sizeChan)

	// Aguarda a função multiply acabar.
	<-multiplyChan

	// Faz a main aguardar as threads finalizarem.
	wg.Wait()
}
