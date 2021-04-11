# GO TO-DO SERVICE

## Build
```sh
# Linux
$ make
# output: bin/main-linux-amd64

# Ambiente atual
$ make cross
# output: bin/main
```

## Deploy da aplicação
```txt
Estrutura dos arquivos

├-
├-settings.json
├-bin/
    ├-main-linux-amd64
```

Executando a aplicação
```sh
# Modo de produção (ENV VAR 'APP_MODE' = 'production')
$ APP_MODE=production ./bin/main-linux-amd64

# Modo default desenvolvimento
$ ./bin/main-linux-amd64
```
