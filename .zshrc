# Initialize the autocompletion
autoload -Uz compinit && compinit -i

## Functions
function dex-fn {
	docker exec -it $1 /bin/sh
}

dstop() { docker stop $(docker ps -a -q); }
dimage() { docker rmi $(docker images -a -q) }
drun () { docker container run -it $1 /bin/sh }

alias ll='ls -ahl'
alias ..='cd ..'
alias pass_generate="openssl rand -base64 $1"
alias docker-rk-login="dex-fn"
alias docker-rk-stop='dstop'
alias docker-rk-stop-rm='docker stop $(docker ps -a -q) && docker rm $(docker ps -a -q)'
alias dpa="docker ps -a"
alias docker-rm-images='dimage'
alias docker-rk-run='drun'
alias dmi='docker images'
alias drmi="docker rmi $1"

ls='ls --color=tty'
grep='grep  --color=auto --exclude-dir={.bzr,CVS,.git,.hg,.svn}'
