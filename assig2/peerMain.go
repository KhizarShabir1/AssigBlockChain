package main
import (
  "encoding/gob"
  a2 "assignment02IBC"
  "net"
  "log"
  "fmt"
  "os"
  "strconv"
  "time"
  // "bufio"

)
type addrs struct {
Name string
Addr string

}

var blockChain * a2.Block


func ListenForTransactions(addres [] addrs, name string, Address string ){
  ln, err := net.Listen("tcp", ":"+Address)

  if err != nil {
    log.Fatal(err)
  }

  for {
    conn, err := ln.Accept()
    if err != nil {
    log.Println(err)
    continue
    }
    log.Println("A client has connected with "+name+" : ", conn.RemoteAddr())

    buf := make([]byte, 4096)

    n, err := conn.Read(buf) //get address of the peer
    if err != nil || n == 0 {
        conn.Close()
        fmt.Println("Closing connection")

    }
    fmt.Println(string(buf[0:n]))
    choice:=string(buf[0:n])

    if choice=="mine"{

      var transact a2.TransSend
      dec := gob.NewDecoder(conn)
      err = dec.Decode(&transact)
      if err != nil {
      //handle error
      }

      // "validate"
      // var transact a2.TransSend
      // fmt.Println(transact.FreeCoin)
      // fmt.Println(transact.Transactions)
      // fmt.Println(transact.Sender)
      go handleValidationAndDistribution(transact, addres ) // myName,


    }else if choice=="broadcast"{
//COde of stopin validation

      bufe := make([]byte, 40960)

      n2, err := conn.Read(bufe) //get address of the peer
      if err != nil || n == 0 {
          conn.Close()
          fmt.Println("Closing connection")

      }
      hashVal:=string(bufe[0:n2])

      existHash:=a2.CheckHashExists(blockChain, hashVal)
      var tmpBlock * a2.Block
      dec := gob.NewDecoder(conn)
      err = dec.Decode(&tmpBlock)
      a2.ListBlocks(tmpBlock)

      if err != nil {
      //handle error
      }


      if existHash==false{  //check that block does not exist beffore this
        fmt.Println("Previous blockchain")
        a2.ListBlocks(blockChain)

        blockChain=tmpBlock
        fmt.Println("new block chai is ")
        a2.ListBlocks(blockChain)
        fmt.Println("new block chai is ")
        //Now lets broadCast to the connections

           for o:=0;o<len(addres);o++{ //lets broadcast blockchain

             conn, err := net.Dial("tcp", "localhost:"+addres[o].Addr)
             if err != nil {
             //handle error
             }

             conn.Write([]byte("broadcast"))
             time.Sleep(1 * time.Second)

             conn.Write([]byte(hashVal))
             time.Sleep(1 * time.Second)

               gobEncoder := gob.NewEncoder(conn)
               err = gobEncoder.Encode(blockChain)
               if err != nil {
                 log.Println(err)
               }

           }

      }else{
        //Just recieve and dont broadCast

        fmt.Println("Block Received but not broadCasted")
        // a2.ListBlocks(blockChain)

      }



    }// Broad ends here

    // s := strings.Split(string(buf[0:n]), ":")
    // nd := addrs{Name: s[0], Addr:s[1]}
    // addresses = append(addresses,nd)



  }


}



