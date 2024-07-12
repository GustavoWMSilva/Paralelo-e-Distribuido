// LORENZO MORE
// GUSTAVO SILVA

package main

import (
	"fmt"
	"sync"
)

type mensagem struct {
	tipo  int    // tipo da mensagem para fazer o controle do que fazer
	corpo [4]int // conteudo da mensagem para colocar os IDs
}

var (
	chans = []chan mensagem{ // vetor de canais para formar o anel de eleicao
		make(chan mensagem),
		make(chan mensagem),
		make(chan mensagem),
		make(chan mensagem),
	}
	controle = make(chan int)
	wg       sync.WaitGroup // wg e usado para esperar o programa terminar
)

func ElectionControler(in chan int) {
	defer wg.Done()

	var temp mensagem

	// comandos para o anel iniciam aqui

	// mudar o processo 0 - canal de entrada 3 - para falho
	temp.tipo = 2
	chans[3] <- temp
	fmt.Printf("\nControle: mudar o processo 0 para falho\n")
	fmt.Printf("Controle: confirmacao de quem falhou: %d\n", <-in) // confirmacao

	// pedindo para processo 1 convocar eleicao
	temp.tipo = 3
	chans[0] <- temp
	fmt.Printf("\nControle: processo 1, convoque nova eleicao\n")
	fmt.Printf("Controle: lider atual: %d\n", <-in) // confirmacao

	// mudar o processo 3 - canal de entrada 2 - para falho
	temp.tipo = 2
	chans[2] <- temp
	fmt.Printf("\nControle: mudar o processo 3 para falho\n")
	fmt.Printf("Controle: confirmacao de quem falhou: %d\n", <-in) // confirmacao

	// pedindo para processo 2 convocar eleicao
	temp.tipo = 3
	chans[1] <- temp
	fmt.Printf("\nControle: processo 2, convoque nova eleicao\n")
	fmt.Printf("Controle: lider atual: %d\n", <-in) // confirmacao

	// pedindo para processo 1 convocar eleicao
	temp.tipo = 3
	chans[0] <- temp
	fmt.Printf("\nControle: processo 2, convoque nova eleicao\n")
	fmt.Printf("Controle: lider atual: %d\n", <-in) // confirmacao

	// ativando processo 0
	temp.tipo = 3
	chans[3] <- temp
	fmt.Printf("\nControle: ativar o processo 0\n")
	fmt.Printf("Controle: lider atual: %d\n", <-in) // confirmacao

	// ativando processo 3
	temp.tipo = 3
	chans[2] <- temp
	fmt.Printf("\nControle: ativar o processo 3\n")
	fmt.Printf("Controle: lider atual: %d\n", <-in) // confirmacao

	// indicando finalizacao dos processos mandando para processo 1
	fmt.Printf("\nControle: finalizando processos um a um\n")
	temp.tipo = 7
	chans[0] <- temp
	chans[2] <- temp
	chans[3] <- temp
	chans[1] <- temp
}

