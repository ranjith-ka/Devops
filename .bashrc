eval keychain --eval --agents ssh --inherit any id_rsa
source $HOME/.keychain/$HOSTNAME-sh

## Functions

function dex-fn() {
    docker exec -it $1 /bin/bash
}

dstop() { docker stop $(docker ps -a -q); }

# workaround for version 0.1.8 of gws'
alias gws="PATH=/usr/local/opt/coreutils/libexec/gnubin:/usr/local/opt/gnu-sed/libexec/gnubin:$PATH gws"

alias ll='ls -ahl'
alias ..='cd ..'
alias vd="vagrant destroy"
alias pass_generate="openssl rand -base64 $1"
alias docker-rk-login="dex-fn"
alias docker-rk-stop='dstop'
alias docker-rk-stop-rm='docker stop $(docker ps -a -q) && docker rm $(docker ps -a -q)'
alias dpa="docker ps -a"
alias ds="dirs -l -v"
_ssh() {
    local cur prev opts
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD - 1]}"
    opts=$(grep '^Host' ~/.ssh/config | grep -v '[?*]' | cut -d ' ' -f 2-)

    COMPREPLY=($(compgen -W "$opts" -- ${cur}))
    return 0
}
complete -F _ssh ssh
complete -F _ssh ssh-copy-id
#complete -F _ssh scp
