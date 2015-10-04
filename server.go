package main

import (
	"fmt"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"net"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"strconv"
	"strings"
	)

	
var record[20] string
var uninvest[20] float64
var count int



type Request1 struct {
						Symb string
						Budget float64
                    }
type Response1 struct {
                         Remainder float64
						 Tradeid int 
						 Rep string
					 }
type Request2 struct {
						Id int
					}
type Response2 struct {
                         Rep string
						 Uninvested float64
					 }

type Stk struct {
		List struct{ 
		Resources [] struct{
          Resource struct  {
               Fields struct {
                             Name string `json:"name"`
							Price string `json:"price"`
							Symbol string `json:"symbol"`
							Ts string `json:"ts"`
							Type string `json:"type"`
							UTCtime string `json:"utctime"`
							Volume string `json:"volume"`
							}`json:"fields"`
                         }`json:"resource"`
					   }`json:"resources"`
                }`json:"list"`
             }

type Stock int

func (t *Stock) Stkresponse(args *Request1, reply *Response1) error {
     
	 var company string
	 var p Stk
	 var remaining float64
	 remaining=0.0
	 var k string
	 var cnt int 
	 var divider float64
	 divider=100.0
	 cnt = 1
	 var m string
	 company = args.Symb
	 budg := args.Budget
	 s:=strings.Split(company,",")
	 for _, cell := range s {
	x:=strings.Split(cell,":")
	y:=x[0]
	temp1 := x[1]
	z:=strings.Split(temp1,"%")
	temp2,err := strconv.ParseFloat(z[0],64)
	if err != nil {
        fmt.Printf("%s", err)
        os.Exit(1)
    } 
	temp3 := temp2/divider
	temp4 := temp3*budg
	a := fmt.Sprint("http://finance.yahoo.com/webservice/v1/symbols/",y,"/quote?format=json")
	response, err := http.Get(a)
    if err != nil {
        fmt.Printf("%s", err)
        os.Exit(1)
    } else {
        defer response.Body.Close()
        contents, err := ioutil.ReadAll(response.Body)
        if err != nil {
            fmt.Printf("%s", err)
            os.Exit(1)
        }
      err=json.Unmarshal(contents, &p)
		   //if err != nil {
		//	panic(err)
		//}
		cost, err :=strconv.ParseFloat(p.List.Resources[0].Resource.Fields.Price,64)
        if err != nil {
        fmt.Printf("%s", err)
        os.Exit(1)
    } 
	number := temp4/cost
	var q int = int(number)
	temp5 := cost*float64(q)
	remains:=temp4-temp5	
	remaining=remaining+remains
	if(cnt==1){
	k = fmt.Sprint(y,":",strconv.Itoa(q),":","$",p.List.Resources[0].Resource.Fields.Price)
	m=k
	cnt=0
	} else {
	m = fmt.Sprint(k,",",y,":",strconv.Itoa(q),":","$",p.List.Resources[0].Resource.Fields.Price)
	k=m
	}
	
	}//ELSE ENDS
	    
            }//FOR ENDS
		record[count]= m
		uninvest[count]=remaining
		count=count+1
		reply.Rep=m
		reply.Tradeid=count
		reply.Remainder=remaining

	 return nil
}



func (t *Stock) Stkidresponse(args *Request2, reply *Response2) error {
     
	 var sign string
	 var m string
	 var k string
	 var p Stk
	 var cnt int
	 cnt=1
	 var portfolio int
	 portfolio = args.Id-1
	 b:=record[portfolio]
	 d:=strings.Split(b,",")
	 
	 for _, cell := range d {
	 s:=strings.Split(cell,":")
	 r:=strings.Split(s[2],"$")
	 price_before, err:=strconv.ParseFloat(r[1],64)
	 if err != nil {
        fmt.Printf("%s", err)
        os.Exit(1)
    } 
	 fmt.Println(price_before)
	 y:=s[0]
	 a := fmt.Sprint("http://finance.yahoo.com/webservice/v1/symbols/",y,"/quote?format=json")
	 fmt.Println(a)
	 response, err := http.Get(a)
     if err != nil {
        fmt.Printf("%s", err)
        os.Exit(1)
    } else {
        defer response.Body.Close()
        contents, err := ioutil.ReadAll(response.Body)
        if err != nil {
            fmt.Printf("%s", err)
            os.Exit(1)
        }
        err=json.Unmarshal(contents, &p)
		   //if err != nil {
		//	panic(err)
		//}
		cost, err :=strconv.ParseFloat(p.List.Resources[0].Resource.Fields.Price,64)
        if err != nil {
        fmt.Printf("%s", err)
        os.Exit(1)
    } 
	if(cost>price_before){
	sign="+"
	}else if(cost<price_before){
	sign="-"
	}else{
	sign="="
	}
	if(cnt==1){
	k = fmt.Sprint(s[0],":",s[1],":",sign,"$",cost)
	m=k
	cnt=0
	}else {
	m =fmt.Sprint(k,",",s[0],":",s[1],":",sign,"$",cost)
	}
	}
	}
	fmt.Println(sign)
	fmt.Println(k)
	reply.Rep = m
	reply.Uninvested=uninvest[portfolio]
	 return nil
}



func main() {

	stock := new(Stock)
	rpc.Register(stock)
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":1455")
	checkError(err)
    listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		jsonrpc.ServeConn(conn)
	}

}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}