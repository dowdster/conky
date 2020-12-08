package main

import (
        "encoding/json"
        "fmt"
        "golang.org/x/net/context"
        "golang.org/x/oauth2"
        "golang.org/x/oauth2/google"
        "google.golang.org/api/calendar/v3"
        "io/ioutil"
        "log"
        "net/http"
        "os"
        "strconv"
        "strings"
        "sync"
        "time"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
        // The file token.json stores the user's access and refresh tokens, and is
        // created automatically when the authorization flow completes for the first
        // time.
        tokFile := "token.json"
        tok, err := tokenFromFile(tokFile)
        if err != nil {
                tok = getTokenFromWeb(config)
                saveToken(tokFile, tok)
        }
        return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
        authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
        fmt.Printf("Go to the following link in your browser then type the "+
                "authorization code: \n%v\n", authURL)

        var authCode string
        if _, err := fmt.Scan(&authCode); err != nil {
                log.Fatalf("Unable to read authorization code: %v", err)
        }

        tok, err := config.Exchange(context.TODO(), authCode)
        if err != nil {
                log.Fatalf("Unable to retrieve token from web: %v", err)
        }
        return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
        f, err := os.Open(file)
        if err != nil {
                return nil, err
        }
        defer f.Close()
        tok := &oauth2.Token{}
        err = json.NewDecoder(f).Decode(tok)
        return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
        fmt.Printf("Saving credential file to: %s\n", path)
        f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
        if err != nil {
                log.Fatalf("Unable to cache oauth token: %v", err)
        }
        defer f.Close()
        json.NewEncoder(f).Encode(token)
}

var MyMaxTrimRight int = 1
var MyCurTrimRight int = 0
var MyMtx = &sync.Mutex{}

func MyTrimRightFunc(r rune) bool {
    MyCurTrimRight++
    return MyCurTrimRight!=MyMaxTrimRight
}

func MyTrimRight(s string, nchars int) string {
    if nchars<=0 {
        return s
    }

    MyMtx.Lock()

    MyMaxTrimRight = nchars
    MyCurTrimRight = 0
    substr := strings.TrimRightFunc(s, MyTrimRightFunc)

    MyMtx.Unlock()
    return substr
}

func PrintMsg(color string, msg string) {
    if color=="zzzzz" {
        os.Stdout.WriteString(color); 
        os.Stdout.WriteString(msg)
        os.Stdout.WriteString(Reset);

    } else {
        os.Stdout.WriteString(msg)
    }
}

var OldReset     = "\033[0m"
var OldBoldWhite = "\033[01m"
var OldBlack     = "\033[30m"
var OldRed       = "\033[31m"
var OldGreen     = "\033[32m"
var OldYellow    = "\033[33m"
var OldBlue      = "\033[34m"
var OldPurple    = "\033[35m"
var OldCyan      = "\033[36m"
var OldGray      = "\033[37m"
var OldBlackBkgnd= "\033[40m"
var OldDarkGray  = "\033[90m"
var OldWhite     = "\033[97m"
var OldGrayBkgnd = "\033[100m"

var BrightBlack string    = "\033[30;1m"
var BrightRed string      = "\033[31;1m"
var BrightGreen string    = "\033[32;1m"
var BrightYellow string   = "\033[33;1m"
var BrightBlue string     = "\033[34;1m"
var BrightMagenta string  = "\033[35;1m"
var BrightCyan string     = "\033[36;1m"
var BrightWhite string    = "\033[37;1m"

var Reset string          = "\033[0m"
var Black string          = "\033[0;30m"
var Red string            = "\033[0;31m"
var Green string          = "\033[0;32m"
var Yellow string         = "\033[0;33m"
var Blue string           = "\033[0;34m"
var Magenta string        = "\033[0;35m"
var Cyan string           = "\033[0;36m"
var White string          = "\033[0;37m"
var BlackBkgnd string     = "\033[0;40m"
var GrayBkgnd string      = "\033[0;100m"


