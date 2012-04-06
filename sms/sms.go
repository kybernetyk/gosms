/* 
A package for sending out SMS. Currently works only with bulksms.com
but may support more services if the need should arise.

Author: Leon Szpilewski / http://nntp.pl
*/
package sms

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type SMSSender interface {
	Send(receivers []string, message string) error
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
func (sms *BulkSMSSMSSender) GetQuote(receivers []string, message string) (err error, quote float64) {
	if sms.RoutingGroup < 1 || sms.RoutingGroup > 3 {
		err = errors.New("Routing Group must be 1, 2 or 3!")
		return
	}
	rtgrp := strconv.Itoa(sms.RoutingGroup)

	endpoint_url := "http://bulksms.vsms.net:5567/eapi/submission/quote_sms/2/2.0"
	rcvrs := strings.Join(receivers, ",")

	// data := url.Values{
	// 	"username":      sms.Username,
	// 	"password":      sms.Password,
	// 	"message":       message,
	// 	"msisdn":        rcvrs,
	// 	"routing_group": rtgrp,
	// }

	data := url.Values{}
	data.Set("username", sms.Username)
	data.Set("password", sms.Password)
	data.Set("message", message)
	data.Set("msisdn", rcvrs)
	data.Set("routing_group", rtgrp)

	c := &http.Client{}
	resp, err := c.PostForm(endpoint_url, data)
	if err != nil {
		return err, 0.0
	}

	respbytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err, 0.0
	}
	resp.Body.Close()

	respstr := strings.Trim(string(respbytes[0:]), "\r\n")
	respitems := strings.Split(respstr, "|")

	respcode, err := strconv.ParseInt(respitems[0], 10, 64)
	if err != nil {
		return err, 0.0
	}

	if respcode != 0 {
		return errors.New(fmt.Sprintf("err code %d: %s", respcode, respitems[1])), 0.0
	}

	quote, err = strconv.ParseFloat(respitems[2], 64)
	return
}

//sends a message to a list of receivers.
//a receiver is a string containing an international telephone number
//without a leading + or 0. to send a sms to a german (+49) number
//you'd use "49172xxxxxx".
//the message strings max len is 160 chars
func (sms *BulkSMSSMSSender) Send(receivers []string, message string) error {
	if sms.RoutingGroup < 1 || sms.RoutingGroup > 3 {
		return errors.New("Routing Group must be 1, 2 or 3!")
	}
	rtgrp := strconv.Itoa(sms.RoutingGroup)

	url := "http://bulksms.vsms.net:5567/eapi/submission/send_sms/2/2.0"
	rcvrs := strings.Join(receivers, ",")

	data := map[string][]string{
		"username":      []string{sms.Username},
		"password":      []string{sms.Password},
		"message":       []string{message},
		"msisdn":        []string{rcvrs},
		"routing_group": []string{rtgrp},
	}

	if sms.Testmode == -1 {
		data["test_always_fail"] = []string{"1"}
	}
	if sms.Testmode == 1 {
		data["test_always_succeed"] = []string{"1"}
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
	respitems := strings.Split(respstr, "|")

	respcode, err := strconv.ParseInt(respitems[0], 10, 64)
	if err != nil {
		return err
	}

	if respcode != 0 {
		return errors.New(fmt.Sprintf("err code %d: %s", respcode, respitems[1]))
	}

	return nil
}
