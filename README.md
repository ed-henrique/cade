# Cadê

Rastreador de pedidos que usa a API dos Correios.

## Configuração do Ambiente de Desenvolvimento

1. Clone este repositório

```sh
git clone git@github.com:ed-henrique/cade.git
cd cade
```

2. Instale as dependências

```sh
npm install
```

3. Copie o arquivo `.env` a partir do `.env.example`

```sh
cp .env.example .env
```

4. Preencha o arquivo `.env`

```env
CHAVE_API_CORREIOS= # Insira sua chave da API dos Correios aqui
```

5. Execute o projeto

```sh
npm run dev
```

6. Acesse [aqui][http://localhost:5173]
