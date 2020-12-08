package main

import (
        "fmt"
        //"log"
        //"os"
        "time"
        "strconv"
        //"io/ioutil"
        //"net/http"
        //"encoding/json"
        //"golang.org/x/net/context"
        //"golang.org/x/oauth2"
        //"golang.org/x/oauth2/google"
        //"google.golang.org/api/calendar/v3"
)

var Reset             = "\033[0m"
var Red               = "\033[31m"
var Green             = "\033[32m"
var Yellow            = "\033[33m"
var Blue              = "\033[34m"
var Purple            = "\033[35m"
var Cyan              = "\033[36m"
var Gray              = "\033[37m"
var White             = "\033[97m"

var CalBoxWidth       = 84
var CalBoxGridDayWidth= 12
var CalBoxTop         = "+-----------------------------------------------------------------------------------+"
var CalBoxGridDay     = "+-----------"
var CalBoxGridEmptyRow= "|           "
var CalBoxEdgeLeft    = "|"
var CalBoxEdgeMid     = "|"
var CalBoxEdgeRight   = "|"
var CalBoxGridDayLeft = "|-----------"
var CalBoxGridDayRight= "+-----------|"

var Months [12]string = [12]string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"} 
var Weekdays [7]string= [7]string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}

func main() {
    /* start=beginning of 1st day in current month*/
    minDate := time.Date(time.Now().Year(), time.Now().Month(), 1, 0/*hour*/, 0/*minute*/, 0/*second*/, 0/*nsec*/, time.Now().Location())
    /* end=end of last day in current month*/
    maxDate := minDate.AddDate(0/*years*/, 1/*months*/, -1/*days*/).Add(time.Hour*23 + time.Minute*59 + time.Second*59)
    fmt.Print(minDate.String() + "\n")
    fmt.Print(maxDate.String() + "\n")

    // draw calendar top box line
    fmt.Print(Gray + CalBoxTop + Reset + "\n")

    // output month-year
    fmt.Print(Gray + CalBoxEdgeLeft + Reset)
    fmt.Print(Yellow); fmt.Printf("%30s %9s %4s %37s", " ", Months[minDate.Month()-1], strconv.Itoa(minDate.Year()), " "); fmt.Print(Reset);
    fmt.Print(Gray + CalBoxEdgeLeft + Reset + "\n")

    // draw top of day grid
    fmt.Print(Gray)
    for x:=0; x<7; x++ {
        fmt.Print(CalBoxGridDay)
    }
    fmt.Print("+" + Reset + "\n")

    // draw days of week
    fmt.Print(Gray + CalBoxEdgeLeft + Reset)
    for _, day := range Weekdays {
        fmt.Print(Yellow); fmt.Printf("%1s%-9s%1s", " ", day, " "); fmt.Print(Reset);
        fmt.Print(Gray + CalBoxEdgeMid + Reset)
    }
    fmt.Print("\n")

    // figure out which weekday the first day in the month falls on
    nstart := 0
    for nstart=0; nstart<len(Weekdays); nstart++ {
        if Weekdays[nstart]==minDate.Weekday().String() {
            break
        }
    }

    monthday := 1
    for week:=0; week<5 && monthday<maxDate.Day(); week++ {
        // each day has 10 rows...
        for row:=0; row<10; row++ {
            // each row has 7 days
            for wday:=0; wday<7; wday++ {
                if row==0 {
                    // horizontal line separating previous week 
                    fmt.Print(CalBoxGridDay)
                    if wday==6 { 
                        fmt.Print("+" + Reset + "\n")
                    }
                } else if row==1 {
                    // day of the month goes here 
                    if week==0 && wday<nstart {
                        // skip weekdays before first day in month
                        fmt.Print(Gray); fmt.Printf("%-12s", CalBoxEdgeLeft); fmt.Print(Reset);
                    } else if monthday<=maxDate.Day() {
                        // in the current month, dump the day etc
                        fmt.Print(Gray + CalBoxEdgeLeft + Reset)
                        fmt.Print(Yellow); fmt.Printf("%4s%2s%5s", " ", strconv.Itoa(monthday), " "); fmt.Print(Reset);
                        monthday++
                        if wday==6 {
                            fmt.Print(Gray + CalBoxEdgeLeft + "\n" + Reset)
                        }
                    } else {
                        // we're in the next month
                        fmt.Print(Gray + CalBoxGridEmptyRow + Reset)
                        if wday==6 {
                            fmt.Print(Gray + CalBoxEdgeLeft + "\n" + Reset)
                        }
                    }
                } else if row==2 {
                    // blank row (just borders)
                    fmt.Print(Gray);
                    fmt.Print(CalBoxGridEmptyRow)
                    if wday==6 {
                        fmt.Print(Gray + CalBoxEdgeLeft + "\n" + Reset)
                    }
                } else {
                    // calendar items go here
                    // blank row (just borders)
                    fmt.Print(Gray);
                    fmt.Print(CalBoxGridEmptyRow)
                    if wday==6 {
                        fmt.Print(Gray + CalBoxEdgeLeft + "\n" + Reset)
                    }
                }
            }
        }
    }
    // draw calendar bottom box line
    fmt.Print(Gray + CalBoxTop + Reset + "\n")
    fmt.Print("\n")
    fmt.Print("\n")
    fmt.Printf("MaxDate.Day: %d\n", maxDate.Day())
    fmt.Print("\n")
}



