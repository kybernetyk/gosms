# gosms

## Description

Simple SMS sending library for Go 1. Currently uses the http API from bulksms.com.

## Installation

Either clone this repository or use goinstall:
	$ [sudo -E] goinstall github.com/jsz/gosms/sms

## Usage

If you want to send a SMS only one method is of interest for you:

	Send(receivers []string, message string) error

It takes a list of telephone numbers who shall receive the message and returns an error if something went wrong.

## Example
	package main

	import (
		"github.com/jsz/gosms"
		import "fmt"
	)

	func main() {
		s := gosms.NewBulkSMSSMSSender("username", "password")
		s.Testmode = 1     //don't send the sms, just perform an API supported test
		s.RoutingGroup = 1 //let's use the cheap eco route

		msg := "hi, this is a test"
		receivers := []string{"49xxxxxx"}  //put a proper tel# here in

		//let's see how much this sms would cost us
		_, quote := s.GetQuote(receivers, msg)
		price := quote * 3.75 * 0.01 //quote is in credits. 1 credit = 3.75 eur cent

		fmt.Printf("the sms will cost us %.4f EUR\n", price)

		//send the sms
		if err := s.Send(receivers, msg); err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("sms sent!")
	}

#License 

(c) Leon Szpilewski 2011, 2012

Licensed under MIT license

#Contact

email: leon.szpilewski@gmail.com
twitter: @fettemama
blog: http://fettemama.org

