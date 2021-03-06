package doauth

import (
	"bytes"
	"context"
	"github.com/digitalocean/godo"
	"github.com/harshpreet93/dopaas/errorcheck"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
)

type TokenSource struct {
	AccessToken string
}

func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

func getFieldFromFile(filepath string, fieldName string) string {
	viper.SetConfigType("yaml")

	dat, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Println("error reading yaml config ", err)
	}
	viper.ReadConfig(bytes.NewBuffer(dat))

	return viper.Get(fieldName).(string) // this would be "steve"
}

func getTokenFromFS() string {
	// get token from ~/.dopaas.yaml
	// file format is:
	// DIGITALOCEAN_ACCESS_TOKEN: "blahhhh"

	dopaasConf, err := homedir.Expand("~/.dopaas.yaml")
	errorcheck.ExitOn(err, "error getting dopaas file")
	return getFieldFromFile(dopaasConf, "DIGITALOCEAN_ACCESS_TOKEN")
}

func Auth() *godo.Client {
	pat := getTokenFromFS()
	tokenSource := &TokenSource{
		AccessToken: pat,
	}

	oauthClient := oauth2.NewClient(context.Background(), tokenSource)
	client := godo.NewClient(oauthClient)
	return client
}
