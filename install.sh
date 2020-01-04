#!/usr/bin/env sh

set -e

is_command() {
    command -v "$1" > /dev/null
}

echoerr() {
    >&2 echo "$@"
}

check_prereqs() {
    if ! is_command curl; then
        echoerr "curl is needed for this script to work"
        exit 1
    fi
    if ! is_command grep; then
        echoerr "grep is needed for this script to work"
        exit 1
    fi
    if ! is_command cut; then
        echoerr "cut is needed for this script to work"
        exit 1
    fi
    if ! is_command sed; then
        echoerr "sed is needed for this script to work"
        exit 1
    fi
}

do_install() {
    target_dir="/usr/local/bin"
    while getopts "b:" arg; do
        case "$arg" in
          b)
              target_dir="$OPTARG"
              ;;
          *)
              exit 1
              ;;
        esac
    done
    shift $((OPTIND - 1))

    version=${1:-}
    [ -z "$version" ] && version="latest"

    set +e
    release_json=$(curl -fsLH 'Accept: application/json' https://github.com/makkes/gitlab-cli/releases/${version})
    set -e
    [ -z "$release_json" ] && echoerr "Unknown version ${version}" && exit 1
    release_tag=$(echo "$release_json"|grep -o '"tag_name":"[^"]*"'|cut -d":" -f2|sed s/\"//g)

    kernel_name=$(uname -s)
    machine=$(uname -m)
    case "${kernel_name}" in
        Darwin)
            os="darwin"
            ;;
        Linux)
            os="linux"
            ;;
        *)
            echoerr "Unsupported OS ${kernel_name}" && exit 1
            ;;
    esac
    [ "${machine}" != "x86_64" ] && echoerr "Unsupported CPU architecture ${machine}" && exit 1
    arch="amd64"

    download_url="https://github.com/makkes/gitlab-cli/releases/download/${release_tag}/gitlab_${release_tag}_${os}_${arch}"
    tmpdir=$(mktemp -d)
    echo "Downloading gitlab ${release_tag}..."
    curl --progress-bar -Lo "${tmpdir}"/gitlab "${download_url}"
    install -t "${target_dir}" "${tmpdir}/gitlab"

    echo "Installed gitlab ${release_tag} into ${target_dir}"
}

check_prereqs
do_install $@
