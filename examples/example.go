package main

import (
	"fmt"
	//"../"
	"github.com/jsz/gosms"
)

func main() {
	sms := gosms.NewBulkSMSSMSSender("USERNAME", "PASSWORD")
	sms.Testmode = 0
	sms.RoutingGroup = 2
	sms.SenderId = "Tabletten"
	
	msg := "https://github.com/jsz/gosms is awesome! -- sent from my Go"

	receivers := []string{"49178xxxxxx", "49172xxxxxxxx", }

	//quote gives you the cost of the sms in credits
	err, quote := sms.GetQuote(receivers, msg)
	if err != nil {
		fmt.Println(err)
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
        fmt.Println(err)
        return
    }

    fmt.Println("sms sent")

}
