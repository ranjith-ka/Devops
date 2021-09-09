# Initialize the autocompletion
autoload -Uz compinit && compinit -i

## Functions
function dex-fn {
	docker exec -it $1 /bin/sh
}

function docker-rm {
	docker rmi $(docker images | awk '{print $1":"$2}' | tail +2 | grep -i $1 | xargs)
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
alias dim='docker images'
alias drmi="docker-rm $1"

ls='ls --color=tty'
grep='grep  --color=auto --exclude-dir={.bzr,CVS,.git,.hg,.svn}'

#### Kubenetes Handy commands

function kube_desc_pod() {
    if [ $# -eq 2 ]; then
        kubectl describe po $1 -n $2
    else
        kubectl describe po $1
    fi
}

function kdex-fn() {
    if [ $# -eq 2 ]; then
        kubectl exec -it $1 /bin/sh -n $2
    else
        kubectl exec -it $1 /bin/sh
    fi
}

alias kctl_po_all='kubectl get po --all-namespaces'
alias kctl_po='kubectl get po'
alias kctl_svc_all='kubectl get svc --all-namespaces'
alias kctl_svc='kubectl get svc'
alias kctl_ds_po='kube_desc_pod'
alias kctl_login="kdex-fn"
alias kctl_busybox='kubectl run -i --tty busybox --image=busybox --restart=Never -- /bin/sh; kubectl delete pod busybox'
