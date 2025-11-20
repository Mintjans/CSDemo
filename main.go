package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	demoinfocs "github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs"
	events "github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs/events"
)

func criarMatrizEventos(p demoinfocs.Parser) {
	//Inicio da partida
	//Halftime
	//Round end type RoundEndOfficial
	//type RoundFreezetimeEnd ¶

}
func onKill(kill events.Kill) {
	var hs string
	if kill.IsHeadshot {
		hs = " (HS)"
	}

	var wallBang string
	if kill.PenetratedObjects > 0 {
		wallBang = " (WB)"
	}

	log.Printf("%s <%v%s%s> %s\n", kill.Killer, kill.Weapon, hs, wallBang, kill.Victim)
}

// writeCSV writes the provided matrix (rows of string fields) to the given file
// nomeArquivo: destination path for the CSV file
// matrix: [][]string where each inner slice is a CSV row
// returns an error on failure so callers can decide how to handle it
func writeCSV(nomeArquivo string, matrix [][]string) error {
	f, err := os.Create(nomeArquivo)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	// ensure file.Close() error is checked
	defer func() {
		if cerr := f.Close(); cerr != nil {
			log.Printf("close error: %v", cerr)
		}
	}()

	w := csv.NewWriter(f)

	for _, row := range matrix {
		if err := w.Write(row); err != nil {
			w.Flush()
			return fmt.Errorf("write error: %w", err)
		}
	}

	// ensure buffered data is flushed to disk
	w.Flush()

	// check for errors that may have occurred during flush/write
	if err := w.Error(); err != nil {
		return fmt.Errorf("csv writer error: %w", err)
	}

	return nil
}

func main() {
	err := demoinfocs.ParseFile("./demos/blast-rivals-2025-season-1-vitality-vs-spirit-bo3--Lq5eQCbcMsmLRiXVp3_m4/vitality-vs-spirit-m1-mirage.dem", func(p demoinfocs.Parser) error {
		p.RegisterEventHandler(onKill)
		//matrizJogador vai ter 2 colunas: idJogador e nomeJogador e vai guardar os dados dos jogadores da partida

		// criar matrizJogador com 11 linhas e 2 colunas
		matrizJogador := make([][]string, 11)
		for i := range matrizJogador {
			matrizJogador[i] = make([]string, 2)
		}
		// nomes das colunas (linha de cabeçalho)
		matrizJogador[0][0] = "idJogador"
		matrizJogador[0][1] = "nomeJogador"

		// preencher matrizJogador com os dados dos jogadores da partida
		players := p.GameState().Participants().Playing()
		for i, player := range players {
			matrizJogador[i+1][0] = fmt.Sprintf("%d", player.SteamID64) // idJogador
			matrizJogador[i+1][1] = player.Name                         // nomeJogador
		}
		log.Printf("players: %v\n", players)
		log.Printf("Matriz Jogador: %v\n", matrizJogador)

		// salvar matrizJogador em um arquivo CSV
		/*if err := writeCSV("matrizJogador.csv", matrizJogador); err != nil {
			log.Printf("erro ao gravar matrizJogador.csv: %v", err)
		}*/

		/*matrizEstadoJogo essa tabela vai servir para guardar os eventos que informam o estado da partida,
		se a partida esta pausada, se os times viraram de lado, se estamos no overtime,
		esse tipo de coisa e a estrutura é a seguinte: tick, tipoEstado, descricaoEstado.*/

		//do{} while {}
		/*for next, erro := p.ParseNextFrame(); next; next, erro = p.ParseNextFrame() {
			if erro != nil {
				log.Panic("fim: ", erro)
			} else {
				log.Printf("%v \n", p.CurrentFrame())
			}
		}*/

		return nil
	})
	if err != nil {
		log.Panic("failed to parse demo: ", err)
	}
}
