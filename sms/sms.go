package sms 

import (
	"os"
	"fmt"
	"http"
	"io/ioutil"
	"strings"
	"strconv"
)

type SMSSender interface {
	Send(receivers []string, message string) os.Error
}

//implements the bulksms.com http API
type BulkSMSSMSSender struct {
	Username string
	Password string

	Testmode     int //0 - no testmode, 1 always succeed, -1 alswys fail
	RoutingGroup int //1 - economy, 2 - standard (default), 3 - premium. see bulksms.com for pricing
}

func NewBulkSMSSMSSender(username, password string) *BulkSMSSMSSender {
	return &BulkSMSSMSSender{
		Username:     username,
		Password:     password,
		Testmode:     0,
		RoutingGroup: 2,
	}
}


//get the total price (in CREDITS! NOT MONEYS!) for sending out a given message to a list of receivers.
//for credit prices see bulksms.com
func (sms *BulkSMSSMSSender) GetQuote(receivers []string, message string) (err os.Error, quote float64) {
	if sms.RoutingGroup < 1 || sms.RoutingGroup > 3 {
		err = os.NewError("Routing Group must be 1, 2 or 3!")
		return
	}
	rtgrp := strconv.Itoa(sms.RoutingGroup)

	url := "http://bulksms.vsms.net:5567/eapi/submission/quote_sms/2/2.0"
	rcvrs := strings.Join(receivers, ",")

	data := map[string]string{
		"username":      sms.Username,
		"password":      sms.Password,
		"message":       message,
		"msisdn":        rcvrs,
		"routing_group": rtgrp,
	}

	c := &http.Client{}
	resp, err := c.PostForm(url, data)
	if err != nil {
		return err, 0.0
	}

	respbytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err, 0.0
	}
	resp.Body.Close()

	respstr := strings.Trim(string(respbytes[0:]), "\r\n")
	respitems := strings.Split(respstr, "|", -1)

	respcode, err := strconv.Atoi64(respitems[0])
	if err != nil {
		return err, 0.0
	}

	if respcode != 0 {
		return os.NewError(fmt.Sprintf("err code %d: %s", respcode, respitems[1])), 0.0
	}

	quote, err = strconv.Atof64(respitems[2])
	return
}


//sends a message to a list of receivers.
//a receiver is a string containing an international telephone number
//without a leading + or 0. to send a sms to a german (+49) number
//you'd use "49172xxxxxx".
//the message strings max len is 160 chars
func (sms *BulkSMSSMSSender) Send(receivers []string, message string) os.Error {
	if sms.RoutingGroup < 1 || sms.RoutingGroup > 3 {
		return os.NewError("Routing Group must be 1, 2 or 3!")
	}
	rtgrp := strconv.Itoa(sms.RoutingGroup)

	url := "http://bulksms.vsms.net:5567/eapi/submission/send_sms/2/2.0"
	rcvrs := strings.Join(receivers, ",")

	data := map[string]string{
		"username":      sms.Username,
		"password":      sms.Password,
		"message":       message,
		"msisdn":        rcvrs,
		"routing_group": rtgrp,
	}

	if sms.Testmode == -1 {
		data["test_always_fail"] = "1"
	}
	if sms.Testmode == 1 {
		data["test_always_succeed"] = "1"
	}

	c := &http.Client{}
	resp, err := c.PostForm(url, data)
	if err != nil {
		return err
	}

	respbytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	resp.Body.Close()

	respstr := string(respbytes[0:])
	respitems := strings.Split(respstr, "|", -1)

	respcode, err := strconv.Atoi64(respitems[0])
	if err != nil {
		return err
	}

	if respcode != 0 {
		return os.NewError(fmt.Sprintf("err code %d: %s", respcode, respitems[1]))
	}

	return nil
}
