package main

// import http and gin 
// the neccessary packages to run this webservice
import (
    "net/http"
    "github.com/gin-gonic/gin"
    "fmt"
    "time"
)

// Define Structs to represent data on transactions 
type transaction struct {
    Payer     string  `json:"payer"`
    Points    int     `json:"points"`
    Timestamp string  `json:"timestamp"`
    Remaining int  `json:"-"`
}

type spend struct {
    Points  int  `json:"points"`
}

type balance struct {
    Payer     string  `json:"payer"`
    Points  int  `json:"points"`
}

var transactions = []transaction{}
var balances = []balance{}

// the format that was given in specs 
// var dateFormat = "2006-01-02T15:04:05-0700";
// var dateFormat = "2006-01-02T15:04:05";
var dateFormat = "2006-01-02T15:04:05Z";
// create routes for the webserivice to interpret endpoints
func main() {
    router := gin.Default()
    fmt.Println("hello")
    router.GET("/checkBalances", getBalances)
    router.PUT("/spend", putSpend) 
    router.POST("/pay", postPay)

    router.Run("localhost:8080")
}

// responds with transactions as JSON.
func getBalances(c *gin.Context) {
    c.IndentedJSON(http.StatusOK, balances)
}

// postPay adds a transaction from JSON received to transactions slice.
func postPay(c *gin.Context) {
    var newTransaction transaction
    // data := &E{}
    // c.Bind(data)
    // binds json to new transaction otherwise error
    if err := c.BindJSON(&newTransaction); err != nil {
        // panic(c)
        return
    }
    // Add the new transaction to the slice.
    newTransaction.Remaining = newTransaction.Points
    transactions = append(transactions, newTransaction)
    // update balances based on the transaction 
    fmt.Println(newTransaction.Points)
    updateBalances(newTransaction)
    
    c.IndentedJSON(http.StatusCreated, newTransaction)
}

// adds or updates balance 
func updateBalances(t transaction){
    found := false 
    for i:= 0; i < len(balances) && !found; i++ {
        if balances[i].Payer == t.Payer {
            found = true
            balances[i].Points += t.Points
            if t.Points < 0 {
                _,j := getOldestWithName(t.Payer)
                fmt.Println(transactions[j])
                transactions[j].Remaining += t.Points
                // fmt.Println(transactions[j].Remaining)
            }
        }
    }
    if !found{
        newBalance := balance{Payer: t.Payer, Points: t.Points}
        balances = append(balances, newBalance)
    } 
}


func putSpend(c *gin.Context) {
    var newSpend spend
    var ledger = []balance{}
    temp := transactions
    // binds json to new transaction otherwise error
    if err := c.BindJSON(&newSpend); err != nil {
        return
    }
    // while there are still points to spend 
    for newSpend.Points > 0 {
        // go through transactions 
        // and update update balance
        
        t,i := getOldestNonSpentTransaction()
        // fmt.Println(t)
        if i == -1 {
            // went into negatives revert and send 400 error 
            transactions = temp
            c.IndentedJSON(http.StatusBadRequest , gin.H{"message": "timestamp not parseable, please follow iso8601"})
            return 
        }else if i == -2{
            transactions = temp
            c.IndentedJSON(http.StatusNotAcceptable , gin.H{"message": "Went Negative"})
            return 
        }
        var newBalance balance
        newBalance, newSpend = spendWithBalance(newSpend,t,i)
        ledger = append(ledger,newBalance)
        // fmt.Println(ledger)
    } 

    for _, b := range balances {
        for j, l := range ledger{
            if b.Payer == l.Payer {
                balances[j].Points += l.Points 
            }
        }
    }


    c.IndentedJSON(http.StatusAccepted, ledger)
}

func getOldestNonSpentTransaction()(transaction , int){
    var result transaction
    var oldest = time.Now() // assuming there is no future transactions saved now
    //  should be  youngest
    index := -2
    for i, t := range transactions {
        // fmt.Println(t)
        var currentTime,err = time.Parse(dateFormat, t.Timestamp)
        if err != nil {
            fmt.Println(err)
            return t, -1
        }
        if(currentTime.Before(oldest) && t.Remaining > 0){
            result = t
            index = i
            oldest = currentTime
        }
    }

    return result, index 
}


func getOldestWithName(payer string)(transaction , int){
    var result transaction
    var oldest = time.Now() // assuming there is no future transactions saved now
    //  should be  youngest
    index := -2
    for i, t := range transactions {
        // fmt.Println(t)
        var currentTime,err = time.Parse(dateFormat, t.Timestamp)
        if err != nil {
            fmt.Println(err)
            return t, -1
        }
        if(currentTime.Before(oldest) && t.Payer == payer && t.Remaining > 0){
            result = t
            index = i
            oldest = currentTime
        }
    }

    return result, index 
}



func spendWithBalance(s spend, t transaction, i int) (balance, spend){
    // spends without going negative find difference
    difference := 0
    if s.Points >= t.Remaining{
        difference =  -1 * t.Remaining
        s.Points -= t.Remaining
        t.Remaining  = 0
    } else {
        difference =  -1*s.Points 
        t.Remaining -= s.Points
        s.Points = 0 
    }
    // updates transaction with the index provided 
    transactions[i] = t 
    // creates a balance using transaction name and difference
    b := balance{Payer:t.Payer, Points: difference}
    return b,s
}

