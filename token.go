package token

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var Ccu ClientConectUser
var retToken RetornoToken
var Token string

type ClientConectUser struct {
	ClientConect ClientHttpUser
}

type ClientHttpUser struct {
	Client *http.Client
	Host   string
	User   string
	Auth   string
}

// RetornoToken estrutura do Token do Sistema
type RetornoToken struct {
	Ambiente string `json:"ambiente"`
	ID       string `json:"id"`
	Perfil   string `json:"perfil"`
	Emissao  string `json:"emissao"`
	ExpiraEm string `json:"expiraEm"`
	Token    string `json:"token"`
}

func GetToken(host, user, pass string) (RetornoToken, error) {
	Ccu = NewClientConectUser(host, user, pass)
	token, err := Ccu.GerarToken()
	if err != nil {
		return token, err
	}
	return token, nil
}

func Start(host, user, pass string) {
	for {
		tk, err := GetToken(host, user, pass)
		Token = tk.Token
		if err != nil {
			log.Panic(err.Error())
		}
		time.Sleep(30 * time.Minute)
	}
}

func NewClientConectUser(host, user, pass string) ClientConectUser {
	return ClientConectUser{
		ClientConect: NewClientUser(host, user, pass),
	}
}

func NewClientUser(enderecoBase, usuario, senha string) ClientHttpUser {

	client := ClientHttpUser{
		Client: &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}},
		Host:   enderecoBase,
		User:   usuario,
		Auth: fmt.Sprintf(
			"Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", usuario, senha))),
		),
	}
	return client
}

func (c ClientHttpUser) NewRequest(method, endpoint string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", c.Host, endpoint), body)
	if err != nil {
		return nil, fmt.Errorf("erro: cmp: newrequest 1: %s", err)
	}
	req.Header.Add("Authorization", c.Auth)
	return req, nil
}

func (ccu ClientConectUser) GerarToken() (RetornoToken, error) {
	req, err := ccu.ClientConect.NewRequest("POST", "", nil)
	if err != nil {
		return retToken, err
	}
	res, err := ccu.ClientConect.Client.Do(req)
	if err != nil {
		//  return Cartao{}, fmt.Errorf("GerarToken 2: %s", err)
		return retToken, err
	}
	if res.StatusCode != 200 && res.StatusCode != 201 {
		return retToken, fmt.Errorf("gerartoken: %s", res.Status)
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return retToken, err
	}
	defer res.Body.Close()

	err = json.Unmarshal(b, &retToken)
	if err != nil {
		return retToken, err
	}

	return retToken, nil
}