func handleValidationAndDistribution(transact a2.TransSend,  addres [] addrs ){  // myName string,
  conn, err := net.Dial("tcp", "localhost:7000")
  if err != nil {
  //handle error
  }
  conn.Write([]byte("validate"))
  time.Sleep(1 * time.Second)
  gobEncoder := gob.NewEncoder(conn)
  err = gobEncoder.Encode(transact)


  buf := make([]byte, 4096)
  n, err := conn.Read(buf) //get validation
  if err != nil || n == 0 {
    conn.Close()
    fmt.Println("Closing connection")

  }
  // "yes"
  // "no"
  recievedValidation:=string(buf[0:n])
  if recievedValidation=="yes"{
    // var transact a2.TransSend
    // fmt.Println(transact.FreeCoin)
    // fmt.Println(transact.Transactions)
    // fmt.Println(transact.Sender)
    var transaction a2.Trans
    transaction.Transactions=append(transaction.Transactions,transact.Transactions)
    transaction.FreeCoin=append(transaction.FreeCoin,transact.FreeCoin)
    transaction.Transactions=append(transaction.Transactions,transact.Sender)
    transaction.FreeCoin=append(transaction.FreeCoin,75) //CoinBase transaction

    transaction.NoOfTrans=2
    // AddMoneyToValidStore(validStore ,"satoshi" , 100 )
    var hashVal string

    hashVal,blockChain = a2.InsertBlock(transaction, blockChain)

    //   fmt.Print(addres[i].Addr)
    //   fmt.Print(addres[i].Name)
    // addres [] addrs
    conn, err := net.Dial("tcp", "localhost:7000") //sendin blockChain for satoshi
    if err != nil {
    //handle error
    }

    conn.Write([]byte("broadcast"))
    time.Sleep(1 * time.Second)
      gobEncoder := gob.NewEncoder(conn)
      err = gobEncoder.Encode(blockChain)
      if err != nil {
        log.Println(err)
      }

   for o:=0;o<len(addres);o++{ //lets broadcast blockchain
     conn, err := net.Dial("tcp", "localhost:"+addres[o].Addr)
     if err != nil {
     //handle error
     }

     conn.Write([]byte("broadcast"))
     time.Sleep(1 * time.Second)
     conn.Write([]byte(hashVal))
     time.Sleep(1 * time.Second)

       gobEncoder := gob.NewEncoder(conn)
       err = gobEncoder.Encode(blockChain)
       if err != nil {
         log.Println(err)
       }

   }
   fmt.Println("Block Broadcated from -> "+transact.Sender)

  }else if recievedValidation=="no" {
    fmt.Println("Invalid Transation recieved from -> "+transact.Sender)

  }



}// function ens here





func main() {



name:= os.Args[1]
Address:=os.Args[2]
// var addresses [] addrs
conn, err := net.Dial("tcp", "localhost:7000")
if err != nil {
//handle error
}
log.Println("A client has connected with satoshi", conn.RemoteAddr())
fmt.Printf(" client connected with satoshi ")

conn.Write([]byte(name+":"+Address))
// var recvdBlock * a2.Block
dec := gob.NewDecoder(conn)
err = dec.Decode(&blockChain)

if err != nil {
//handle error
}
var addres [] addrs
// dec := gob.NewDecoder(conn)
// gob.Register(net.TCPConn)
err = dec.Decode(&addres)
if err != nil {
//handle error
}

a2.ListBlocks(blockChain)
fmt.Println("My connections are following Broadcated from -> ")
for i:=0;i<len(addres);i++{

  fmt.Print(addres[i].Addr)
  fmt.Print(" : ")
  fmt.Print(addres[i].Name)
  fmt.Println("\n")

}

go ListenForTransactions(addres , name , Address)
for {

var option int
option= -1
fmt.Println("Enter 1 for transaction : ")
_, err= fmt.Scan(&option)
if option==1{
    // reader := bufio.NewReader(os.Stdin)

    var coins int
    fmt.Println("Please enter no of FreeCoins do you want to send : ")
    _, err= fmt.Scan(&coins)
    for i:=0;i<len(addres);i++{
      fmt.Println("Enter "+strconv.Itoa(i)+" to send to "+addres[i].Name)
    }
    var tranC int

    _, err= fmt.Scan(&tranC)
    conn, err := net.Dial("tcp", "localhost:7000")
    if err != nil {
    //handle error
    }
    conn.Write([]byte("chooseMiner"))
    time.Sleep(1 * time.Second)
    //Encodin transaction
    var transact a2.TransSend
    transact.Transactions=addres[tranC].Name
    transact.FreeCoin=coins
    transact.Sender=name

    gobEncoder := gob.NewEncoder(conn)
    err = gobEncoder.Encode(transact)
    if err != nil {
      log.Println(err)
    }
    // conn.Close()



}



}



}
