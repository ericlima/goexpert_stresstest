package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Estrutura para armazenar os resultados do teste
type Resultados struct {
	totalRequests       int
	status200           int
	outrosStatus        map[int]int
	tempoTotalExecucao  time.Duration
	mu                  sync.Mutex
}

func main() {
	// Parâmetros de entrada via CLI
	url := flag.String("url", "", "URL do serviço a ser testado (obrigatório)")
	requests := flag.Int("requests", 0, "Número total de requests (obrigatório)")
	concurrency := flag.Int("concurrency", 1, "Número de chamadas simultâneas (obrigatório)")
	flag.Parse()

	// Validação dos parâmetros
	if *url == "" || *requests <= 0 || *concurrency <= 0 {
		fmt.Println("Uso: --url=<URL> --requests=<número total> --concurrency=<concorrência>")
		return
	}

	fmt.Printf("Iniciando teste de carga para URL: %s\n", *url)
	fmt.Printf("Número total de requests: %d\n", *requests)
	fmt.Printf("Concorrência: %d\n", *concurrency)

	// Inicialização dos resultados
	resultados := &Resultados{
		totalRequests:       *requests,
		outrosStatus:        make(map[int]int),
	}

	// Canal para sincronização das goroutines
	wg := &sync.WaitGroup{}
	requestsCh := make(chan struct{}, *concurrency)

	// Registro do tempo de início
	inicio := time.Now()

	// Lançar requests
	for i := 0; i < *requests; i++ {
		wg.Add(1)
		requestsCh <- struct{}{} // Controle de concorrência

		go func() {
			defer wg.Done()
			fazerRequest(*url, resultados)
			<-requestsCh
		}()
	}

	// Aguarda a conclusão de todos os requests
	wg.Wait()

	// Registro do tempo total
	resultados.tempoTotalExecucao = time.Since(inicio)

	// Exibe o relatório final
	imprimirRelatorio(resultados)
}

// Função para realizar um request HTTP
func fazerRequest(url string, resultados *Resultados) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Erro ao realizar request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Registro do status HTTP retornado
	resultados.mu.Lock()
	if resp.StatusCode == 200 {
		resultados.status200++
	} else {
		resultados.outrosStatus[resp.StatusCode]++
	}
	resultados.mu.Unlock()
}

// Função para exibir o relatório final
func imprimirRelatorio(resultados *Resultados) {
	fmt.Println("\n--- Relatório de Teste de Carga ---")
	fmt.Printf("Tempo total de execução: %v\n", resultados.tempoTotalExecucao)
	fmt.Printf("Total de requests realizados: %d\n", resultados.totalRequests)
	fmt.Printf("Total de requests com status 200: %d\n", resultados.status200)
	fmt.Println("Distribuição de outros status HTTP:")
	for status, count := range resultados.outrosStatus {
		fmt.Printf("  Status %d: %d\n", status, count)
	}
	fmt.Println("-----------------------------------")
}
