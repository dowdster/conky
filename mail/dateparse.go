package main

import (
        "log"
        "strconv"
        "strings"
        "time"
)


func main() {

    a := [11]string{
        "Wed, 24 Jun 2020 23:38:48 +0000 (UTC)", 
        "Wed, 3 Jun 2020 13:37:18 -0400",
        "Sun, 31 May 2020 00:42:05 +0000", 
        "Sun, 31 May 2020 00:40:12 +0000",
        "Tue, 5 May 2020 12:31:21 -0400",
        "Thu, 16 Apr 2020 22:27:56 +0000 (UTC)", 
        "Thu, 16 Apr 2020 09:40:47 +0000",
        "Thu, 16 Apr 2020 00:12:20 +0000 (UTC)",
        "Mon, 6 Apr 2020 12:30:06 -0500",
        "Mon, 6 Apr 2020 12:29:58 -0500 ",
        "25 Jun 2020 01:04:21 -0400"}

    log.Print("Starting up....\n")

    log.Printf("Crazy concat strings:\n")
    for _, itm:=range a {
        a     := strings.Split(itm, "(")
        itm    = a[0]
        a      = strings.Split(itm, ",")
        if len(a)==1 {
            itm= a[0]
        } else {
            itm= a[1]
        }
        t,_   := time.Parse("2 Jan 2006 15:04:05 -0700", strings.TrimSpace(itm))
        t      = t.Local()

        s     := t.Weekday().String()[0:3] + " " +
                t.Month().String()[0:3] + " " + 
                strconv.Itoa(t.Day()) + " " +
                strconv.Itoa(t.Hour()) + ":" +
                strconv.Itoa(t.Minute()) 
        log.Printf("'%v'----->'%v'\n", itm, s)
    }

    log.Printf("\nTime.Format:\n")
    for _, itm:=range a {
        a     := strings.Split(itm, "(")
        itm    = a[0]
        a      = strings.Split(itm, ",")
        if len(a)==1 {
            itm= a[0]
        } else {
            itm= a[1]
        }
        t,_   := time.Parse("2 Jan 2006 15:04:05 -0700", strings.TrimSpace(itm))
        t    = t.Local()
        s   := t.Format("Mon Jan 2 03:04PM") 
        log.Printf("'%v'----->'%v'\n", itm, s)
    }
}