var CalBoxWidth       = 140
var CalBoxGridDayWidth= 20
var CalBoxTop         = "+--------------------------------------------------------------------------------------------------------------------------------------------------+"
var CalBoxGridDay     = "+--------------------"
var CalBoxGridEmptyRow= "|                    "
var CalBoxEdgeLeft    = "|"
var CalBoxEdgeMid     = "|"
var CalBoxEdgeRight   = "|"
var CalBoxGridDayLeft = "|------------------------"
var CalBoxGridDayRight= "+------------------------|"

var Months      [12]string = [12]string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"} 
var Weekdays    [7]string= [7]string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
var EventColors [7]string= [7]string{Green, Blue, Magenta, Cyan, BlackBkgnd, GrayBkgnd, BrightWhite}

func DrawCal(calEvents [][]string, minDate time.Time, maxDate time.Time) {

   smsg := ""

   // map that matches calendar & text color
   colormap := make(map[string]string)
   curcolor := 0

   // draw calendar top box line
   PrintMsg(White, CalBoxTop + "\n")

   // output month-year
   PrintMsg(White, CalBoxEdgeLeft)
   smsg = fmt.Sprintf("%65s %9s %4s %65s", " ", Months[minDate.Month()-1], strconv.Itoa(minDate.Year()), " ")
   PrintMsg(BrightYellow, smsg)
   PrintMsg(White, CalBoxEdgeLeft + "\n")

   // draw top of day grid
   for x:=0; x<7; x++ {
       PrintMsg(White, CalBoxGridDay)
   }
   PrintMsg(White, "+\n")

   // draw days of week
   PrintMsg(White, CalBoxEdgeLeft)
   for _, day := range Weekdays {
       smsg = fmt.Sprintf("%6s%-14s", " ", day)
       PrintMsg(BrightYellow, smsg)
       fmt.Print(White, CalBoxEdgeMid)
   }
   PrintMsg(White, "\n")

   // figure out which weekday the first day in the month falls on
   nstart := 0
   for nstart=0; nstart<len(Weekdays); nstart++ {
       if Weekdays[nstart]==minDate.Weekday().String() {
           break
       }
   }

   // get last day
   nend   := maxDate.Day()
   today  := time.Now().Day()
   curday := 1

   for week:=0; week<5; week++ {
       daystart := curday
       // each day has 10 rows...
       for row:=0; row<10; row++ {
           // each row has 7 days
           for wday:=0; wday<7; wday++ {
               if row==0 {
                   // horizontal line separating previous week 
                   PrintMsg(White, CalBoxGridDay)
                   if wday==6 { 
                       PrintMsg(White, "+\n")
                   }
               } else if row==1 {
                   // day of the month goes here 
                   if (week==0 && wday<nstart) || (week==5 && wday>nend) {
                       // skip weekdays outside of the current month
                       smsg = fmt.Sprintf("%-21s", CalBoxEdgeLeft)
                       PrintMsg(White, smsg)
                       continue
                   }

                   if curday<=maxDate.Day() {
                       // in the current month, dump the day etc
                       PrintMsg(White, CalBoxEdgeLeft)
                       smsg = fmt.Sprintf("%9s%2s%9s", " ", strconv.Itoa(curday), " ")
                       if curday==today {
                            PrintMsg(BrightRed, smsg)
                       } else {
                            PrintMsg(BrightYellow, smsg)
                       }
                       curday++
                       if wday==6 {
                            PrintMsg(White, CalBoxEdgeLeft+"\n")
                       }
                   } else {
                      // we're in the next month
                      fmt.Print(White + CalBoxGridEmptyRow + Reset)
                       if wday==6 {
                           PrintMsg(White, CalBoxEdgeLeft+"\n")
                       }
                   }
               } else if row==2 {
                   // blank row (just borders)
                   PrintMsg(White, CalBoxGridEmptyRow)
                   if wday==6 {
                       PrintMsg(White, CalBoxEdgeLeft+"\n")
                   }
                   curday++
               } else {
                   if (week==0 && wday<nstart) || (week==5 && wday>nend) {
                       // skip weekdays before first day in month
                       smsg = fmt.Sprintf("%-21s", CalBoxEdgeLeft)
                       fmt.Print(White, smsg)
                       continue
                   } 
                   // calendar items go here
                   if curday<len(calEvents) && len(calEvents[curday-1])>0 {
                        evntitm := calEvents[curday-1][0]
                        // event string format "calendar_id"|"event_summary"|"event_time"
                        info := strings.SplitN(evntitm, "|", -1)
                        // find cal_id in colormap
                        if colormap[info[0]]=="" {
                            colormap[info[0]] = EventColors[curcolor]
                            curcolor++
                        }
                        PrintMsg(White, CalBoxEdgeLeft)
                        // print event info
                        smsg = fmt.Sprintf("%-5s %-14s", strings.TrimRight(info[2], " "), MyTrimRight(info[1], len(info[1])-12))
                        // set event color & print event
                        PrintMsg(colormap[info[0]], smsg)

                        // lastly delete the item
                        tmp := calEvents[curday-1][1:]
                        calEvents[curday-1] = tmp
                   } else {
                        PrintMsg(White, CalBoxGridEmptyRow)
                   }
                    
                   curday++

                   if wday==6 {
                        PrintMsg(White, CalBoxEdgeLeft+"\n")
                   }
               }
           }

           if row!=9 {
             curday = daystart
           }
       }
   }
   // draw calendar bottom box line
   PrintMsg(White, CalBoxTop+"\n")
   PrintMsg(White,"\n")
   PrintMsg(White,"\n")
   PrintMsg(White,"\n")
}

