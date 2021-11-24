#!/bin/zsh

cmd_locate='locate <text> (find keyword everywhere; "updatedb" may be in order)'
cmd_whereis='whereis <text> (locate binary and manpage in standard Linux places + $PATH)'
cmd_which='which <text> (locate binary in $PATH)'
cmd_find='find . -name <text> (find pattern in the current folder and below)'
cmd=$(echo "$cmd_locate\n$cmd_whereis\n$cmd_which\n$cmd_find" | fzf --print-query --phony )
#echo $cmd

# variable expansion magic from https://unix.stackexchange.com/questions/29724/how-to-properly-collect-an-array-of-lines-in-zsh
arr=("${(@f)cmd}")

if [ -z "$arr[1]" ]
then
  echo "empty text"
  exit
fi

if [[ "$arr[2]" = "$cmd_locate" ]]
then
  echo "issuing 'locate $arr[1]'..."
  locate "$arr[1]"
  exit
fi

if [[ "$arr[2]" = "$cmd_whereis" ]]
then
  echo "issuing 'whereis $arr[1]'..."
  whereis "$arr[1]"
  exit
fi

if [[ "$arr[2]" = "$cmd_which" ]]
then
  echo "issuing 'which $arr[1]'..."
  which "$arr[1]"
  exit
fi

if [[ "$arr[2]" = "$cmd_find" ]]
then
  echo "issuing 'find . -name $arr[1]'..."
  find . -name "$arr[1]"
  exit
fi