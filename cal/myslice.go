package main

import (
    "bytes"
    "fmt"
    "log"
    "os"
    "strings"
    "sync"
    "time"
)

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

func testtrim() {
    str := "This is some random text"
    sub1:= MyTrimRight(str, 10)
    sub2:= MyTrimRight(sub1, 5)
    log.Printf("Orig: %s\nSub1: %s\nSub2: %s\n", str, sub1, sub2)

    str = "Mobility"
    sub1 = MyTrimRight(str, len(str)-12)
    log.Printf("Orig: %s\nSub1: %s\n", str, sub1)
}

func testoutputfile() {
    f, _ := os.Create("/tmp/gocal.log")
    defer f.Close()
    f.WriteString("Hello world\n")
    f.WriteString("Goodbye world\n")
}

func testslice() {
    s := make([]string, 3)
    log.Println("emp:", s)

    s[0] = "a"
    s[1] = "b"
    s[2] = "c"
    log.Println("set:", s)
    log.Println("get:", s[2])
    log.Println("len:", len(s))

    s = append(s, "d")
    s = append(s, "e", "f")
    log.Println("apd:", s)

    s1 := []int{}
    log.Printf("len=%d, cap=%d %v\n", len(s1), cap(s1), s1)

    s1 = make([] int, 10)
    s1[0] = 5
    log.Printf("len=%d, cap=%d %v\n", len(s1), cap(s1), s1)

    s1 = s1[:4]
    s1[1] = 5
    s1[2] = 5
    log.Printf("len=%d, cap=%d %v\n", len(s1), cap(s1), s1)

    s1 = s1[:6]
    s1[3] = 5
    s1[4] = 5
    s1[5] = 5
    log.Printf("len=%d, cap=%d %v\n", len(s1), cap(s1), s1)

    s1 = s1[:10]
    s1[6] = 5
    s1[7] = 5
    s1[8] = 5
    s1[9] = 5
    log.Printf("len=%d, cap=%d %v\n", len(s1), cap(s1), s1)

    //resize
    daycap := cap(s1)
    tmp := make([]int, daycap+10)
    if daycap>0  {
        // copy existing items
        for i:=range s1 {
            tmp[i] = s1[i]
        }
        s1 = tmp
    }
    s1 = s1[:11]
    s1[10] = 6
    s1 = s1[:12]
    s1[11] = 6
    log.Printf("len=%d, cap=%d %v\n", len(s1), cap(s1), s1)

}

func testlogger() {
    var (
        buf bytes.Buffer
        logger = log.New(&buf, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
    )

    date := "2020-06-28"
    _,err := time.Parse("2006-01-2", date)
    if err!=nil {
        logger.Fatal("Unable to parse date (%v).  Err: %v", date, err)
    } else {
        logger.Println("Parse date (%v) OK", "2020-06-06")
    }
    tz := time.Now().Local().String()
    logger.Println("TZ: %v", tz)
    zone,offset := time.Now().Local().Zone()
    logger.Println("Z: %v, Offset: %v", zone, offset)

    fmt.Print(&buf)
}

var Reset     = "\033[0m"
var BoldWhite = "\033[01m"
var Black     = "\033[30m"
var Red       = "\033[31m"
var Green     = "\033[32m"
var Yellow    = "\033[33m"
var Blue      = "\033[34m"
var Purple    = "\033[35m"
var Cyan      = "\033[36m"
var Gray      = "\033[37m"
var BlackBkgnd= "\033[40m"
var DarkGray  = "\033[90m"
var White     = "\033[97m"

func teststdout() {
    // general make sure it works
    os.Stdout.WriteString("your 1st message here\n")
    fmt.Fprintf(os.Stdout, "your 2nd message here\n")
    // try different colors
    os.Stdout.WriteString(BoldWhite); os.Stdout.WriteString("os.Stdout.WriteString BoldWhite\n");os.Stdout.WriteString(Reset);
    os.Stdout.WriteString(Black); os.Stdout.WriteString("os.Stdout.WriteString Black\n");os.Stdout.WriteString(Reset);
    os.Stdout.WriteString(Red); os.Stdout.WriteString("os.Stdout.WriteString Red\n");os.Stdout.WriteString(Reset);
    os.Stdout.WriteString(Green); os.Stdout.WriteString("os.Stdout.WriteString Green\n");os.Stdout.WriteString(Reset);
    os.Stdout.WriteString(Yellow); os.Stdout.WriteString("os.Stdout.WriteString Yellow\n");os.Stdout.WriteString(Reset);
    os.Stdout.WriteString(Blue); os.Stdout.WriteString("os.Stdout.WriteString Blue\n");os.Stdout.WriteString(Reset);
    os.Stdout.WriteString(Purple); os.Stdout.WriteString("os.Stdout.WriteString Purple\n");os.Stdout.WriteString(Reset);
    os.Stdout.WriteString(Cyan); os.Stdout.WriteString("os.Stdout.WriteString Cyan\n");os.Stdout.WriteString(Reset);
    os.Stdout.WriteString(BlackBkgnd); os.Stdout.WriteString("os.Stdout.WriteString BlackBkgnd\n");os.Stdout.WriteString(Reset);
    os.Stdout.WriteString(DarkGray); os.Stdout.WriteString("os.Stdout.WriteString DarkGray\n");os.Stdout.WriteString(Reset);
    os.Stdout.WriteString(White); os.Stdout.WriteString("os.Stdout.WriteString White\n");os.Stdout.WriteString(Reset);

    //
    fmt.Fprintf(os.Stdout, "\nyour 2nd message here\n")
    os.Stdout.WriteString(BoldWhite); fmt.Fprintf(os.Stdout, "fmt.Fprintf BoldWhite\n");os.Stdout.WriteString(Reset);
    os.Stdout.WriteString(Black); fmt.Fprintf(os.Stdout, "fmt.Fprintf Black\n");os.Stdout.WriteString(Reset);
    os.Stdout.WriteString(Red); fmt.Fprintf(os.Stdout, "fmt.Fprintf Red\n");os.Stdout.WriteString(Reset);
    os.Stdout.WriteString(Green); fmt.Fprintf(os.Stdout, "fmt.Fprintf Green\n");os.Stdout.WriteString(Reset);
    os.Stdout.WriteString(Yellow); fmt.Fprintf(os.Stdout, "fmt.Fprintf Yellow\n");os.Stdout.WriteString(Reset);
    os.Stdout.WriteString(Blue); fmt.Fprintf(os.Stdout, "fmt.Fprintf Blue\n");os.Stdout.WriteString(Reset);
    os.Stdout.WriteString(Purple); fmt.Fprintf(os.Stdout, "fmt.Fprintf Purple\n");os.Stdout.WriteString(Reset);
    os.Stdout.WriteString(Cyan); fmt.Fprintf(os.Stdout, "fmt.Fprintf Cyan\n");os.Stdout.WriteString(Reset);
    os.Stdout.WriteString(BlackBkgnd); fmt.Fprintf(os.Stdout, "fmt.Fprintf BlackBkgnd\n");os.Stdout.WriteString(Reset);
    os.Stdout.WriteString(DarkGray); fmt.Fprintf(os.Stdout, "fmt.Fprintf DarkGray\n");os.Stdout.WriteString(Reset);
}

func main() {
    teststdout()
}

