package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/emersion/go-smtp"
	brevo "github.com/getbrevo/brevo-go/lib"

	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	//"html/template"
)

// The Backend implements SMTP server methods.
type Backend struct{}

type EmailData struct {
	From       string
	SenderName string
	To         []string
	ToName     []string
	Content    string
	Subject    string
	port       int
	host       string
	apiKey     string
	byPassAuth bool
}

type segredosType struct {
	brevoApiKey string
}

var globalEmailData EmailData
var globalSegredos segredosType


// pega segredos das variaveis de ambiente
func getSegredos() {
	password, exists := os.LookupEnv("BREVO_APIKEY")
	if !exists || password == "" {
		fmt.Println("Variável de ambiente BREVO_APIKEY não está definida ou está vazia!")
		log.Fatal("Variável não definida! Não posso continuar.")
	} else {
		fmt.Println("Variável de ambiente BREVO_APIKEY está definida.")
		// Use a senha aqui
		globalSegredos.brevoApiKey=password
	}
}

func (bkd *Backend) NewSession(_ *smtp.Conn) (smtp.Session, error) {
	log.Println("NewSession")
	return &Session{}, nil
}

// A Session is returned after EHLO.
type Session struct {
	auth bool
}

func (s *Session) AuthPlain(username, password string) error {
	log.Println("AuthPlain")
	if username != "username" || password != "password" {
		return smtp.ErrAuthFailed
	}
	s.auth = true
	return nil
}

func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	log.Println("Mail")
	if globalEmailData.byPassAuth {
		s.auth = true
	}
	if !s.auth {
		log.Println("Mail Erro Autenticação")
		return smtp.ErrAuthRequired
	}
	log.Println("Mail from:", from)
	globalEmailData.From = from
	return nil
}

func (s *Session) Rcpt(to string) error {
	if !s.auth {
		log.Println("Rcpt Erro Autenticação")
		return smtp.ErrAuthRequired
	}
	log.Println("Rcpt to:", to)
	globalEmailData.To = append(globalEmailData.To, to)
	globalEmailData.ToName = append(globalEmailData.ToName, "")
	return nil
}

func (s *Session) Data(r io.Reader) error {
	if !s.auth {
		log.Println("Data Erro Autenticação")
		return smtp.ErrAuthRequired
	}
	if b, err := ioutil.ReadAll(r); err != nil {
		return err
	} else {
		log.Println("Data:", string(b))
		globalEmailData.Content = string(b)
		enviaEmail()
	}
	return nil
}

func (s *Session) Reset() {
	log.Println("Reset")
}

func (s *Session) Logout() error {
	log.Println("Logout")
	return nil
}

func formatToJSON1(email, name string) string {
	if name == "" { name = "sem nome" }
	data := struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}{
		Email: email,
		Name:  name,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	fmt.Println("string(jsonData):", string(jsonData))
	return string(jsonData)
}

func formatToJSON(toEmail, toName []string) string {
	formattedToJSON := ""
	fmt.Println("ToName:", toName)
	fmt.Println("len(toName):", len(toName))
	fmt.Println("toEmail:", toEmail)
	fmt.Println("len(toEmail):", len(toEmail))
	for i, toE := range toEmail {
		fmt.Println("i:", i)
		fmt.Println("toE:", toE)
		fmt.Println("toName[i]:", toName[i])
		formattedToJSON += formatToJSON1(toE, toName[i])
		if i < len(toEmail)-1 {
			formattedToJSON += ","
		}
	}
	fmt.Println("formattedToJSON:", formattedToJSON)
	return formattedToJSON
}

func enviaEmail() {

	fmt.Println("Preparando envio de e-mail...")

	var ctx context.Context
	cfg := brevo.NewConfiguration()
	//Configure API key authorization: api-key
	cfg.AddDefaultHeader("api-key", globalEmailData.apiKey)
	//Configure API key authorization: partner-key
	//cfg.AddDefaultHeader("partner-key", globalEmailData.apiKey)
	//
	//cfg.AddDefaultHeader("content-type", "application/json")
	//cfg.AddDefaultHeader("accept", "application/json")

	br := brevo.NewAPIClient(cfg)
	result, resp, err := br.AccountApi.GetAccount(ctx)
	if err != nil {
		fmt.Println("Error when calling AccountApi->get_account: ", err.Error())
		return
	}
	fmt.Println("GetAccount Object:", result, " GetAccount Response: ", resp)

	fmt.Println("")

	payloadJson := `
	{
		"sender": {
		  "name": "!senderName!", 
		  "email": "!senderEmail!"
		},
		"to": [
			!jsonTo!
		],
		"subject": "!subject!",
		"htmlContent": !HTMLContent!
	  }
	  `
	payloadJson = strings.Replace(payloadJson, "!senderName!", globalEmailData.SenderName, -1)
	payloadJson = strings.Replace(payloadJson, "!senderEmail!", globalEmailData.From, -1)
	payloadJson = strings.Replace(payloadJson, "!subject!", globalEmailData.Subject, -1)

	//globalEmailData.Content = string( template.HTML(globalEmailData.Content) )

	// Serializa a estrutura em JSON
	jsonData, err := json.Marshal(globalEmailData.Content)
	if err != nil {
		log.Fatalf("Erro ao serializar JSON: %v", err)
	}

	globalEmailData.Content = string(jsonData)

	payloadJson = strings.Replace(payloadJson, "!HTMLContent!", globalEmailData.Content, -1)

	jsonTo := `
	    {
		"email": "xyz@yahoo.com",
		"name": "nome"
	  }`

	jsonTo = formatToJSON(globalEmailData.To, globalEmailData.ToName)

	fmt.Println("jsonTo:", jsonTo)

	payloadJson = strings.Replace(payloadJson, "!jsonTo!", jsonTo, -1)

	payload := []byte(payloadJson)

	fmt.Println("Iniciando o request para ", globalEmailData.host)
	req, err := http.NewRequest("POST", globalEmailData.host, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Erro ao criar a requisição:", err)
		return
	}

	fmt.Println("irei tentar enviar o e-mail")

	req.Header.Set("api-key", globalEmailData.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("accept", "application/json")

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		fmt.Println("Erro ao fazer a requisição:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Status Code:", resp.Status)
	if resp.StatusCode == 400 {
		fmt.Println("Bad Request!!")
		fmt.Println("payloadJson:", payloadJson)
	}
	fmt.Println("Resp:", resp)

	return
}

func main() {

	// load secrets
	getSegredos()

	be := &Backend{}

	s := smtp.NewServer(be)

	s.Addr = ":1025"
	s.Domain = "localhost"
	s.ReadTimeout = 100 * time.Second
	s.WriteTimeout = 100 * time.Second
	s.MaxMessageBytes = 1024 * 1024
	s.MaxRecipients = 50
	s.AllowInsecureAuth = false

	globalEmailData = EmailData{
		host:       "https://api.brevo.com/v3/smtp/email",
		port:       443,
		apiKey:     globalSegredos.brevoApiKey,
		From:       "xyz@yahoo.com",
		SenderName: "Máquina",
		To:         []string{},
		ToName:     []string{},
		Content:    "Teste de envio de e-mail",
		Subject:    "Assunto do e-mail",
		byPassAuth: true,
	}

	globalEmailData.Content = "Teste de envio de e-mail"

	//enviaEmail()

	log.Println("Starting server at", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}


//todo verificar domínio do e-mail para evitar span
//todo colocar como parametro a porte do servidor, endereço da api.