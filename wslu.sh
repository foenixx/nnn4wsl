#!/bin/zsh

cmd_cmd_tab="Open cmd.exe in the new tab"
cmd_cmd_pane="Open cmd.exe in the new pane"
cmd_pwsh_tab="Open PowerShell in the new tab"
cmd_pwsh_pane="Open PowerShell in the new pane"
cmd_vs="Open in VSCode"
cmd_copy_clip="Copy to win clipboard"
cmd=$(echo "$cmd_cmd_tab\n$cmd_cmd_pane\n$cmd_pwsh_tab\n$cmd_pwsh_pane\n$cmd_vs\n$cmd_copy_clip" | fzf )
echo $cmd

#printf "%s" "$fullpath"
winpath=$(wslpath -w "$(pwd)" 2> /dev/null || printf "%s" "$1")
printf "%s" "$winpath"


if [[ "$cmd" = "$cmd_cmd_pane" ]]
then
  wslrun --gui wt.exe -w 0 split-pane  --profile {0caa0dad-35be-5f56-a8ff-afceeeaa6101} -d "$winpath"
  exit
fi

if [[ "$cmd" = "$cmd_cmd_tab" ]]
then
  wslrun --gui wt.exe -w 0 new-tab  --profile {0caa0dad-35be-5f56-a8ff-afceeeaa6101} -d "$winpath"
  exit
fi


if [[ "$cmd" = "$cmd_pwsh_pane" ]]
then
  wslrun --gui wt.exe -w 0 split-pane  --profile {574e775e-4f2a-5b96-ac1e-a2962a402336} -d "$winpath"
  exit
fi

if [[ "$cmd" = "$cmd_pwsh_tab" ]]
then
  wslrun --gui wt.exe -w 0 new-tab  --profile {574e775e-4f2a-5b96-ac1e-a2962a402336} -d "$winpath"
  exit
fi
