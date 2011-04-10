package main

import (
	"fmt"
)

func main() {
	sms := NewBulkSMSSMSSender("joorek", "warbird")
	sms.Testmode = 0
	sms.RoutingGroup = 1
	
	msg := "lol, hi. das ist ein test!"
	receivers := []string{"492161531987"}//,"491722579081"}

	err, quote := sms.GetQuote(receivers, msg)
	if err != nil {
		fmt.Println(err.String())
		return
	} 
	price := quote * 3.75 * 0.01   //mad math skills calculate price in MONEYS
	
	//we're cheap!
	if quote > 2.0 {
	    fmt.Printf("sorry, but %.2f credits (%.2f EUR) is too much for a sms!\n", quote, price)
	    return
	}
	
	fmt.Printf("this sms will cost %.4f eur\n", price)
	

    if err := sms.Send(receivers, msg); err != nil {
        fmt.Println(err.String())
        return
    }

    fmt.Println("sms sent")

}
