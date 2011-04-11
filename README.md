## Description

Simple SMS sending library in Go. Currently uses the http API from bulksms.com.

## Installation

Not a package yet ... TODO

## Usage

If you want to send a SMS only one method is of interest for you:

	Send(receivers []string, message string) os.Error

It takes a list of telephone numbers who shall receive the message and returns an os.Error is something went wrong.

## Example

	sms := NewBulkSMSSender("username", "password)
	sms.Testmode = 1			//don't send the sms, just perform an API supported test
	sms.RoutingGroup = 1	//let's use the cheap eco route
	
	msg := "hi, this is a test"
	receivers := []string{"49xxxxxxxx"}
	
	//let's see how much this sms would cost us
	_, quote := sms.GetQuote(receivers, msg)
	price := quote * 3.75 * 0.01	//quote is in credits. 1 credit = 3.75 eur cent

	fmt.Printf("the sms will cost us %.4f EUR\n", price)
	
	//send the sms
	if err := sms.Send(receivers, msg); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("sms sent!")
