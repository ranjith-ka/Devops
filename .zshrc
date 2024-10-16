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

alias k='kubectl'
alias k_po_all='kubectl get po --all-namespaces'
alias k_po='kubectl get po'
alias kpo='kubectl get po'
alias k_svc_all='kubectl get svc --all-namespaces'
alias k_svc='kubectl get svc'
alias k_ds_po='kube_desc_pod'
alias k_login="kdex-fn"
alias k_busybox='kubectl run -i --tty busybox --image=busybox --restart=Never -- /bin/sh; kubectl delete pod busybox'

dive_image() {
  local IMAGE_NAME="${1}"
  local TMP_FILE=/tmp/dive-tmp-image.tar
  docker save "$IMAGE_NAME" > $TMP_FILE && dive $TMP_FILE --source=docker-archive && rm $TMP_FILE
}

# Function to remove all "exited" containers
remove_exited_containers() {
  local exited_containers=$(docker ps -a -f status=exited -q)
  if [ -n "$exited_containers" ]; then
    docker rm $exited_containers
  fi
}

# Function to remove all "created" containers
remove_created_containers() {
  local created_containers=$(docker ps -a -f status=created -q)
  if [ -n "$created_containers" ]; then
    docker rm $created_containers
  fi
}

# Function to remove all stopped containers (both "exited" and "created")
remove_stopped_containers() {
  docker container prune -f
}

# Function to remove all untagged images
remove_untagged_images() {
  local untagged_images=$(docker images -f "dangling=true" -q)
  if [ -n "$untagged_images" ]; then
    docker rmi $untagged_images
  fi
}

cleanup_docker() {
  remove_exited_containers
  remove_created_containers
  remove_stopped_containers
  remove_untagged_images
}