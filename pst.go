package main

import (
	"fmt"
	"log"
	"net/http"
)

type Paste struct {
	Id   string
	Body []byte
}

func main() {
	err := http.ListenAndServe(":8080", handler())
	if err != nil {
		log.Fatal(err)
	}
}

func handler() http.Handler {
	r := http.NewServeMux()
	r.HandleFunc("/", home)
	return r
}

func home(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		fmt.Fprintf(w, Help)
	case "POST":
		s := r.FormValue("f1")
		fmt.Fprintf(w, s)
	}
}

var Help string = `
pst(1)                               PST                                  pst(1)

NAME

    pst: command line pastebin.


TL;DR

    ~$ echo Hello world. | curl -F 'f:1=<-' pst
    http://pst/fpW


GET

    pst/ID
        raw

    pst/ID/
        default syntax language (by filetype, if provided)
        append #n-LINENO to link directly to a particular line
        uses pygments, see pygments documentation for details

    pst/ID/LANG
        explicitly set language

    pst/ID+
        console highlighting (default)

    pst/ID+LANG
        console highlighting (explicitly set language)


POST

    pst/

        f:N    contents or attached file.

    where N is a unique number within request. (This allows you to post
    multiple files at once.)

    returns: http://pst/id for N in request


DELETE

    pst/ID
        delete ID.


EXAMPLES

    Anonymous, unnamed paste, two ways:

        cat file.ext | curl -F 'f:1=<-' pst
        curl -F 'f:1=@file.ext' pst


    Delete ID, two ways:

        curl -n -X DELETE pst/ID
        curl -F 'rm=ID' pst


CLIENT

    A client is maintained at pst/client

        curl pst/client > pst
        chmod +x pst
        ./pst -h

    Or if you wish, paste the following function into $HOME/.bashrc:

        pst() {
            local opts
            local OPTIND
            [ -f "$HOME/.netrc" ] && opts='-n'
            while getopts ":hd:i:n:" x; do
                case $x in
                    h) echo "pst [-d ID] [-i ID] [-n N] [opts]"; return;;
                    d) $echo curl $opts -X DELETE pst/$OPTARG; return;;
                    i) opts="$opts -X PUT"; local id="$OPTARG";;
                    n) opts="$opts -F read:1=$OPTARG";;
                esac
            done
            shift $(($OPTIND - 1))
            [ -t 0 ] && {
                local filename="$1"
                shift
                [ "$filename" ] && {
                    curl $opts -F f:1=@"$filename" $* pst/$id
                    return
                }
                echo "^C to cancel, ^D to send."
            }
            curl $opts -F f:1='<-' $* pst/$id
        }

    Then open a new shell and type ` + "`pst -h`" + `


CAVEATS:
    Paste at your risk. Be nice please.

    The codebase for pst is intended to be free and open-source. It is not
    published at the moment because the author doesn't want to publish code
    that isn't pretty and pleasant (and is also deeply lazy about getting it
    there).
`
