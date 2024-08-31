# Clone do Git

Esse projeto é um desafio da **[https://devgym.com.br](https://app.devgym.com.br/challenges/5b56d4a1-378c-41f0-9c91-7a9577d00671)**.

Implementação de uma versão básica de um sistema de versão como o git.

## Instalação

```bash
go install github.com/LuizGuilherme13/git-clone
```

## Como usar

### 1. Inicializando o diretório

```bash
~./go/bin/git-clone init
```

Irá criar uma pasta oculta '.backup' dentro do diretório atual, onde armazenará as copias dos arquivos:

### 2. Adicionando os arquivos

```bash
~./go/bin/git-clone add <nome_do_arquivo>...
```

Cria uma cópia do arquivo passado dentro da pasta '.backup':

### 2. Verificando o estado dos arquivos

```bash
~./go/bin/git-clone status
```

Exibe o estado atual de cada arquivo, sendo:

**Untracked** (O Arquivo nunca foi adicionado);

**Changes not staged** (O Arquivo mudou desde a última vez que ele foi adicionado);

**Changes to be commited** (O Arquivo está em staged e pronto para ser commitado);

### 3. Commitando

```bash
~./go/bin/git-clone commit "your_message_here"
```

Armazenará o estado atual do projeto em uma estrutura dentro de '.backup/objects'

### 4. Log

```bash
~./go/bin/git-clone log
```

Exibe o histórico de commits

## Work in progress...

- rm: Remover um arquivo do Index e da pasta própria pasta.

- restore: Desfaz um arquivo adicionado que ainda não foi commitado.

- checkout: Você pode ir pra qualquer commit do projeto e ver o estado do projeto naquele momento.

- Implementar branches.
