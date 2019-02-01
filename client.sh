#!/usr/bin/env bash

# Examples:
#     pst hello.txt world.py     # paste files.
#     echo Hello world. | pst    # read from STDIN.
#     pst                        # Paste in terminal.

pst() {
    local PST_HOST=https://pst.zidhuss.tech
    [ -t 0 ] && {

        [ $# -gt 0 ] && {
            for filename in "$@"
            do
                if [ -f "$filename" ]
                then
                    curl -F f:1=@"$filename" $PST_HOST
                else
                    echo "file '$filename' does not exist!"
                fi
            done
            return
        }

        echo "^C to cancel, ^D to send."
    }
    curl -F f:1='<-' $PST_HOST
}

pst $*
