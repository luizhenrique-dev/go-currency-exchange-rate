## Desenvolvimento de Sistemas em Go - Cotação do Dólar

**1. Estrutura do Projeto:**
   - client.go
   - server.go

**2. Funcionalidade do client.go:**
   - Realiza uma requisição HTTP no server.go solicitando a cotação do dólar.
   - O client.go precisará receber do server.go apenas o valor atual do câmbio (campo "bid" do JSON)
   - Salva a cotação atual em um arquivo "cotacao.txt" no formato: Dólar: {valor}
   - Usa o package "context" com um timeout máximo de 300ms para receber o resultado do server.go.

**3. Funcionalidade do server.go:**
   - Consuma a API de câmbio de Dólar para Real (https://economia.awesomeapi.com.br/json/last/USD-BRL).
   - Retorna, no formato JSON, o resultado da cotação para o cliente.
   - Usa o package "context" para registrar no banco de dados SQLite cada cotação recebida.
   - Define um timeout máximo de 200ms para chamar a API de cotação do dólar e 10ms para persistir os dados no banco.

**4. Endpoints:**
   - `/cotacao`: Endpoint que o client.go acessará para obter a cotação.
   - `/list`: Endpoint opcional para buscar no banco de dados todas as cotações e retornar seus dados em um array de JSON.

**5. Logs e Timeout:**
   - Os contextos deverão retornar erro nos logs caso o tempo de execução seja insuficiente.

**6. Para executar o server:**
   - A porta a ser utilizada pelo servidor HTTP será a 8080, portanto ela precisa estar disponível.
   - No terminal, navegue até o diretório onde o projeto do aplicativo está localizado. Em seguida, execute o seguinte comando para executá-lo:
  
```sh
cd server && go run server.go database.go
```

**7. Para executar o client:**
   - No terminal, navegue até o diretório onde o projeto do aplicativo está localizado. Em seguida, execute o seguinte comando para executá-lo:

```sh
cd client && go run client.go
```