func main() {
        b, err := ioutil.ReadFile("credentials1.json")
        if err != nil {
                log.Fatalf("Unable to read client secret file: %v", err)
        }

        // If modifying these scopes, delete your previously saved token.json.
        config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
        if err != nil {
                log.Fatalf("Unable to parse client secret file to config: %v", err)
        }
        client := getClient(config)

        srv, err := calendar.New(client)
        if err != nil {
                log.Fatalf("Unable to retrieve Calendar client: %v", err)
        }

        /* get ALL user calendars */
        cals, err:= srv.CalendarList.List().Do()

        /* start=beginning of 1st day in current month*/
        startT := time.Date(time.Now().Year(), time.Now().Month(), 1, 0/*hour*/, 0/*minute*/, 0/*second*/, 0/*nsec*/, time.Now().Location())
        /* end=end of last day in current month*/
        endT   := startT.AddDate(0/*years*/, 1/*months*/, -1/*days*/).Add(time.Hour*23 + time.Minute*59 + time.Second*59)

        maxdays := endT.Day()+1
        calEvents := make([][]string, maxdays)

        /* iterate calendars */
        for _, calitem := range cals.Items {
            // enumerate events in calendar
            events, err := srv.Events.List(calitem.Id).
                                       ShowDeleted(false).
                                       SingleEvents(true).
                                       TimeMin(startT.Format(time.RFC3339)).
                                       TimeMax(endT.Format(time.RFC3339)).
                                       MaxResults(100).
                                       OrderBy("startTime").
                                       Do()

            if (err != nil) {
                log.Fatalf("Unable to enumerate events in calendar (%s).  Error: %v", calitem.Id, err)
                continue;
            }
            
            // copy events to struct
            for _, eventitem := range events.Items {
                var date time.Time
                date, err = time.Parse(time.RFC3339, eventitem.Start.DateTime)
                if err!=nil {
                    //date = eventitem.Start.Date
                    date, err = time.Parse("2006-01-02", eventitem.Start.Date)
                }
                if err!=nil {
                    log.Fatalf("Unable to extract datetime from item (%s).", eventitem.Summary, eventitem.Start.DateTime, eventitem.Start.Date)
                    continue
                }

                // build up the new event string
                s := calitem.Id + "|"
                s += eventitem.Summary + "|" 
                s += strconv.Itoa(date.Local().Hour()) + ":"
                if date.Local().Minute()<10 {
                    s += "0"
                }
                s += strconv.Itoa(date.Local().Minute())

                // add event string to the slice
                day := date.Day()-1
                calEvents[day] = append(calEvents[day], s)

                //log.Printf("Add event (%v) to Day %v\n", s, day)
            }
        } 
        DrawCal(calEvents, startT, endT)
}

