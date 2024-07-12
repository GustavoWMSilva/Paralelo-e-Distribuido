# README.md

## Trabalho Prático de Distribuídos

### Autores
- Gustavo Willian Martins da Silva
- Lorenzo Duarte More

### Introdução

Em sistemas distribuídos, é comum que um processo especial funcione como coordenador, desempenhando um papel central para os outros processos. No entanto, esse coordenador pode falhar por diversos motivos, como problemas de rede ou hardware. Quando isso ocorre, é necessário eleger um novo coordenador para garantir que o sistema continue funcionando corretamente.

Este trabalho implementa um algoritmo de eleição distribuída baseado em anel lógico. Quando um coordenador falha, um processo ativo deve iniciar uma eleição para escolher um novo líder com base em um critério de prioridade.

### Descrição do Algoritmo

O algoritmo de eleição baseado em anel lógico funciona da seguinte maneira:

1. **Formação do Anel**:
   - Os processos se conectam formando uma sequência lógica.
   - Cada processo conhece o anel, mas envia mensagens apenas para o próximo processo ativo na sequência.

2. **Detecção de Falha**:
   - Quando o coordenador falha, o processo que detecta a falha inicia uma eleição.
   - O processo coloca sua prioridade em uma mensagem e envia para o próximo processo ativo.

3. **Processo de Eleição**:
   - A mensagem de eleição circula pelo anel.
   - Cada processo ativo insere sua prioridade na mensagem.
   - Quando a mensagem retorna ao iniciador da eleição, este verifica a maior prioridade e define o novo líder.

4. **Comunicação do Novo Líder**:
   - O processo iniciador informa a todos os processos ativos sobre o novo líder.
   - Todos os processos passam a monitorar o novo coordenador.

### Implementação

A implementação foi realizada na linguagem Go, utilizando goroutines e canais para simular os processos e a comunicação entre eles. A estrutura básica inclui:

- **Tipos de Mensagem**:
  - `tipo 2`: Indica que um processo deve falhar.
  - `tipo 3`: Inicia uma eleição.
  - `tipo 4`: Coloca a prioridade (ID) na mensagem de eleição.
  - `tipo 5`: Informa o novo líder.
  - `tipo 6`: Mensagem de confirmação de liderança.
  - `tipo 7`: Finaliza o processo.

- **Controle de Eleição**:
  - Um processo de controle simula falhas e reativações de processos, bem como inicia eleições conforme necessário.

### Simulação e Testes

Para testar o algoritmo, foram realizados os seguintes passos:

1. **Falha do Coordenador Inicial**:
   - O processo 0 (coordenador inicial) é configurado para falhar.
   - O controle confirma a falha e solicita ao processo 1 que inicie uma nova eleição.
   - O processo 3 é eleito como novo líder.

2. **Falha de um Novo Coordenador**:
   - O processo 3 é configurado para falhar.
   - O controle solicita ao processo 2 que inicie uma nova eleição.
   - O processo 2 é eleito como novo líder.

3. **Reativação de Processos**:
   - O processo 0 é reativado e solicita uma nova eleição, mas não vence devido à prioridade menor.
   - O processo 3 é reativado e vence a eleição, tornando-se o novo líder.

4. **Finalização dos Processos**:
   - O controle envia mensagens de finalização para todos os processos, encerrando a simulação.

### Execução do Código

Para executar a simulação, certifique-se de ter o ambiente Go configurado. Em seguida, compile e execute o programa:

```sh
go run main.go
```

### Conclusão

Este trabalho demonstrou a implementação de um algoritmo distribuído de eleição baseado em anel lógico. O algoritmo foi capaz de eleger novos coordenadores de forma eficiente em um ambiente simulado, garantindo a continuidade do sistema mesmo em caso de falhas.

### Contato

Para mais informações, entre em contato com os desenvolvedores:

- Gustavo Willian Martins da Silva
- Lorenzo Duarte More

---

Este trabalho foi realizado como parte da disciplina de Sistemas Distribuídos, sob a orientação do professor DeRose.
