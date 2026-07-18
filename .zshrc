# Enable Powerlevel10k instant prompt. Should stay close to the top of ~/.zshrc.
if [[ -r "${XDG_CACHE_HOME:-$HOME/.cache}/p10k-instant-prompt-${(%):-%n}.zsh" ]]; then
  source "${XDG_CACHE_HOME:-$HOME/.cache}/p10k-instant-prompt-${(%):-%n}.zsh"
fi

# Path to Oh My Zsh installation.
export ZSH="$HOME/.oh-my-zsh"

# Theme and plugins.
ZSH_THEME="robbyrussell"
plugins=(git)

# Path setup.
path+=(/opt/homebrew/bin "$HOME/bin")
export PATH

# Load Oh My Zsh.
if [[ -f "$ZSH/oh-my-zsh.sh" ]]; then
  source "$ZSH/oh-my-zsh.sh"
fi

PROMPT='${ret_status} %{$fg[cyan]%}%~%{$reset_color%} $(git_prompt_info)'

# History settings.
HISTFILE="$HOME/.zhistory"
HISTSIZE=1000
SAVEHIST=1000
setopt INC_APPEND_HISTORY
setopt SHARE_HISTORY
setopt HIST_IGNORE_DUPS
setopt HIST_IGNORE_ALL_DUPS
setopt HIST_VERIFY
setopt HIST_EXPIRE_DUPS_FIRST

# Completion and syntax highlighting.
autoload -Uz compinit
compinit -i
source /opt/homebrew/share/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh 2>/dev/null || true

# kubectl completion and kubecolor support if available.
if command -v kubectl >/dev/null 2>&1; then
  source <(kubectl completion zsh) 2>/dev/null || true
  if command -v kubecolor >/dev/null 2>&1; then
    alias kubectl=kubecolor
    compdef kubecolor=kubectl
  fi
fi

# Aliases.
alias ll='ls -ahl'
alias ..='cd ..'
alias dpa='docker ps -a'
alias dim='docker images'
alias k='kubectl'
alias kpo='kubectl get po'
alias kpoall='kubectl get po --all-namespaces'
alias ksvc='kubectl get svc'
alias ksvcall='kubectl get svc --all-namespaces'
alias kx='kubectx'
alias kxn='kubens'

if [[ "$OSTYPE" == darwin* ]]; then
  alias ls='ls -G'
else
  alias ls='ls --color=auto'
fi

alias grep='grep --color=auto --exclude-dir={.bzr,CVS,.git,.hg,.svn}'
alias pass_generate='openssl rand -base64'

# Docker helpers.
docker_exec_shell() {
  if [[ -z "$1" ]]; then
    echo "Usage: docker_exec_shell <container>"
    return 1
  fi
  docker exec -it "$1" /bin/sh
}

docker_remove_images() {
  if [[ -z "$1" ]]; then
    echo "Usage: docker_remove_images <pattern>"
    return 1
  fi
  docker images --format '{{.Repository}}:{{.Tag}}' | grep -i "$1" | xargs -r docker rmi
}

docker_stop_all() {
  docker ps -a -q | xargs -r docker stop
}

docker_rmi_all() {
  docker images -a -q | xargs -r docker rmi
}

docker_run_shell() {
  if [[ -z "$1" ]]; then
    echo "Usage: docker_run_shell <image>"
    return 1
  fi
  docker container run -it "$1" /bin/sh
}

alias docker-rk-login='docker_exec_shell'
alias docker-rk-stop='docker_stop_all'
alias docker-rk-stop-rm='docker stop $(docker ps -a -q) && docker rm $(docker ps -a -q)'
alias docker-rk-run='docker_run_shell'
alias docker-rm-images='docker_rmi_all'
alias docker-rm='docker_remove_images'
alias drmi='docker_remove_images'

# Kubernetes handy commands.
kube_desc_pod() {
  local pod="$1"
  local ns="$2"

  if [[ -z "$pod" ]]; then
    echo "Usage: kube_desc_pod <pod> [namespace]"
    return 1
  fi

  kubectl describe pod "$pod"${ns:+ -n "$ns"}
}

kdex_fn() {
  local pod="$1"
  local ns="$2"

  if [[ -z "$pod" ]]; then
    echo "Usage: kdex_fn <pod> [namespace]"
    return 1
  fi

  kubectl exec -it "$pod" /bin/sh${ns:+ -n "$ns"}
}

alias k_ds_po='kube_desc_pod'
alias k_login='kdex_fn'
alias k_busybox='kubectl run -i --tty busybox --image=busybox --restart=Never -- /bin/sh; kubectl delete pod busybox'

# Powerlevel10k prompt.
if [[ -f /opt/homebrew/share/powerlevel10k/powerlevel10k.zsh-theme ]]; then
  source /opt/homebrew/share/powerlevel10k/powerlevel10k.zsh-theme
fi
[[ -f ~/.p10k.zsh ]] && source ~/.p10k.zsh

# VS Code shell integration.
if [[ "$TERM_PROGRAM" == "vscode" ]] && command -v code >/dev/null 2>&1; then
  . "$(code --locate-shell-integration-path zsh)"
fi
