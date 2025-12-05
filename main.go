package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	demoinfocs "github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs"
	"github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs/common"
	events "github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs/events"
)

func criarMatrizJogadores(players []*common.Player) {
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
	for i, player := range players {
		matrizJogador[i+1][0] = fmt.Sprintf("%d", player.SteamID64) // idJogador
		matrizJogador[i+1][1] = player.Name                         // nomeJogador
	}

	// salvar matrizJogador em um arquivo CSV
	if err := writeCSV("dadosJogador.csv", matrizJogador); err != nil {
		log.Printf("erro ao gravar dadosJogador.csv: %v", err)
	}
}

func criarMatrizes(demoData demoinfocs.Parser) {
	/*matrizEventos essa tabela vai servir para guardar os eventos que informam o estado da partida,
	se a partida esta pausada, se os times viraram de lado, se estamos no overtime,
	esse tipo de coisa e a estrutura é a seguinte: tick, tipoEstado, descricaoEstado.*/

	/*criando matrizEventos com cabeçalho.
	tick: o tick do jogo em que o evento ocorreu, numero inteiro;
	tipoEstado: o tipo do estado da partida, vai ser um numero inteiro (depois devo aprender a fazer enum);
		- 1: inicio do round;
		- 2: freezetime end;
		- 3: round end;
		- 4: halftime;
		- 5: inicio da partida;
	nomeEstado: nome do evento, string.
	*/
	matrizEventos := [][]string{{"tick", "tipoEstado", "nomeEstado"}}

	// variável para guardar a lista de jogadores
	var player []*common.Player

	// variáveis para guardar o tick de inicio e fim da partida
	var tickinicio int
	var tickfim int

	/*parser.RegisterEventHandler example:
			parser.RegisterEventHandler(func(e events.WeaponFired) {
	    		fmt.Printf("%s fired his %s\n", e.Shooter.Name, e.Weapon.Type)
			})
	*/

	//Inicio da partida: MatchStart signals that the match has started.
	demoData.RegisterEventHandler(func(inicioPartida events.MatchStart) {
		// chamar a função criarMatrizJogadores ao inicio da partida, quando os jogadores já estão conectados no servidor
		criarMatrizJogadores(demoData.GameState().Participants().Playing())
		// guardar o tick de inicio da partida
		tickinicio = demoData.CurrentFrame()
		// guardar o evento de inicio da partida na matrizEventos
		matrizEventos = append(matrizEventos, []string{
			fmt.Sprintf("%d", demoData.CurrentFrame()), // tick
			"5",             // tipoEstado
			"inicioPartida", // nomeEstado
		})
	})

	//Halftime: TeamSideSwitch signals that teams are switching sides.
	demoData.RegisterEventHandler(func(halftime events.TeamSideSwitch) {
		matrizEventos = append(matrizEventos, []string{
			fmt.Sprintf("%d", demoData.CurrentFrame()), // tick
			"4",        // tipoEstado
			"halftime", // nomeEstado
		})
	})

	//RoundEnd: RoundEnd signals that a round just finished.
	demoData.RegisterEventHandler(func(roundEnd events.RoundEnd) {
		matrizEventos = append(matrizEventos, []string{
			fmt.Sprintf("%d", demoData.CurrentFrame()), // tick
			"3",        // tipoEstado
			"roundEnd", // nomeEstado
		})
	})

	//type RoundFreezetimeEnd: RoundFreezetimeEnd signals that the freeze time is over.
	demoData.RegisterEventHandler(func(freezetimeEnd events.RoundFreezetimeEnd) {
		matrizEventos = append(matrizEventos, []string{
			fmt.Sprintf("%d", demoData.CurrentFrame()), // tick
			"2",             // tipoEstado
			"freezetimeEnd", // nomeEstado
		})
	})

	//roundStart: RoundStart signals that a new round has started.
	demoData.RegisterEventHandler(func(roundStart events.RoundStart) {
		matrizEventos = append(matrizEventos, []string{
			fmt.Sprintf("%d", demoData.CurrentFrame()), // tick
			"1",          // tipoEstado
			"roundStart", // nomeEstado
		})
	})

	//AnnouncementWinPanelMatch: signals that the 'win panel' has been displayed. I guess that's the final scoreboard.
	// salvar matrizEventos em um arquivo CSV ao final da partida
	demoData.RegisterEventHandler(func(fim events.AnnouncementWinPanelMatch) {
		// guardar o tick de fim da partida
		tickfim = demoData.CurrentFrame()
		if err := writeCSV("estadoJogo.csv", matrizEventos); err != nil {
			log.Printf("erro ao gravar estadoJogo.csv: %v", err)
		}
	})

	// criar matrizPosições com cabeçalho: tick, idJogador, posiçãoX, posiçãoY, posiçãoZ.
	matrizPosicoes := [][]string{{"tick", "idJogador", "posiçãoX", "posiçãoY", "posiçãoZ"}}

	//do{} while {} loop para retirar todos dados de posição de cada frame
	for next, erro := demoData.ParseNextFrame(); next; next, erro = demoData.ParseNextFrame() {
		if erro != nil {
			log.Printf("erro ao parsear frame: %v", erro)
			break // ou continue, dependendo do caso
		}
		if !next {
			break
		}

		if tickfim != 0 {
			// já encerrado, opcionalmente sair do loop
			break
		}
		if tickinicio != 0 {
			// obter a lista de membros no servidor que estão jogando
			player = demoData.GameState().Participants().Playing()

			for _, p := range player {
				pos := p.Position()
				// adicionar a posição do jogador na matrizPosicoes
				matrizPosicoes = append(matrizPosicoes, []string{
					fmt.Sprintf("%d", demoData.CurrentFrame()), // tick
					fmt.Sprintf("%d", p.SteamID64),             // idJogador
					fmt.Sprintf("%.2f", pos.X),                 // posiçãoX
					fmt.Sprintf("%.2f", pos.Y),                 // posiçãoY
					fmt.Sprintf("%.2f", pos.Z),                 // posiçãoZ
				})
			}
		}
	}

	// salvar matrizPosicoes em um arquivo CSV
	if err := writeCSV("posicoesJogadores.csv", matrizPosicoes); err != nil {
		log.Printf("erro ao gravar posicoesJogadores.csv: %v", err)
	}
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
	err := demoinfocs.ParseFile("./demos/blast-rivals-2025-season-2-spirit-vs-vitality-bo3-KtWhzrlsNkWCCS0U9BIlr3/spirit-vs-vitality-m1-dust2.dem", func(p demoinfocs.Parser) error {
		criarMatrizes(p)
		return nil
	})
	if err != nil {
		log.Panic("failed to parse demo: ", err)
	}
}