func ElectionStage(TaskId int, in chan mensagem, out chan mensagem, leader int) {
	defer wg.Done()

	// variaveis locais que indicam se este processo e o lider e se esta ativo

	var actualLeader int
	var bFailed bool = false // todos iniciam sem falha

	actualLeader = leader // indicacao do lider que foi passada como parametro

	for {
	temp := <-in // ler mensagem
	if temp.tipo == 3 || temp.tipo == 4 {
		fmt.Printf("Processo %2d: recebi mensagem %d, [ %d, %d, %d, %d ]\n", TaskId, temp.tipo, temp.corpo[0], temp.corpo[1], temp.corpo[2], temp.corpo[3])
	}
		
	switch temp.tipo {
	// torna falho
	case 2:
	{
		bFailed = true
		fmt.Printf("Processo %2d: falhei, enviando mensagem pro controle\n", TaskId)
		//fmt.Printf("Processo %2d: lider atual %d\n", TaskId, leader)
		leader = -5
		controle <- TaskId
	}
	// volta como era antes
	case 3:
	{
		if bFailed == true {
			fmt.Printf("Processo %2d: estou ativo novamente, chamando nova eleicao\n", TaskId)
		}
		// processo ativo
		bFailed = false
		leader = -5
		temp.tipo = 4
		for i := range temp.corpo {
			temp.corpo[i] = -5 // desativa todos
		}
		temp.corpo[TaskId] = TaskId
		fmt.Printf("Processo %2d: colocando meu id na mensagem: [ %d, %d, %d, %d ]\n", TaskId, temp.corpo[0], temp.corpo[1], temp.corpo[2], temp.corpo[3])
		out <- temp
	}
	// colocar ID no corpo
	case 4:
	{
		if temp.corpo[TaskId] == TaskId {
			// Percorra o array a partir do segundo elemento (índice 1) e compare cada elemento com o valor máximo
			fmt.Printf("Processo %2d: recebi todas ids: [ %d, %d, %d, %d ]\n", TaskId, temp.corpo[0], temp.corpo[1], temp.corpo[2], temp.corpo[3])
			for i := 0; i < len(temp.corpo); i++ {
				if temp.corpo[i] > temp.corpo[0] {
					temp.corpo[0] = temp.corpo[i]
				}
			}

			leader = temp.corpo[0]
			temp.corpo[1] = TaskId
			temp.tipo = 5
			fmt.Printf("Processo %2d: vencedor eleicao: processo %d \n", TaskId, leader)
			fmt.Printf("Processo %2d: informando novo lider aos outros processos\n", TaskId)
			fmt.Printf("Processo %2d: lider atualizado %d\n", TaskId, leader)
		} else if bFailed == false {
			temp.corpo[TaskId] = TaskId
			fmt.Printf("Processo %2d: colocando meu id na mensagem: [ %d, %d, %d, %d ]\n", TaskId, temp.corpo[0], temp.corpo[1], temp.corpo[2], temp.corpo[3])
		} else {
			fmt.Printf("Processo %2d:  Processo falho, pula\n", TaskId)
		}
		out <- temp
	}
	// mandando ID do novo lider para todos
	case 5:
	{
		if TaskId != temp.corpo[1] {
			if bFailed == false {
				leader = temp.corpo[0]
				fmt.Printf("Processo %2d: lider atualizado: %d \n", TaskId, leader)
			} else {
				fmt.Printf("Processo %2d: Processo falho, pula\n", TaskId)
			}
		} else {
			temp.tipo = 6
			controle <- leader
		}
		out <- temp

	}
	case 6:
	{
		// parte para sair do laco
		break
	}

	// terminando processos
	case 7:
	{
		if bFailed == true {
			fmt.Printf("Processo %2d ja foi encerrado.\n", TaskId)
			return
		} else {
			fmt.Printf("Processo %2d: mensagem para finalizar recebida.\n", TaskId)
			bFailed = true
			return	
		}
	}
	default:
	{
		fmt.Printf("Processo %2d: nao conheco este tipo de mensagem\n", TaskId)
		fmt.Printf("Processo %2d: lider atual %d\n", TaskId, actualLeader)
	}
	}
}
fmt.Printf("%2d: terminei \n", TaskId)
}

func main() {

	wg.Add(5) // Adicione uma contagem de quatro, um para cada goroutine

	// criar os processos do anel de eleição

	go ElectionStage(0, chans[3], chans[0], 0) // este e o lider
	go ElectionStage(1, chans[0], chans[1], 0) // nao e lider, e o processo 0
	go ElectionStage(2, chans[1], chans[2], 0) // nao e lider, e o processo 0
	go ElectionStage(3, chans[2], chans[3], 0) // nao e lider, e o processo 0

	fmt.Println("\n   Anel de processos criado")

	// criar o processo controlador

	go ElectionControler(controle)

	fmt.Println("\n   Processo controlador criado")

	wg.Wait() // Esperar pelas goroutines terminarem
}
