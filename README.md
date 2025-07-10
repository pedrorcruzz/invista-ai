

<p align="center">
  <img src="public/logo.png" alt="Logo InvistAI" width="180"/>
</p>

<h1 align="center">Versão CLI</h1>



<p align="center">
  <b>Gerencie seus investimentos. Controle seus gastos. Tudo no seu terminal.</b>
</p>

---
> **🔔 Nota:** Este repositório reúne o antigo `smart-spending-checker CLI` (também criado por mim) em uma única ferramenta para **gestão de investimentos** e **controle inteligente de gastos**.  
> Todos os recursos agora estão centralizados em um só lugar!


## ✨ Funcionalidades

- 📈 <b>Controle de Investimentos</b> — Adicione/edite dados mensais, veja lucros brutos/líquidos e visualize seu progresso.
- 🧠 <b>Gestor Inteligente de Gastos</b> — Planeje compras, gerencie parcelas e receba recomendações inteligentes.
- 💾 <b>Dados Locais</b> — Todos os seus dados são salvos localmente em arquivos JSON simples.
- 🖥️ <b>Interface Bonita no Terminal</b> — Menus modernos com bordas para uma experiência CLI agradável.
- 🐚 <b>CLI Universal</b> — Use com <code>go run</code>, construa um binário ou chame de scripts <code>fish</code>, <code>zsh</code>, <code>sh</code> em qualquer lugar.

---

## 🚀 Primeiros Passos

### 1. Clone o Repositório

```sh
git clone https://github.com/pedrorcruz/invista-ai-cli
cd invista-ai-cli
```

### 2. Rodar com Go

```sh
go run main.go
```

### 3. Buildar & Usar em Qualquer Lugar

```sh
go build -o invista-ai
./invista-ai
```

---

## Automatizando o Acesso de Qualquer Lugar no Terminal

Para rodar o InvistAI de qualquer diretório no seu terminal, você pode criar um script e uma função (ou alias).

### 1. Crie um Script Shell

Crie um arquivo chamado `invista-ai.sh` (ou qualquer nome que preferir) em um diretório de sua escolha (ex: `~/.dotfiles/scripts`). Adicione o conteúdo abaixo, **trocando o caminho do `cd` para o local correto do seu projeto**:

```bash
#!/bin/bash

cd ~/Developer/Scripts/invista-ai  # ⚠️ TROQUE PELO SEU CAMINHO REAL

./invista-ai  # ⚠️ TROQUE PELO NOME DO SEU BINÁRIO
sleep 1.3
clear
```

### 2. Torne o Script Executável

Dê permissão de execução ao script:

```bash
chmod +x invista-ai.sh
```

### 3. Crie uma Função (Fish) ou Alias (Zsh/Bash)

#### Fish Shell

Adicione a função abaixo ao seu arquivo ~/.config/fish/config.fish:

```fish
function invista-ai
    set prev_dir (pwd)
    cd ~/.dotfiles/scripts # ⚠️ TROQUE PELO DIRETÓRIO DO SEU SCRIPT
    ./invista-ai.sh
    cd $prev_dir
end
```

#### Zsh/Bash

Adicione o alias abaixo ao seu ~/.zshrc ou ~/.bashrc:

```bash
alias invista-ai="cd ~/.dotfiles/scripts && ./invista-ai.sh && cd -" # ⚠️ TROQUE PELO DIRETÓRIO DO SEU SCRIPT
```

### 4. Recarregue sua Configuração do Shell

Após adicionar a função ou alias, recarregue sua configuração:

#### Fish

```bash
source ~/.config/fish/config.fish
```

#### Zsh

```bash
source ~/.zshrc
```

#### Bash

```bash
source ~/.bashrc
```

Agora você pode rodar o InvistAI de qualquer diretório apenas digitando `invista-ai` no terminal.

---

## 🧩 Menus

### Menu Principal

```
1. Ver resumo completo (vertical)
2. Ver resumo completo (horizontal)
3. Adicionar/editar mês
4. Gestor Inteligente de Gastos
5. Sair
```

### Gestor Inteligente de Gastos

```
1. Adicionar produto
2. Remover produto
3. Listar meses
4. Atualizar lucro mensal
5. Editar produto
6. Antecipar parcelas
7. Configurar porcentagem segura
8. Voltar ao menu principal
```

- Você pode selecionar produtos pelo número ou digitando o nome!
- Todos os menus são exibidos em caixinhas para clareza e estilo.

---

## 📦 Onde os Dados São Salvos

- Dados de investimentos: <code>dados.json</code>
- Dados do gestor de gastos: <code>data/products.json</code>

---

## 📝 Licença & Créditos

- LICENÇA [MIT](https://github.com/pedrorcruzz/invista-ai/blob/develop/LICENSE)
- Criado por [Pedro Rosa](https://github.com/pedrorcruzz)

---

<p align="center">
  <b>Gerencie sua vida financeira direto do terminal!</b>
</p>
