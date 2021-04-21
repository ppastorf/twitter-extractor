# Projeto SSC0158 - Grupo 13
## Especificação
https://edisciplinas.usp.br/mod/assign/view.php?id=3544468

## Integrantes
- Douglas Tallach, 9292708
  douglas.tallach@usp.br
- Juliano Fantozzi, 9791218
  juliano.fantozzi@usp.br
- Pedro Pastorello Fernandes, 10262502
  pedropastorf@usp.br
- Vinicius Leite Ribeiro, 10388200
  viniribeiro@usp.br

## Contribuindo
- Clonar o repo e criar uma branch nova
    - Cada integrante deve usar branches próprias para codar
    - Vamos integrar tudo com Merge Requests, será a unica forma de commitar na Master

- É uma boa configurar seu usuário e email do git com os definidos na planilha dos grupos
    - git config --local user.email <email>
    - git config --local user.name <name>

## Acesso a nossa VM
- username = gcloud13
- portas = 21111, 21112, 21113

ssh $username@andromeda.lasdpc.icmc.usp.br -p $porta

#### Execução
Para executar o que foi desenvolvido até agora, vocês vão precisar de algumas credenciais, que não foram commitadas no repositório por motivos obvios. Vocês podem entrar em contato conosco que podemos passar elas de alguma outra maneira. Vocês também pode usar credenciais próprias para os serviços se quiserem.

*obs: todos os caminhos de arquivo são relativos a raiz do repositório*

**1. worker-extractor**
Para executar o código do worker de extração:

```bash
cd app/worker-extractor
go run .
```

Você vai precisar de um arquivo de credenciais no seguinte formato em app/`worker-producer/secret/credentials.yaml`:

```yaml
twitter:
  CONSUMER_KEY: "<credencial_twitter>"
  CONSUMER_KEY_SECRET: "<credencial_twitter>"
  ACCESS_TOKEN: "<credencial_twitter>"
  ACCESS_TOKEN_SECRET: "<credencial_twitter>"
```

**2. Deploy de infra**
Todos os scripts relacionados a infra-estrutura estão na pasta `deploy/`. Vocês vão precisar que os seguintes arquivos de credencial existam na pasta `secret/`. Os arquivos exportam os valores como variáveis de ambiente para facilitar o acesso aos valores em scripts.


- `secret/aws_access`
Esse arquivo será usado para se autenticar na AWS para criar o ambiente de desenvolvimento. Uma chave privada para a conexão ssh com as máquinas da AWS deve exister nesta mesma pasta.

```bash 
export AWS_ACCESS_KEY_ID="<credencial_aws>"
export AWS_SECRET_ACCESS_KEY="<credencial_aws>"
export AWS_PRIVATE_KEY_FILE="<caminho_para_chave_privada_aws>"
```

- `secret/db_password`
Esse arquivo seta a senha do banco. Pode ser qualquer valor armazenável em uma variável de ambiente.

```bash
export DB_PASSWORD="<senha_do_banco>"
```

- `secret/gitlab_deploy_token`
Esse arquivo será usado para publicar imagens de container no nosso registry do Gitlab.

```bash 
export GITLAB_DEPLOY_TOKEN="<token_gitlab>"
export GITLAB_DEPLOY_USER="<token_gitlab>"
```

- `secret/lab_access`
Senha de ssh da máquina do lab do nosso grupo.

```bash 
export LAB_SSH_PASS="<senha_vm_lab>"
```