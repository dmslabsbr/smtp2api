FROM golang
LABEL author="Daniel S"

# Defina o diretório de trabalho
WORKDIR /go/src/smtp2api

# Define e reseta  env
#ENV http_proxy=""
#ENV https_proxy="" 
#ENV HTTP_PROXY="" 
#ENV HTTPS_PROXY=""

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./

# Baixe todas as dependências
RUN go mod download && go mod verify


# Copie o resto dos arquivos do diretório
COPY . .

# Instale o pacote
RUN go build -v -o .

# Define o ponto de entrada
#ENTRYPOINT [ "/go/src/smtp2api" ]

CMD ["./smtp2api"]

EXPOSE 1025