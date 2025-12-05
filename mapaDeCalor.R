library(readr)

tabela_jogador <- read_csv("dadosJogador.csv",
  col_types = cols(
    idJogador = col_character(),
    nomeJogador = col_character()
  )
)
tabela_estado_jogo <- read_csv("estadoJogo.csv",
  col_types = cols(
    tick = col_integer(),
    tipoEstado = col_integer(),
    nomeEstado = col_character()
  )
)
tabela_posicoes <- read_csv(
  "posicoesJogadores.csv",
  col_types = cols(
    idJogador = col_character(),
    tick = col_integer(),
    posiçãoX = col_double(),
    posiçãoY = col_double(),
    posiçãoZ = col_double()
  )
)

#salva o ultimo tick de início da rodada
#pra caso a partida tenha sido reiniciada
inicio <- max(tabela_estado_jogo$tick[tabela_estado_jogo$tipoEstado == 5])
#último (ou maior)

#filtra a tabela de estado
#removendo os ticks anteriores ao reinicio da partida
tabela_estado_jogo <- tabela_estado_jogo[tabela_estado_jogo$tick >= inicio, ]

#mudar de indice para tick
#salva os ticks de cada inicio de rodada (fim do timefreeze)
ticks_inicio_round <- tabela_estado_jogo$tick[tabela_estado_jogo$tipoEstado == 2] # nolint: line_length_linter.

#salva os ticks de cada fim de rodada
ticks_fim_round    <- tabela_estado_jogo$tick[tabela_estado_jogo$tipoEstado == 3] # nolint: line_length_linter.

#sabendo que apos o inicio da partida todo round que inicia uma hora termina
#o tamanho dos vetores de inicios e fins de rounds deve ser igual
if (length(ticks_inicio_round) != length(ticks_fim_round)) {
  stop("Erro: número de inícios e fins de rounds não correspondem.")
}else {
  # remover posições após o fim do último round
  # e as posições anteriores ao reinício da partida
  if (length(ticks_fim_round) >= 1) {
    tabela_posicoes <- tabela_posicoes[tabela_posicoes$tick <= ticks_fim_round[length(ticks_fim_round)], ] # nolint: line_length_linter.
    tabela_posicoes <- tabela_posicoes[tabela_posicoes$tick >= inicio, ] # nolint: line_length_linter.
  }

  tabela_posicoes_limpa <- data.frame(
    idJogador = character(),
    tick = integer(),
    posiçãoX = double(),
    posiçãoY = double(),
    posiçãoZ = double()
  )

  #removendo intervalos entre fim de um round e início do próximo
  if (length(ticks_fim_round) > 1) {
    for (k in seq_len(length(ticks_fim_round) - 1)) {
      gap_start <- ticks_inicio_round[k]
      gap_end   <- ticks_fim_round[k]
      # adicionando os dados validos à nova tabela
      tabela_posicoes_limpa <- rbind(tabela_posicoes_limpa, tabela_posicoes[tabela_posicoes$tick > gap_start & tabela_posicoes$tick < gap_end, ]) # nolint: line_length_linter.
    }
  }
}

#arrumando a escala das posições
#o resultado são coordenadas em pixels relativas à imagem do radar
#formula:
#(x - map.PosX) / map.Scale , (map.PosY - y) map.Scale
tabela_posicoes_limpa$posiçãoX <- (tabela_posicoes_limpa$posiçãoX - (-2476)) / 4.4 # nolint: line_length_linter.
tabela_posicoes_limpa$posiçãoY <- -((3239 - tabela_posicoes_limpa$posiçãoY) / 4.4) # nolint: line_length_linter.

for (nome in tabela_jogador$nomeJogador) {
  if (nome == "donk") {
    donk_id <- tabela_jogador$idJogador[tabela_jogador$nomeJogador == nome]
    donk_posicoes <- subset(tabela_posicoes_limpa, idJogador == donk_id)
  }else if (nome == "ZywOo") {
    zywoo_id <- tabela_jogador$idJogador[tabela_jogador$nomeJogador == nome]
    zywoo_posicoes <- subset(tabela_posicoes_limpa, idJogador == zywoo_id)
  }
}

library(ggplot2)
library(ggpubr)
library(png) # Or jpeg, depending on your image type

img <- readPNG("de_dust2_radar_psd.png") # Load the image file
p <- ggplot(donk_posicoes, aes(x = posiçãoX, y = posiçãoY)) +
  background_image(img) +
  stat_density_2d(aes(fill = ..density..), geom = "raster", contour = FALSE, alpha = 0.5) +
  coord_fixed(ratio = 1) +
  scale_fill_distiller(palette = "Spectral", direction = -1) +
  scale_x_continuous(expand = c(0, 0), limits = c(0, 1024)) +
  scale_y_continuous(expand = c(0, 0), limits = c(0, -1024)) +
  theme(
    legend.position = "none"
  )
print(p)