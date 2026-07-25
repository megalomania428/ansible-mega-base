#!/usr/bin/env bash
set -ueo pipefail
umask 0022
MY_BIN="$(realpath "$0")"
MY_PATH="$(dirname "${MY_BIN}")"
cd "${MY_PATH}/../.."
# shellcheck disable=1091
source "${MY_PATH}/../prepare.sh"
sce='default'
LOG_PATH="/tmp/molecule-$(/usr/bin/env date '+%Y%m%d%H%M%S.%3N')"
printf "\n\n\nmolecule [create] action\n"
# shellcheck disable=2154
ANSIBLE_LOG_PATH="${LOG_PATH}-0create" "${_appimage_bin}" molecule -v create -s "${sce}"
n=1
run_group() {
  local tag="${1}"
  local last="${tag:-}"
  local prefix="${last##*-}"
  [[ -n "${prefix}" ]] && prefix="-${prefix}"
  if [[ -z "${tag:-}" ]]; then
    args=" -- --check"
  else
    args=" -- -t ${tag} --check"
  fi
  printf "\n\n\nmolecule [converge] %s check\n" "${tag:-empty}"
  # shellcheck disable=2086
  ANSIBLE_LOG_PATH="${LOG_PATH}-$(printf %02d $n)converge${prefix}-check" \
    "${_appimage_bin}" molecule -v converge -s "${sce}"${args}
  ((n++))
  for mode in action check; do
    args=''
    if [[ "${mode}" == 'check' ]]; then
      if [[ -z "${tag:-}" ]]; then
        args=" -- --check"
      else
        args=" -- -t ${tag} --check"
      fi
    elif [[ -n "${tag:-}" ]]; then
      args=" -- -t $tag"
    fi
    for stage in converge idempotence; do
      printf "\n\n\nmolecule [%s] %s %s\n" "${stage}" "${mode}" "${tag}"
      # shellcheck disable=2086
      ANSIBLE_LOG_PATH="${LOG_PATH}-$(printf %02d $n)${stage}${prefix}-${mode}" \
        "${_appimage_bin}" molecule -v "${stage}" -s "${sce}"${args}
      ((n++))
    done
  done
}
run_group ''
