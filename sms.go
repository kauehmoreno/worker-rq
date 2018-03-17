package main

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// Trial credentials
var ACCOUNTSID = os.Getenv("ACCOUNTSID")
var AUTHTOKEN = os.Getenv("AUTHTOKEN")
var URL = "https://api.twilio.com/2010-04-01/Accounts/" + ACCOUNTSID + "/Messages.json"
var TWILIONUMBER = os.Getenv("TWILIONUMBER")

type SMS struct {
	Number       string `json:"number"`
	Email        string `json:"email"`
	ConfirmToken string `json:"confirmToken,omitempty"`
}

func ErrorSms(err error, number string, email string, msg string) {
	log.WithFields(log.Fields{
		"error":  err.Error(),
		"number": number,
		"email":  email,
		"time":   time.Now(),
	}).Error(msg)
}

// TODO Implement redis queue to insert it into
// and take a time sleep of 1 seconds to send sms based
// on politicy of Twilio EUA
func (sms *SMS) SendSMS() {
	v := url.Values{}
	v.Set("To", sms.Number)
	v.Set("From", TWILIONUMBER)
	sms.GenerateToken()
	msg := fmt.Sprintf("Access Token: %s", sms.ConfirmToken)
	v.Set("Body", msg)

	rb := *strings.NewReader(v.Encode())

	client := &http.Client{}

	req, erro := http.NewRequest("POST", URL, &rb)
	if erro != nil {
		ErrorSms(erro, sms.Number, sms.Email, "Error on Configure Request based on endpoint and Config")
	} else {
		req.SetBasicAuth(ACCOUNTSID, AUTHTOKEN)
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		// PAUSE CAUSED BY rules of USA number on twilio API
		time.Sleep(time.Second)

		resp, err := client.Do(req)
		if err != nil {

			ErrorSms(err, sms.Number, sms.Email, "Error on POST message to TWILIO")
		} else {
			if resp.StatusCode == 201 || resp.StatusCode == 200 {
				log.WithFields(log.Fields{
					"status": resp.StatusCode,
				}).Warning("Enviado com sucesso")
				sms.Set(15 * time.Minute)
			} else {
				log.WithFields(log.Fields{
					"status": resp.StatusCode,
					"resp":   resp.Request.Body,
					"header": resp.Request.Header,
				}).Error("Houve alguem erro na API")
			}
		}
	}

}

// GenerateToken and assign into ConfirmToken
func (sms *SMS) GenerateToken() {
	token := make([]byte, 3)
	rand.Read(token)
	sms.ConfirmToken = fmt.Sprintf("%x", token)
}
