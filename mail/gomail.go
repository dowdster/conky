package main

import (
        "encoding/json"
        "fmt"
        "io/ioutil"
        "log"
        "net/http"
        "os"
        "sort"
        "strconv"
        "strings"
        "time"
        "golang.org/x/net/context"
        "golang.org/x/oauth2"
        "golang.org/x/oauth2/google"
        "google.golang.org/api/gmail/v1"
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

func getargs() map[string]string {

    m := make(map[string]string)
    for _, itm:=range os.Args {
        a := strings.Split(itm,"=")
        if len(a)>1 {
            //log.Printf("Adding arg '%v'='%v'\n", a[0], a[1])
            m[a[0]] = a[1]
        }
    }
    if m["search"]=="" {
        if m["label"]!="" {
            m["search"] = "in:" + m["label"] + " "
        } 
        if m["status"]!="" {
            m["search"] += "is:" + m["status"] + " "
        }
        if m["date"]!="" {
            m["search"] += "newer_than:" + m["date"] + " "
        }
    }
    if m["user"]=="" {
        m["user"] = "me"
    }
    if m["maxmessages"]=="" {
        m["maxmessages"]=strconv.Itoa(100)
    }
    if m["verbose"]=="" {
        m["verbose"] = "false"
    } else {
        m["verbose"] = strings.ToLower(m["verbose"])
    }
    return m
}

type msgdata struct {
    key          int64
    date         time.Time
    subject      string
    from         string
}

func getlabels() {
        /*
        r, err := srv.Users.Labels.List(args["user"]).Do()
        if err != nil {
                log.Fatalf("Unable to retrieve labels: %v", err)
        }
        if len(r.Labels) == 0 {
                fmt.Println("No labels found.")
                return
        }
        fmt.Println("Labels:")
        for _, l := range r.Labels {
                ll, lerr := srv.Users.Labels.Get(args["user"], l.Id).Do()
                if lerr != nil {
                        log.Fatalf("Unable to get label (%v): %v", l.Name, lerr)
                }
                fmt.Printf("- %s Tot: %v, Unread: %v, ThreadTot: %v, ThreadUnread: %v\n", 
                            ll.Id,
                            ll.MessagesTotal,
                            ll.MessagesUnread,
                            ll.ThreadsTotal,
                            ll.ThreadsUnread)
        }
        */
}

func gettime(dtin string) time.Time {
    dtout, err := time.Parse(time.RFC3339, "0001-01-01T00:00:00Z")

    if len(dtin)>0 {
        //
        // dates can have the timezone specifier OR NOT...
        // Thu, 16 Apr 2020 00:12:20 +0000 (UTC)
        // Mon, 6 Apr 2020 12:30:06 -0500
        //
        a     := strings.Split(dtin, "(")
        dttmp := a[0]
    
        //
        // Dates can also not include the DayofWeek eg. Mon,
        //
        a      = strings.Split(dttmp, ",")
        if len(a)==1 {
            dttmp = a[0]
        } else {
            dttmp = a[1]
        }

        dtout,err = time.Parse("2 Jan 2006 15:04:05 -0700", strings.TrimSpace(dttmp))
        if err!=nil {
            fmt.Printf("Unable to parse date '%v'.  Error: %v\n", dtin, err)
        }
    }
    return dtout.Local()
}

func getfrom(fromin string, friendly bool) string {
    fromout := ""

    if len(fromin)>0 {
        // 
        // sender strings can look like any of these:
        // "Sam Morgan, EXACT Sports" <director@exactsports.com> 
        // Digital Federal Credit Union <dcu@dcu.org> Your DCU eStatement is ready to view
        // dcu@dcu.org 
        //
        // try and split the friendly name & <emailaddress> into elements of an array
        arr := strings.Split(fromin, "<")
        if len(arr)>1 {
            //
            // The input from string looks like one of these:
            // "Sam Morgan, EXACT Sports" <director@exactsports.com> 
            // Digital Federal Credit Union <dcu@dcu.org> Your DCU eStatement is ready to view
            //
            // trim spaces to left & right of friendly
            arr[0] = strings.TrimSpace(arr[0])
            // ...and surround with "" if not already
            if strings.Index(arr[0], "\"")==-1 {
                arr[0] = "\"" + arr[0] + "\""
            }
            // remove trailing > from emailaddress 
            arr[1] = strings.Replace(arr[1], ">", "", 1)
            // ...and trim spaces left & right
            arr[1] = strings.TrimSpace(arr[1])
        } else {
            //
            // Probably the input from string looks like this:
            // dcu@dcu.org 
            arr = make([]string,2)
            arr[0] = "\"" + strings.TrimSpace(fromin) + "\""
            arr[1] = arr[0]
        }
        if friendly {
            fromout = arr[0]
        } else {
            fromout = arr[1]
        }
    }
    return fromout
}

func main() {
        // process command line args
        args := getargs()
        verbose := args["verbose"]=="true"
        
        if verbose {
            log.Print("Starting up....\n")
        }

        b, err := ioutil.ReadFile("credentials.json")
        if err != nil {
                log.Fatalf("Unable to read client secret file: %v", err)
        }

        // If modifying these scopes, delete your previously saved token.json.
        if verbose {
            log.Print("Authenticating & Authorizing....\n")
        }
        config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
        if err != nil {
                log.Fatalf("Unable to parse client secret file to config: %v", err)
        }
        client := getClient(config)

        if verbose {
            log.Print("Creating GMAIL client\n")
        }
        srv, err := gmail.New(client)
        if err != nil {
                log.Fatalf("Unable to retrieve Gmail client: %v", err)
        }

        maxmsgs,_ := strconv.Atoi(args["maxmessages"])
        totmsgs   := 0
        ptkn      := ""
        mymsgs    := make([]msgdata, 0, maxmsgs)

        var md msgdata
  
        for {
            if totmsgs>=maxmsgs {
                break
            }

            req  := srv.Users.Messages.List(args["user"]).MaxResults(int64(maxmsgs)).Q(args["search"])
            if ptkn!="" {
                req.PageToken(ptkn)
            }

            if verbose {
                log.Printf("Requesting chunk of messages matching query '%v'...\n", args["search"])
            }
            m, err := req.Do()
            if err != nil {
                log.Fatalf("Unable to get messages: %v", err)
            }

            if verbose {
                log.Printf("Processing %v messages...\n", len(m.Messages))
            }

            for _, msg := range m.Messages {
                msgdata, err := srv.Users.Messages.Get(args["user"], msg.Id).Do()
                if err != nil {
                    log.Printf("Unable to get message (%v): %v", msg.Id, err)
                    continue;
                }
            
                // reset our struct... 
                md.date, err = time.Parse(time.RFC3339, "0001-01-01T00:00:00Z"); md.from=""; md.subject=""; md.key=0;

                // load struct with values from cur msg
                md.key = msgdata.InternalDate
                if md.key==0 {
                    if verbose {
                        log.Printf("Unable to get message %v internaldate", msg.Id)
                    }
                    continue;
                }
                
                for _, hdr:=range msgdata.Payload.Headers {
                    //log.Printf("Hdr: %v, Value: %v\n", hdr.Name, hdr.Value)
                    if hdr.Name == "Date" {
                        md.date    = gettime(hdr.Value)
                    } else if hdr.Name == "From" {
                        md.from = getfrom(hdr.Value, true/*friendly*/)
                    } else if hdr.Name == "Subject" {
                        md.subject = hdr.Value
                    }

                    if !md.date.IsZero() && md.from!="" && md.subject!="" {
                        break
                    }
                }

                //fmt.Printf("%d) %v %v %v\n", totmsgs, md.date, md.from, md.subject)

                if !md.date.IsZero() && md.from!="" && md.subject!="" && md.key!=0 {
                    mymsgs = append(mymsgs, md) 
                    //fmt.Printf("%d) %v %v %v\n", totmsgs, dt, fm, sub)
                    totmsgs++
                } else if verbose {
                    fmt.Printf("INCOMPLETE MSG:  d:'%v' s:'%v' fm:'%v' k:'%v'\n", 
                                md.date, md.subject, md.from, md.key)
                }

                if totmsgs>=maxmsgs {
                    break
                }
            }

            if len(m.Messages)<maxmsgs {
                break
            }
    }

    if verbose {
       log.Printf("Sorting %v messages...\n", len(mymsgs))
    }
    sort.SliceStable(mymsgs, func(i, j int) bool {
        return mymsgs[i].key > mymsgs[j].key
    })

    if verbose {
        log.Printf("Displaying %v messages...\n", len(mymsgs))
    }
    fmt.Print("----------------------------------------------------------------------------\n")
    idx := 1
    for _, m:=range mymsgs {
        fmt.Printf("%02d) %v %v %v\n", 
                    idx, 
                    m.date.Format("Mon Jan 2 03:04PM"),
                    m.from, 
                    m.subject)
        idx++
    }
    fmt.Print("----------------------------------------------------------------------------\n")
}

