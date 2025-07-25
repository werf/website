#!/usr/bin/env bash
# This script contains functions with different licenses. Some functions specifically annotated  with the author/license information.
# Functions that are not annotated with author/license information released under Apache License 2.0: https://www.apache.org/licenses/LICENSE-2.0.txt
set -uo pipefail

main() {
  REQUIRED_BASH_VERSION="3.2.57"
  REQUIRED_GIT_VERSION="2.18.0"

  validate_bash_version "$REQUIRED_BASH_VERSION"

  OPT_DEFAULT_CI="no"
  OPT_DEFAULT_INTERACTIVE="auto"
  OPT_DEFAULT_JOIN_DOCKER_GROUP="auto"
  OPT_DEFAULT_SETUP_BIN_PATH="auto"
  OPT_DEFAULT_AUTOACTIVATE_WERF="auto"
  OPT_DEFAULT_INSTALL_WERF_SYSTEM_DEPENDENCIES="auto"
  OPT_DEFAULT_SETUP_BUILDAH="auto"
  OPT_DEFAULT_VERIFY_TRDL_SIGNATURE="yes"
  OPT_DEFAULT_SHELL="auto"
  OPT_DEFAULT_WERF_AUTOACTIVATE_VERSION="2"
  OPT_DEFAULT_WERF_AUTOACTIVATE_CHANNEL="stable"
  OPT_DEFAULT_INITIAL_TRDL_VERSION="0.7.0"
  OPT_DEFAULT_TRDL_TUF_REPO="https://tuf.trdl.dev"
  OPT_DEFAULT_TRDL_PGP_PUBKEY_URL="https://trdl.dev/trdl-client.asc"
  OPT_DEFAULT_WERF_TUF_REPO="https://tuf.werf.io"
  OPT_DEFAULT_WERF_TUF_ROOT_VERSION="1"
  OPT_DEFAULT_WERF_TUF_ROOT_SHA="b7ff6bcbe598e072a86d595a3621924c8612c7e6dc6a82e919abe89707d7e3f468e616b5635630680dd1e98fc362ae5051728406700e6274c5ed1ad92bea52a2"
  OPT_REGEX_SHELL='^(bash|zsh|auto)?$'
  OPT_REGEX_WERF_AUTOACTIVATE_VERSION='^[.0-9]+$'
  OPT_REGEX_WERF_AUTOACTIVATE_CHANNEL='.+'
  OPT_REGEX_INITIAL_TRDL_VERSION='^[0-9]+.[0-9]+.[0-9]+$'
  OPT_REGEX_TRDL_TUF_REPO='^https?://.+'
  OPT_REGEX_TRDL_PGP_PUBKEY_URL='^https?://.+'
  OPT_REGEX_WERF_TUF_REPO='^https?://.+'
  OPT_REGEX_WERF_TUF_ROOT_VERSION='^[0-9]+$'
  OPT_REGEX_WERF_TUF_ROOT_SHA='^[a-fA-F0-9]+$'
  OPT_REGEX_USER='^[a-z_][a-z0-9_-]{0,31}$'
  eval "$(getoptions getoptions_config) abort 'Failure parsing options.'"

  OPT_CI="${OPT_CI:-$OPT_DEFAULT_CI}"
  if [ "$OPT_CI" == "yes" ]; then
    OPT_INTERACTIVE="${OPT_INTERACTIVE:-"no"}"
    OPT_JOIN_DOCKER_GROUP="${OPT_JOIN_DOCKER_GROUP:-"no"}"
    OPT_SETUP_BIN_PATH="${OPT_SETUP_BIN_PATH:-"no"}"
    OPT_AUTOACTIVATE_WERF="${OPT_AUTOACTIVATE_WERF:-"no"}"
    OPT_INSTALL_WERF_SYSTEM_DEPENDENCIES="${OPT_INSTALL_WERF_SYSTEM_DEPENDENCIES:-"no"}"
    OPT_SETUP_BUILDAH="${OPT_SETUP_BUILDAH:-"no"}"
  else
    OPT_INTERACTIVE="${OPT_INTERACTIVE:-$OPT_DEFAULT_INTERACTIVE}"
    OPT_JOIN_DOCKER_GROUP="${OPT_JOIN_DOCKER_GROUP:-$OPT_DEFAULT_JOIN_DOCKER_GROUP}"
    OPT_SETUP_BIN_PATH="${OPT_SETUP_BIN_PATH:-$OPT_DEFAULT_SETUP_BIN_PATH}"
    OPT_AUTOACTIVATE_WERF="${OPT_AUTOACTIVATE_WERF:-$OPT_DEFAULT_AUTOACTIVATE_WERF}"
    OPT_INSTALL_WERF_SYSTEM_DEPENDENCIES="${OPT_INSTALL_WERF_SYSTEM_DEPENDENCIES:-$OPT_DEFAULT_INSTALL_WERF_SYSTEM_DEPENDENCIES}"
    OPT_SETUP_BUILDAH="${OPT_SETUP_BUILDAH:-$OPT_DEFAULT_SETUP_BUILDAH}"
  fi
  OPT_VERIFY_TRDL_SIGNATURE="${OPT_VERIFY_TRDL_SIGNATURE:-$OPT_DEFAULT_VERIFY_TRDL_SIGNATURE}"
  OPT_SHELL="${OPT_SHELL:-$OPT_DEFAULT_SHELL}"
  OPT_WERF_AUTOACTIVATE_VERSION="${OPT_WERF_AUTOACTIVATE_VERSION:-$OPT_DEFAULT_WERF_AUTOACTIVATE_VERSION}"
  OPT_WERF_AUTOACTIVATE_CHANNEL="${OPT_WERF_AUTOACTIVATE_CHANNEL:-$OPT_DEFAULT_WERF_AUTOACTIVATE_CHANNEL}"
  OPT_INITIAL_TRDL_VERSION="${OPT_INITIAL_TRDL_VERSION:-$OPT_DEFAULT_INITIAL_TRDL_VERSION}"
  OPT_TRDL_TUF_REPO="${OPT_TRDL_TUF_REPO:-$OPT_DEFAULT_TRDL_TUF_REPO}"
  OPT_TRDL_PGP_PUBKEY_URL="${OPT_TRDL_PGP_PUBKEY_URL:-$OPT_DEFAULT_TRDL_PGP_PUBKEY_URL}"
  OPT_WERF_TUF_REPO="${OPT_WERF_TUF_REPO:-$OPT_DEFAULT_WERF_TUF_REPO}"
  OPT_WERF_TUF_ROOT_VERSION="${OPT_WERF_TUF_ROOT_VERSION:-$OPT_DEFAULT_WERF_TUF_ROOT_VERSION}"
  OPT_WERF_TUF_ROOT_SHA="${OPT_WERF_TUF_ROOT_SHA:-$OPT_DEFAULT_WERF_TUF_ROOT_SHA}"
  OPT_USER="${OPT_USER:-$(get_user)}"
  
  export HOME="$(eval echo ~$OPT_USER)"

  check_user_existence "$OPT_USER"
  install_werf_system_dependencies "$OPT_INSTALL_WERF_SYSTEM_DEPENDENCIES"
  ensure_cmds_available uname git grep tee install
  validate_git_version "$REQUIRED_GIT_VERSION"

  setup_buildah "$OPT_SETUP_BUILDAH" "$OPT_USER"

  [[ "$OPT_VERIFY_TRDL_SIGNATURE" == "no" ]] && ensure_cmds_available gpg

  declare arch
  arch="$(get_arch)" || abort "Failure getting system architecture."

  declare os
  os="$(get_os)" || abort "Failure getting OS."

  declare shell
  if [[ "$OPT_SHELL" == "auto" ]]; then
    set_shell || abort "Failure getting user shell."
  else
    shell="$OPT_SHELL"
  fi

  [[ $os == "linux" ]] && propose_joining_docker_group "$OPT_JOIN_DOCKER_GROUP" "$OPT_USER"
  setup_trdl_bin_path "$shell" "$OPT_SETUP_BIN_PATH"
  install_trdl "$os" "$arch" "$OPT_TRDL_TUF_REPO" "$OPT_TRDL_PGP_PUBKEY_URL" "$OPT_INITIAL_TRDL_VERSION" "$OPT_VERIFY_TRDL_SIGNATURE" "$OPT_USER"
  add_trdl_werf_repo "$OPT_WERF_TUF_REPO" "$OPT_WERF_TUF_ROOT_VERSION" "$OPT_WERF_TUF_ROOT_SHA" "$OPT_USER" 
  enable_automatic_werf_activation "$shell" "$OPT_AUTOACTIVATE_WERF" "$OPT_WERF_AUTOACTIVATE_VERSION" "$OPT_WERF_AUTOACTIVATE_CHANNEL"
  finalize "$OPT_WERF_AUTOACTIVATE_VERSION" "$OPT_WERF_AUTOACTIVATE_CHANNEL"
}

getoptions_config() {
  setup   REST on:"yes" no:"no" init:"" help:usage -- "Usage: $(basename "$0") [options]" ''
  msg -- "Options (use + instead of - for short options to disable them):"

  flag OPT_CI --{no-}ci -- \
    "CI mode. Configures a few other options in a way recommended for CI (i.e. enables non-interactive mode). Full list of configured options can be found in the installer script itself. Default: $OPT_DEFAULT_CI"

  flag OPT_INTERACTIVE -i --{no-}interactive -- \
    "If run non-interactively, then choose recommended answers for all prompts. Default: $OPT_DEFAULT_INTERACTIVE"
  flag OPT_VERIFY_TRDL_SIGNATURE --{no-}verify-trdl-signature -- \
    "Verify PGP signature of trdl binary. Default: $OPT_DEFAULT_VERIFY_TRDL_SIGNATURE"
  flag OPT_JOIN_DOCKER_GROUP --{no-}join-docker-group -- \
    "Add user to the docker group? Default: $OPT_DEFAULT_JOIN_DOCKER_GROUP"
  flag OPT_SETUP_BIN_PATH -p --{no-}setup-bin-path -- \
    "Add \"$HOME/bin\" to \$PATH for trdl? Default: $OPT_DEFAULT_SETUP_BIN_PATH"
  flag OPT_AUTOACTIVATE_WERF -a --{no-}autoactivate-werf -- \
    "Enable werf autoactivation? Default: $OPT_DEFAULT_AUTOACTIVATE_WERF"
  flag OPT_INSTALL_WERF_SYSTEM_DEPENDENCIES --{no-}install-werf-system-depedencies -- \
    "Install system dependencies for werf? Default: $OPT_DEFAULT_INSTALL_WERF_SYSTEM_DEPENDENCIES"
  flag OPT_SETUP_BUILDAH --{no-}setup-buildah -- \
    "Install and set up Buildah backend? Default: $OPT_DEFAULT_SETUP_BUILDAH"
  
  param OPT_SHELL -s --shell validate:"validate_option_by_regex '$OPT_REGEX_SHELL'" -- \
    "Shell, for which werf/trdl should be set up. Default: $OPT_DEFAULT_SHELL. Allowed values regex: $OPT_REGEX_SHELL"
  param OPT_WERF_AUTOACTIVATE_VERSION -v --version validate:"validate_option_by_regex '$OPT_REGEX_WERF_AUTOACTIVATE_VERSION'" -- \
    "Autoactivated (if enabled) werf version. Default: $OPT_DEFAULT_WERF_AUTOACTIVATE_VERSION. Allowed values regex: $OPT_REGEX_WERF_AUTOACTIVATE_VERSION"
  param OPT_WERF_AUTOACTIVATE_CHANNEL -c --channel validate:"validate_option_by_regex '$OPT_REGEX_WERF_AUTOACTIVATE_CHANNEL'" -- \
    "Autoactivated (if enabled) werf channel. Default: $OPT_DEFAULT_WERF_AUTOACTIVATE_CHANNEL. Allowed values regex: $OPT_REGEX_WERF_AUTOACTIVATE_CHANNEL"
  param OPT_INITIAL_TRDL_VERSION --initial-trdl-version validate:"validate_option_by_regex '$OPT_REGEX_INITIAL_TRDL_VERSION'" -- \
    "Initially installed version of trdl (self-updates automatically). Default: $OPT_DEFAULT_INITIAL_TRDL_VERSION. Allowed values regex: $OPT_REGEX_INITIAL_TRDL_VERSION"
  param OPT_TRDL_TUF_REPO --trdl-tuf-repo validate:"validate_option_by_regex '$OPT_REGEX_TRDL_TUF_REPO'" -- \
    "trdl TUF-repository URL. Default: $OPT_DEFAULT_TRDL_TUF_REPO. Allowed values regex: $OPT_REGEX_TRDL_TUF_REPO"
  param OPT_TRDL_PGP_PUBKEY_URL --trdl-pgp-pubkey-url validate:"validate_option_by_regex '$OPT_REGEX_TRDL_PGP_PUBKEY_URL'" -- \
    "trdl GPG public key URL. Default: $OPT_DEFAULT_TRDL_PGP_PUBKEY_URL. Allowed values regex: $OPT_REGEX_TRDL_PGP_PUBKEY_URL"
  param OPT_WERF_TUF_REPO --werf-tuf-repo validate:"validate_option_by_regex '$OPT_REGEX_WERF_TUF_REPO'" -- \
    "werf TUF-repository URL. Default: $OPT_DEFAULT_WERF_TUF_REPO. Allowed values regex: $OPT_REGEX_WERF_TUF_REPO"
  param OPT_WERF_TUF_ROOT_VERSION --werf-tuf-root-version validate:"validate_option_by_regex '$OPT_REGEX_WERF_TUF_ROOT_VERSION'" -- \
    "werf TUF-repository root version. Default: $OPT_DEFAULT_WERF_TUF_ROOT_VERSION. Allowed values regex: $OPT_REGEX_WERF_TUF_ROOT_VERSION"
  param OPT_WERF_TUF_ROOT_SHA --werf-tuf-root-sha validate:"validate_option_by_regex '$OPT_REGEX_WERF_TUF_ROOT_SHA'" -- \
    "werf TUF-repository root SHA-hash. Default: $OPT_DEFAULT_WERF_TUF_ROOT_SHA. Allowed values regex: $OPT_REGEX_WERF_TUF_ROOT_SHA"
  param OPT_USER -u --user validate:"validate_option_by_regex '$OPT_REGEX_USER'" -- \
    "User for which subuid subguid should be set up. Allowed values regex: $OPT_REGEX_USER"

  disp    :usage  -h --help
}

validate_option_by_regex() {
  declare value="$OPTARG"
  declare valid_regex="$1"

  [[ $value =~ $valid_regex ]] || return 1
}

validate_bash_version() {
  declare required_bash_version="$1"

  declare current_bash_version="${BASH_VERSINFO[0]}.${BASH_VERSINFO[1]}.${BASH_VERSINFO[2]}"
  compare_versions "$current_bash_version" "$required_bash_version"
  [[ $? -gt 1 ]] && abort "Bash version must be at least \"$required_bash_version\"."
}

validate_git_version() {
  declare required_git_version="$1"

  declare current_git_version
  current_git_version="$(git --version | awk '{print $3}' | awk -F. '{printf("%s.%s.%s", $1, $2, $3)}')" || abort "Unable to parse git version."
  compare_versions "$current_git_version" "$required_git_version"
  case $? in
    1)  log::info "Current version is greater than or equal to required version." ;;
    2)  log::info "Git version is too old, proceeding to upgrade." 
        if is_command_exists apt-get; then
          log::info "Debian-based system detected. Adding PPA and installing Git..."
          if ! is_command_exists add-apt-repository; then
            log::info "Installing add-apt-repository..."
            install_package software-properties-common
          fi
          log::info "Adding PPA for Git..."
          run_as_root "add-apt-repository ppa:git-core/ppa"
          install_package git
        fi
        ;;
    0)  log::info "Versions are equal." ;;
    *)  abort "Unexpected error comparing versions." ;;
  esac
}

get_arch() {
  declare arch
  arch="$(uname -m)" || abort "Can't get system architecture."

  case "$arch" in
    x86_64) arch="amd64" ;;
    arm64 | armv8* | aarch64*) arch="arm64" ;;
    i386 | i486 | i586 | i686) abort "werf is not available for x86 architecture." ;;
    arm | armv7*) abort "werf is not available for 32-bit ARM architectures." ;;
    *) abort "werf is not available for \"$arch\" architecture." ;;
  esac

  printf '%s' "$arch"
}

get_os() {
  declare os
  os="$(uname | tr '[:upper:]' '[:lower:]')" || abort "Can't detect OS."

  case "$os" in
    linux*) os="linux" ;;
    darwin*) os="darwin" ;;
    msys* | cygwin* | mingw* | nt | win*) abort "This installer does not support Windows." ;;
    *) abort "Unknown OS." ;;
  esac

  printf '%s' "$os"
}

set_shell() {
  shell="$(printf '%s' "$SHELL" | rev | cut -d'/' -f1 | rev)" || abort "Can't determine login shell."

  declare default_shell="$shell"
  declare answer
  while :; do
    case "$shell" in
      bash | zsh)
        printf '[INPUT REQUIRED] Current login shell is "%s". Press ENTER to setup werf for this shell or choose another one.\n[b]ash/[z]sh/[a]bort? Default: %s.\n' "$shell" "$default_shell"
        ;;
      *)
        default_shell="bash"
        printf '[INPUT REQUIRED] Current login shell is "%s". This shell is unsupported. Choose another shell or abort installation.\n[b]ash/[z]sh/[a]bort? Default: %s.\n' "$shell" "$default_shell"
        ;;
    esac

    if [[ $OPT_INTERACTIVE == "no" ]]; then
      answer=""
    else
      read -r answer
    fi

    case "${answer:-$default_shell}" in
      [bB] | [bB][aA][sS][hH])
        shell="bash"
        break
        ;;
      [zZ] | [zZ][sS][hH])
        shell="zsh"
        break
        ;;
      [aA] | [aA][bB][oO][rR][tT])
        printf 'Aborted by the user.\n' 1>&2
        exit 0
        ;;
      *)
        printf 'Invalid choice, please retry.\n' 1>&2
        continue
        ;;
    esac
  done
}

propose_joining_docker_group() {
  declare override_join_docker_group="$1"
  local user="$2"

  [[ $override_join_docker_group == "no" ]] && return 0
  if is_command_exists docker; then
    is_user_in_group "$user" docker && return 0
    is_root && return 0
    ensure_cmds_available usermod
    if [[ $override_join_docker_group == "auto" ]]; then
      prompt_yes_no_skip 'werf needs access to the Docker daemon. Add current user to the "docker" group? (root required)' "yes" || return 0
    fi
    run_as_root "usermod -aG docker \"$user\"" || abort "Can't add user \"$user\" to group \"docker\"."
  else
    return 0
  fi    
}

setup_trdl_bin_path() {
  declare shell="$1"
  declare override_setup_trdl_bin_path="$2"

  [[ $override_setup_trdl_bin_path == "no" ]] && return 0

  case "$shell" in
    bash)
      declare active_bash_profile_file="$(get_active_bash_profile_file)"
      if [[ $override_setup_trdl_bin_path == "auto" ]]; then
        prompt_yes_no_skip "trdl is going to be installed in \"$HOME/bin/\". Add this directory to your \$PATH in \"$HOME/.bashrc\" and \"$active_bash_profile_file\"? (strongly recommended)" "yes" || return 0
      fi
      add_trdl_bin_path_setup_to_file "$HOME/.bashrc"
      add_trdl_bin_path_setup_to_file "$active_bash_profile_file"
      ;;
    zsh)
      declare active_zsh_profile_file="$(get_active_zsh_profile_file)"
      if [[ $override_setup_trdl_bin_path == "auto" ]]; then
        prompt_yes_no_skip "trdl is going to be installed in \"$HOME/bin/\". Add this directory to your \$PATH in \"$HOME/.zshrc\" and \"$active_zsh_profile_file\"? (strongly recommended)" "yes" || return 0
      fi
      add_trdl_bin_path_setup_to_file "$HOME/.zshrc"
      add_trdl_bin_path_setup_to_file "$active_zsh_profile_file"
      ;;
    *) abort "Shell \"$shell\" is not supported." ;;
  esac
}

add_trdl_bin_path_setup_to_file() {
  declare file="$1"

  declare path_append_cmd='[[ "$PATH" == *"$HOME/bin:"* ]] || export PATH="$HOME/bin:$PATH"'
  grep -qsxF -- "$path_append_cmd" "$file" && log::info "Skipping adding \"$HOME/bin/\" to \$PATH in \"$file\": already added." && return 0

  append_line_to_file "$path_append_cmd" "$file"
}

install_trdl() {
  declare os="$1"
  declare arch="$2"
  declare trdl_tuf_repo="$3"
  declare trdl_gpg_pubkey_url="$4"
  declare trdl_version="$5"
  declare skip_signature_verification="$6"
  declare user="$7"

  is_trdl_installed_and_up_to_date "$trdl_version" && log::info "Skipping trdl installation: already installed in \"$HOME/bin/\"." && return 0
  log::info "Installing trdl to \"$HOME/bin/\"."

  run_as_user "$user" "mkdir -p \"$HOME/bin\"" || abort "Can't create \"$HOME/bin\" directory."

  declare tmp_dir
  tmp_dir="$(run_as_user "$user" "mktemp -d")" || abort "Can't create temporary directory."

  run_as_user "$user" "$(download ${trdl_tuf_repo}/targets/releases/${trdl_version}/${os}-${arch}/bin/trdl $tmp_dir/trdl)"
  run_as_user "$user" "$(download ${trdl_tuf_repo}/targets/signatures/${trdl_version}/${os}-${arch}/bin/trdl.sig $tmp_dir/trdl.sig)"


  if [[ $skip_signature_verification == "no" ]]; then
    log::info "Verifying trdl binary signatures."
    run_as_user "$user" "$(download $trdl_gpg_pubkey_url $tmp_dir/trdl.pub)"
    run_as_user "$user" "gpg -q --import <$tmp_dir/trdl.pub" || abort "Can't import trdl public gpg key."
    run_as_user "$user" "gpg -q --verify $tmp_dir/trdl.sig $tmp_dir/trdl" || abort "Can't verify trdl binary with a gpg signature."
  fi

  run_as_user "$user" "install $tmp_dir/trdl \"$HOME/bin\"" || abort "Can't install trdl binary from \"$tmp_dir/trdl\" to \"$HOME/bin\"."
  run_as_user "$user" "rm -rf $tmp_dir" 2>/dev/null
}

is_trdl_installed_and_up_to_date() {
  declare required_trdl_version="$1"

  [[ -x "$HOME/bin/trdl" ]] || return 1

  declare current_trdl_version
  current_trdl_version="$("$HOME/bin/trdl" version 2>/dev/null | cut -c2-)" || return 1

  compare_versions "$current_trdl_version" "$required_trdl_version"
  test $? -lt 2
  return
}

add_trdl_werf_repo() {
  declare werf_tuf_repo_url="$1"
  declare werf_tuf_repo_root_version="$2"
  declare werf_tuf_repo_root_sha="$3"
  declare user="$4"

  log::info "Adding werf repo to trdl."
  run_as_user "$user" "\"$HOME/bin/trdl\" add werf $werf_tuf_repo_url $werf_tuf_repo_root_version $werf_tuf_repo_root_sha" || abort "Can't add \"werf\" repo to trdl."
}

enable_automatic_werf_activation() {
  declare shell="$1"
  declare override_werf_enable_autoactivation="$2"
  declare werf_autoactivate_version="$3"
  declare werf_autoactivate_channel="$4"

  [[ $override_werf_enable_autoactivation == "no" ]] && return 0

  case "$shell" in
    bash)
      declare active_bash_profile_file="$(get_active_bash_profile_file)"
      if [[ $override_werf_enable_autoactivation == "auto" ]]; then
        prompt_yes_no_skip "Add automatic werf activation to \"$HOME/.bashrc\" and \"$active_bash_profile_file\"? (recommended for interactive usage, not recommended for CI)" "yes" || return 0
      fi
      add_automatic_werf_activation_to_file "$HOME/.bashrc" "$werf_autoactivate_version" "$werf_autoactivate_channel"
      add_automatic_werf_activation_to_file "$active_bash_profile_file" "$werf_autoactivate_version" "$werf_autoactivate_channel"
      ;;
    zsh)
      declare active_zsh_profile_file="$(get_active_zsh_profile_file)"
      if [[ $override_werf_enable_autoactivation == "auto" ]]; then
        prompt_yes_no_skip "Add automatic werf activation to \"$HOME/.zshrc\" and \"$active_zsh_profile_file\"? (recommended for interactive usage, not recommended for CI)" "yes" || return 0
      fi
      add_automatic_werf_activation_to_file "$HOME/.zshrc" "$werf_autoactivate_version" "$werf_autoactivate_channel"
      add_automatic_werf_activation_to_file "$active_zsh_profile_file" "$werf_autoactivate_version" "$werf_autoactivate_channel"
      ;;
    *) abort "Shell \"$shell\" is not supported." ;;
  esac
}

add_automatic_werf_activation_to_file() {
  declare file="$1"
  declare werf_autoactivate_version="$2"
  declare werf_autoactivate_channel="$3"

  declare werf_activation_cmd_grep1="trdl use werf"
  declare werf_activation_cmd_grep2="trdl\" use werf"

  declare werf_activation_cmd="! { which werf | grep -qsE \"^$HOME/.trdl/\"; } && [[ -x \"\$HOME/bin/trdl\" ]] && source \$(\"\$HOME/bin/trdl\" use werf \"$werf_autoactivate_version\" \"$werf_autoactivate_channel\")"

  if grep -qsF -- "$werf_activation_cmd_grep1" "$file" || \
    grep -qsF -- "$werf_activation_cmd_grep2" "$file"; then
    log::info "Skipping adding werf activation to \"$file\": already found \"trdl use werf\" command."
    return 0
  fi

  append_line_to_file "$werf_activation_cmd" "$file"
}

get_active_bash_profile_file() {
  declare result="$HOME/.bash_profile"
  for file in "$HOME/.bash_profile" "$HOME/.bash_login" "$HOME/.profile"; do
    [[ -r $file ]] && result="$file" && break
  done
  printf '%s' "$result"
}

get_active_zsh_profile_file() {
  declare result="$HOME/.zprofile"
  for file in "$HOME/.zprofile" "$HOME/.zlogin"; do
    [[ -r $file ]] && result="$file" && break
  done
  printf '%s' "$result"
}

finalize() {
  declare werf_autoactivate_version="$1"
  declare werf_autoactivate_channel="$2"

  log::info "werf installation finished successfully!"
  log::info "Open new shell session if you have enabled werf autoactivation or activate werf manually with:\n$ source \$(\"$HOME/bin/trdl\" use werf \"$werf_autoactivate_version\" \"$werf_autoactivate_channel\")"
}

append_line_to_file() {
  declare line="$1"
  declare file="$2"

  create_file "$file"
  echo >> "$file"
  declare append_cmd="tee -a '$file' <<< '$line' 1>/dev/null"
  [[ -w $file ]] || append_cmd="run_as_root '$append_cmd'"
  eval $append_cmd || abort "Can't append line \"$line\" to file \"$file\"."
}

create_file() {
  declare file="$1"

  [[ -e $file ]] && return 0

  declare stderr
  if ! stderr="$(touch "$file" 1>/dev/null 2>&1)" && [[ $stderr == *"Permission denied"* ]]; then
    run_as_root "touch '$file'" || abort "Unable to create file \"$file\"."
  fi
}

download() {
  declare url="$1"
  declare path="$2"

  declare download_cmd
  if is_command_exists wget; then
    download_cmd="wget -qO '$path' '$url'"
  elif is_command_exists curl; then
    download_cmd="curl -sSLo '$path' '$url'"
  else
    abort 'Neither "wget" nor "curl" available. Install and retry.'
  fi

  eval $download_cmd || abort "Unable to download \"$url\" to \"$path\"."
}

is_user_in_group() {
  declare user="$1"
  declare group="$2"

  id -nG "$user" 2>/dev/null | grep -wqs "$group"
  return
}

ensure_cmds_available() {
  declare cmd
  for cmd in "${@}"; do
    is_command_exists "$cmd" || abort "\"$cmd\" required but not available."
  done
}

run_as_root() {
  declare cmd="$1"

  if is_root; then
    eval $cmd || return 1
  else
    [[ $OPT_INTERACTIVE == "no" ]] && log::warn "Non-interactive mode enabled, but current user doesn't have enough privilege, so the next command will be called with sudo (this might hang): sudo $cmd"
    ensure_cmds_available sudo
    eval sudo $cmd || return 1
  fi

  return 0
}

get_user() {
  id -un
}

run_as_user() {
  local user="$1"
  shift
  local cmd="$@"
  if [[ "$user" == "$(get_user)" ]]; then
    eval "$cmd" || return 1
  else
    ensure_cmds_available sudo
    [[ $OPT_INTERACTIVE == "no" ]] && log::warn "Non-interactive mode is enabled, the following command will be invoked using sudo for the specified user (this might hang): sudo -u \"$user\" \"$cmd\""
    sudo -u "$user" bash -c "$cmd" || return 1
  fi

} 

is_root() {
  test "$(id -u)" -eq 0
  return
}

is_command_exists() {
  declare command="$1"
  command -v "$command" 1>/dev/null 2>&1
  return
}

# Returns 0 if answer is "yes".
# Returns 1 if answer is "skip".
# Aborts if answer is "abort".
# Returns exit code corresponding to the $def_arg if non-interactive.
prompt_yes_no_skip() {
  declare prompt_msg="$1"
  declare def_arg="$2"

  case "$def_arg" in
    [yY] | [yY][eE][sS])
      def_arg="yes"
      ;;
    [aA] | [aA][bB][oO][rR][tT])
      def_arg="abort"
      ;;
    [sS] | [sS][kK][iI][pP])
      def_arg="skip"
      ;;
  esac

  declare answer
  while :; do
    printf '[INPUT REQUIRED] %s\n[y]es/[a]bort/[s]kip? Default: %s.\n' "$prompt_msg" "$def_arg"

    if [[ $OPT_INTERACTIVE == "no" ]]; then
      answer=""
    else
      read -r answer
    fi

    case "${answer:-$def_arg}" in
      [yY] | [yY][eE][sS])
        return 0
        ;;
      [aA] | [aA][bB][oO][rR][tT])
        printf 'Aborted by the user.\n'
        exit 0
        ;;
      [sS] | [sS][kK][iI][pP])
        return 1
        ;;
      *)
        printf 'Invalid choice, please retry.\n'
        continue
        ;;
    esac
  done
}

# Returns 0 if versions equal.
# Returns 1 if $version1 is greater than $version2.
# Returns 2 if $version1 is less than $version2.
compare_versions() {
  declare version1="$1"
  declare version2="$2"

  for ver in "$version1" "$version2"; do
    [[ $ver =~ ^[.0-9]*$ ]] || abort "Invalid version passed to compare_versions(): $ver"
  done

  [[ $version1 == "$version2" ]] && return 0

  declare IFS=.
  declare -a ver1 ver2
  read -r -a ver1 <<<"$version1"
  read -r -a ver2 <<<"$version2"

  for ((i = ${#ver1[@]}; i < ${#ver2[@]}; i++)); do
    ver1[i]=0
  done

  for ((i = 0; i < ${#ver1[@]}; i++)); do
    if [[ -z ${ver2[i]} ]]; then
      ver2[i]=0
    fi
    if ((10#${ver1[i]} > 10#${ver2[i]})); then
      return 1
    fi
    if ((10#${ver1[i]} < 10#${ver2[i]})); then
      return 2
    fi
  done

  return 0
}

install_werf_system_dependencies() {
  local override_install_werf_system_dependencies="$1"

  [[ $override_install_werf_system_dependencies == "no" ]] && return 0

  if [[ $override_install_werf_system_dependencies == "auto" ]]; then
    prompt_yes_no_skip "Install system dependencies for werf?" "skip" || return 0
  fi

  local missing_packages=""

  is_command_exists git || missing_packages="$missing_packages git"
  is_command_exists curl || missing_packages="$missing_packages curl"
  is_command_exists gpg || missing_packages="$missing_packages gpg"

  if [ "$missing_packages" == "" ]; then
    return 0
  fi

  install_package "$missing_packages"
}

setup_buildah() {

  local override_setup_buildah="$1"
  local user="$2"

  [[ $override_setup_buildah == "no" ]] && return 0

  if [[ $override_setup_buildah == "auto" ]]; then
    prompt_yes_no_skip "Install and set up Buildah backend?" "skip" || return 0
  fi

  is_command_exists buildah || install_buildah "$user"
    
  if is_command_exists sysctl; then
  # Check if kernel.unprivileged_userns_clone exists
    if ! [ -z "$(sysctl -ne kernel.unprivileged_userns_clone)" ]; then
      KERNEL_SETTING="$(sysctl -ne kernel.unprivileged_userns_clone)"
      if [ "$KERNEL_SETTING" = "0" ]; then
      echo 'kernel.unprivileged_userns_clone = 1' | run_as_root "tee -a /etc/sysctl.conf"
      run_as_root "sysctl -w kernel.unprivileged_userns_clone=1"
      log::info "Updated kernel.unprivileged_userns_clone value"
      fi
    fi

    # Check if kernel.apparmor_restrict_unprivileged_userns exists
    if ! [ -z "$(sysctl -ne kernel.apparmor_restrict_unprivileged_userns)" ]; then
      KERNEL_APPARMOR="$(sysctl -ne kernel.apparmor_restrict_unprivileged_userns)"
      if [ "$KERNEL_APPARMOR" = "1" ]; then
        log::info "Set kernel.apparmor_restrict_unprivileged_userns to 0"
        echo "kernel.apparmor_restrict_unprivileged_userns = 0" | run_as_root "tee -a /etc/sysctl.conf" \ &&
          run_as_root "sysctl -w kernel.apparmor_restrict_unprivileged_userns=0"
      fi
    fi

    # Check the sysctl value: user.max_user_namespaces
    USER_NAMESPACES_SETTING="$(sysctl -n user.max_user_namespaces)"
    if [[ "$USER_NAMESPACES_SETTING" < "15000" ]]; then
      echo 'user.max_user_namespaces = 15000' | run_as_root "tee -a /etc/sysctl.conf"
      run_as_root "sysctl -w user.max_user_namespaces=15000"
      log::info "Updated user.max_user_namespaces value"
    fi
  else 
    log::warn "sysctl not found"
  fi
}

get_linux_distro() {
  declare distro
  if [[ -f /etc/os-release ]]; then
    . "/etc/os-release"
    distro="$NAME"
  elif type lsb_release 1>/dev/null 2>&1; then
    distro="$(lsb_release -si)"
  elif [[ -f "/etc/lsb-release" ]]; then
    . "/etc/lsb-release"
    distro="$DISTRIB_ID"
  elif [[ -f "/etc/debian_version" ]]; then
    distro="debian"
  elif [[ -f "/etc/SuSe-release" ]]; then
    distro="suse"
  elif [[ -f "/etc/redhat-release" ]]; then
    distro="redhat"
  else
    abort "Unable to detect Linux distro."
  fi
  echo "$distro"
}

get_distro_version() {
  declare version
  if [[ -f /etc/os-release ]]; then
    . "/etc/os-release"
    version="$VERSION_ID"
  elif type lsb_release 1>/dev/null 2>&1; then
    version=$(lsb_release -r | awk -F"\t" '{print $2}')
  elif [[ -f "/etc/lsb-release" ]]; then
    . "/etc/lsb-release"
    version="$DISTRIB_RELEASE"
  elif [[ -f "/etc/debian_version" ]]; then
    version=$(cat /etc/debian_version)
  elif [[ -f "/etc/SuSe-release" ]]; then
    version=$(cat /etc/SuSe-release | grep -oP 'VERSION = \K[^\s]+')
  elif [[ -f "/etc/redhat-release" ]]; then
    version=$(cat /etc/redhat-release | sed -E 's/[^0-9]*([0-9]+\.[0-9]+).*/\1/')
  else
    echo "Unable to detect Linux distro."
    return 1
  fi
  echo "$version"
}

install_package() {
  local package="$1"
  distro=$(get_linux_distro)
  
  log::info "Detected Linux distribution: $distro"

  case "$distro" in
    *CentOS*|*Red\ Hat*|*Fedora*)
      run_as_root "yum -y install $package"
      ;;
    *Debian*|*Ubuntu*)
      run_as_root "apt-get -y update"
      run_as_root "DEBIAN_FRONTEND=noninteractive apt-get -y install $package"
      ;;
    *SUSE*|*openSUSE*)
      run_as_root "zypper install $package"
      ;;
    *RHEL*)
      if [[ "$distro" == *"7"* ]]; then
        run_as_root "yum -y install $package"
      elif [[ "$distro" == *"8"* ]]; then
        run_as_root "dnf module install -y $package"  
      fi
      ;;
    *)
      abort "Unsupported distribution: $distro"
      ;;
  esac
}

install_buildah(){
  local user=$1
  distro=$(get_linux_distro)
  distro_version=$(get_distro_version)
  log::info "Installing Buildah and configuring user namespaces..."

  case "$distro" in
    *Debian*)
      install_package "buildah uidmap"
      ;;
    *Ubuntu*)
      compare_versions "$distro_version" "20.04"
      case $? in
        0|2)
          echo "deb http://download.opensuse.org/repositories/devel:/kubic:/libcontainers:/stable/xUbuntu_${distro_version}/ /" \
            | run_as_root "tee /etc/apt/sources.list.d/devel-kubic-libcontainers-stable.list"

          curl -Ls "https://download.opensuse.org/repositories/devel:kubic:libcontainers:stable/xUbuntu_${distro_version}/Release.key" \
            | run_as_root "apt-key add -"
          install_package "wget ca-certificates gnupg2 libfuse3-dev fuse-overlayfs"
          ;;
      esac
      install_package "buildah uidmap"
      ;;
    *CentOS*|*Red\ Hat*|*Fedora*)
      run_as_root "yum groupinstall -y 'Development Tools'"
      install_package "curl-devel expat-devel gettext-devel openssl-devel zlib-devel perl autoconf"
      ;;
    *SUSE*|*openSUSE*)
      install_package "buildah shadow"
      ;;
    *RHEL*)
      if [[ "$distro" == *"8"* ]]; then
          run_as_root "dnf module enable -y container-tools:2.0"
          install_package "buildah uidmap"
      elif [[ "$distro" == *"7"* ]]; then
          run_as_root "subscription-manager repos --enable=rhel-7-server-extras-rpms"
          install_package "buildah uidmap"
      fi
      ;;
    *)
      abort "Unsupported distribution."
      ;;
  esac

  set_user_subuids_subgids "$user"

  # Migrate Podman system if applicable
  if is_command_exists podman; then
    log::info "Migrating Podman system configuration..."
    run_as_root "podman system migrate"
  fi

  log::info "Buildah installation and configuration completed."
}


get_free_range() {
  local start=100000
  local end=6000000
  local range_size=65536

  local used_ranges
  if [[ -s /etc/subuid || -s /etc/subgid ]]; then
    used_ranges=$(awk -F: '{print $2 ":" $3}' /etc/subuid /etc/subgid 2>/dev/null | sort -n)
  else
    echo "$start:$range_size"
    return 0
  fi

  while read -r range; do
    local range_start range_size_in_use
    range_start=$(echo "$range" | cut -d: -f1)
    range_size_in_use=$(echo "$range" | cut -d: -f2)
    local range_end=$((range_start + range_size_in_use - 1))

    if ((start + range_size - 1 < range_start)); then
      echo "$start:$range_size"
      return 0
    fi
    start=$((range_end + 1))
  done <<< "$used_ranges"

  echo "$start:$range_size"
  return 0
}

set_user_subuids_subgids() {

  local min_range_size=65536
  local current_user="$1"
  local free_range=$(get_free_range)

  if grep -q "^$current_user:" /etc/subuid; then
    local current_range=$(awk -F: -v user="$current_user" '$1 == user {print $3}' /etc/subuid)
    if (( current_range < min_range_size )); then
      abort "User $current_user has a subuid range smaller than $min_range_size. Please adjust it manually."
    else
      log::info "User $current_user already has a valid subuid range."
    fi
  else
    echo "$current_user:$free_range" | run_as_root "tee -a /etc/subuid"
    log::info "Assigned subuid range $free_range to user $current_user."
  fi

  if grep -q "^$current_user:" /etc/subgid; then
    local current_range=$(awk -F: -v user="$current_user" '$1 == user {print $3}' /etc/subgid)
    if (( current_range < min_range_size )); then
      abort "User $current_user has a subgid range smaller than $min_range_size. Please adjust it manually."
    else
      log::info "User $current_user already has a valid subgid range."
    fi
  else
    echo "$current_user:$free_range" | run_as_root "tee -a /etc/subgid"
    log::info "Assigned subgid range $free_range to user $current_user."
  fi
}

check_user_existence() {
  local user="$1"

  if ! id "$user" &>/dev/null; then
    abort "User '$user' does not exist."
  fi
}

log::info() {
  declare msg="$1"

  printf '[INFO] %b\n' "$msg"
}

log::warn() {
  declare msg="$1"

  printf '[WARNING] %b\n' "$msg" >&2
}

log::err() {
  declare msg="$1"

  printf '[ERROR] %b\n' "$msg" >&2
}

abort() {
  declare msg="$1"

  printf '[FATAL] %b\n' "$msg" >&2
  printf '[FATAL] Aborting.\n' >&2
  exit 1
}

# Author of this function: ko1nksm (Koichi Nakashima)
# License of this function: CC0-1.0 License https://github.com/ko1nksm/getoptions/blob/v3.3.0/LICENSE
getoptions() {
        _error='' _on=1 _no='' _export='' _plus='' _mode='' _alt='' _rest='' _def=''
        _flags='' _nflags='' _opts='' _help='' _abbr='' _cmds='' _init=@empty IFS=' '
        [ $# -lt 2 ] && set -- "${1:?No parser definition}" -
        [ "$2" = - ] && _def=getoptions_parse

        i='                                     '
        while eval "_${#i}() { echo \"$i\$@\"; }"; [ "$i" ]; do i=${i#?}; done

        quote() {
                q="$2'" r=''
                while [ "$q" ]; do r="$r${q%%\'*}'\''" && q=${q#*\'}; done
                q="'${r%????}'" && q=${q#\'\'} && q=${q%\'\'}
                eval "$1=\${q:-\"''\"}"
        }
        code() {
                [ "${1#:}" = "$1" ] && c=3 || c=4
                eval "[ ! \${$c:+x} ] || $2 \"\$$c\""
        }
        sw() { sw="$sw${sw:+|}$1"; }
        kv() { eval "${2-}${1%%:*}=\${1#*:}"; }
        loop() { [ $# -gt 1 ] && [ "$2" != -- ]; }

        invoke() { eval '"_$@"'; }
        prehook() { invoke "$@"; }
        for i in setup flag param option disp msg; do
                eval "$i() { prehook $i \"\$@\"; }"
        done

        args() {
                on=$_on no=$_no export=$_export init=$_init _hasarg=$1 && shift
                while loop "$@" && shift; do
                        case $1 in
                                -?) [ "$_hasarg" ] && _opts="$_opts${1#-}" || _flags="$_flags${1#-}" ;;
                                +?) _plus=1 _nflags="$_nflags${1#+}" ;;
                                [!-+]*) kv "$1"
                        esac
                done
        }
        defvar() {
                case $init in
                        @none) : ;;
                        @export) code "$1" _0 "export $1" ;;
                        @empty) code "$1" _0 "${export:+export }$1=''" ;;
                        @unset) code "$1" _0 "unset $1 ||:" "unset OPTARG ||:; ${1#:}" ;;
                        *)
                                case $init in @*) eval "init=\"=\${${init#@}}\""; esac
                                case $init in [!=]*) _0 "$init"; return 0; esac
                                quote init "${init#=}"
                                code "$1" _0 "${export:+export }$1=$init" "OPTARG=$init; ${1#:}"
                esac
        }
        _setup() {
                [ "${1#-}" ] && _rest=$1
                while loop "$@" && shift; do kv "$1" _; done
        }
        _flag() { args '' "$@"; defvar "$@"; }
        _param() { args 1 "$@"; defvar "$@"; }
        _option() { args 1 "$@"; defvar "$@"; }
        _disp() { args '' "$@"; }
        _msg() { args '' _ "$@"; }

        cmd() { _mode=@ _cmds="$_cmds${_cmds:+|}'$1'"; }
        "$@"
        cmd() { :; }
        _0 "${_rest:?}=''"

        _0 "${_def:-$2}() {"
        _1 'OPTIND=$(($#+1))'
        _1 'while OPTARG= && [ $# -gt 0 ]; do'
        [ "$_abbr" ] && getoptions_abbr "$@"

        args() {
                sw='' validate='' pattern='' counter='' on=$_on no=$_no export=$_export
                while loop "$@" && shift; do
                        case $1 in
                                --\{no-\}*) i=${1#--?no-?}; sw "'--$i'|'--no-$i'" ;;
                                --with\{out\}-*) i=${1#--*-}; sw "'--with-$i'|'--without-$i'" ;;
                                [-+]? | --*) sw "'$1'" ;;
                                *) kv "$1"
                        esac
                done
                quote on "$on"
                quote no "$no"
        }
        setup() { :; }
        _flag() {
                args "$@"
                [ "$counter" ] && on=1 no=-1 v="\$((\${$1:-0}+\$OPTARG))" || v=''
                _3 "$sw)"
                _4 '[ "${OPTARG:-}" ] && OPTARG=${OPTARG#*\=} && set "noarg" "$1" && break'
                _4 "eval '[ \${OPTARG+x} ] &&:' && OPTARG=$on || OPTARG=$no"
                valid "$1" "${v:-\$OPTARG}"
                _4 ';;'
        }
        _param() {
                args "$@"
                _3 "$sw)"
                _4 '[ $# -le 1 ] && set "required" "$1" && break'
                _4 'OPTARG=$2'
                valid "$1" '$OPTARG'
                _4 'shift ;;'
        }
        _option() {
                args "$@"
                _3 "$sw)"
                _4 'set -- "$1" "$@"'
                _4 '[ ${OPTARG+x} ] && {'
                _5 'case $1 in --no-*|--without-*) set "noarg" "${1%%\=*}"; break; esac'
                _5 '[ "${OPTARG:-}" ] && { shift; OPTARG=$2; } ||' "OPTARG=$on"
                _4 "} || OPTARG=$no"
                valid "$1" '$OPTARG'
                _4 'shift ;;'
        }
        valid() {
                set -- "$validate" "$pattern" "$1" "$2"
                [ "$1" ] && _4 "$1 || { set -- ${1%% *}:\$? \"\$1\" $1; break; }"
                [ "$2" ] && {
                        _4 "case \$OPTARG in $2) ;;"
                        _5 '*) set "pattern:'"$2"'" "$1"; break'
                        _4 "esac"
                }
                code "$3" _4 "${export:+export }$3=\"$4\"" "${3#:}"
        }
        _disp() {
                args "$@"
                _3 "$sw)"
                code "$1" _4 "echo \"\${$1}\"" "${1#:}"
                _4 'exit 0 ;;'
        }
        _msg() { :; }

        [ "$_alt" ] && _2 'case $1 in -[!-]?*) set -- "-$@"; esac'
        _2 'case $1 in'
        _wa() { _4 "eval 'set -- $1' \${1+'\"\$@\"'}"; }
        _op() {
                _3 "$1) OPTARG=\$1; shift"
                _wa '"${OPTARG%"${OPTARG#??}"}" '"$2"'"${OPTARG#??}"'
                _4 "$3"
        }
        _3 '--?*=*) OPTARG=$1; shift'
        _wa '"${OPTARG%%\=*}" "${OPTARG#*\=}"'
        _4 ';;'
        _3 '--no-*|--without-*) unset OPTARG ;;'
        [ "$_alt" ] || {
                [ "$_opts" ] && _op "-[$_opts]?*" '' ';;'
                [ ! "$_flags" ] || _op "-[$_flags]?*" - 'OPTARG= ;;'
        }
        [ "$_plus" ] && {
                [ "$_nflags" ] && _op "+[$_nflags]?*" + 'unset OPTARG ;;'
                _3 '+*) unset OPTARG ;;'
        }
        _2 'esac'
        _2 'case $1 in'
        "$@"
        rest() {
                _4 'while [ $# -gt 0 ]; do'
                _5 "$_rest=\"\${$_rest}" '\"\${$(($OPTIND-$#))}\""'
                _5 'shift'
                _4 'done'
                _4 'break ;;'
        }
        _3 '--)'
        [ "$_mode" = @ ] || _4 'shift'
        rest
        _3 "[-${_plus:++}]?*)" 'set "unknown" "$1"; break ;;'
        _3 '*)'
        case $_mode in
                @)
                        _4 "case \$1 in ${_cmds:-*}) ;;"
                        _5 '*) set "notcmd" "$1"; break'
                        _4 'esac'
                        rest ;;
                +) rest ;;
                *) _4 "$_rest=\"\${$_rest}" '\"\${$(($OPTIND-$#))}\""'
        esac
        _2 'esac'
        _2 'shift'
        _1 'done'
        _1 '[ $# -eq 0 ] && { OPTIND=1; unset OPTARG; return 0; }'
        _1 'case $1 in'
        _2 'unknown) set "Unrecognized option: $2" "$@" ;;'
        _2 'noarg) set "Does not allow an argument: $2" "$@" ;;'
        _2 'required) set "Requires an argument: $2" "$@" ;;'
        _2 'pattern:*) set "Does not match the pattern (${1#*:}): $2" "$@" ;;'
        _2 'notcmd) set "Not a command: $2" "$@" ;;'
        _2 '*) set "Validation error ($1): $2" "$@"'
        _1 'esac'
        [ "$_error" ] && _1 "$_error" '"$@" >&2 || exit $?'
        _1 'echo "$1" >&2'
        _1 'exit 1'
        _0 '}'

        [ "$_help" ] && eval "shift 2; getoptions_help $1 $_help" ${3+'"$@"'}
        [ "$_def" ] && _0 "eval $_def \${1+'\"\$@\"'}; eval set -- \"\${$_rest}\""
        _0 '# Do not execute' # exit 1
}

# Author of this function: ko1nksm (Koichi Nakashima)
# License of this function: CC0-1.0 License https://github.com/ko1nksm/getoptions/blob/v3.3.0/LICENSE
getoptions_abbr() {
        abbr() {
                _3 "case '$1' in"
                _4 '"$1") OPTARG=; break ;;'
                _4 '$1*) OPTARG="$OPTARG '"$1"'"'
                _3 'esac'
        }
        args() {
                abbr=1
                shift
                for i; do
                        case $i in
                                --) break ;;
                                [!-+]*) eval "${i%%:*}=\${i#*:}"
                        esac
                done
                [ "$abbr" ] || return 0

                for i; do
                        case $i in
                                --) break ;;
                                --\{no-\}*) abbr "--${i#--\{no-\}}"; abbr "--no-${i#--\{no-\}}" ;;
                                --*) abbr "$i"
                        esac
                done
        }
        setup() { :; }
        for i in flag param option disp; do
                eval "_$i()" '{ args "$@"; }'
        done
        msg() { :; }
        _2 'set -- "${1%%\=*}" "${1#*\=}" "$@"'
        [ "$_alt" ] && _2 'case $1 in -[!-]?*) set -- "-$@"; esac'
        _2 'while [ ${#1} -gt 2 ]; do'
        _3 'case $1 in (*[!a-zA-Z0-9_-]*) break; esac'
        "$@"
        _3 'break'
        _2 'done'
        _2 'case ${OPTARG# } in'
        _3 '*\ *)'
        _4 'eval "set -- $OPTARG $1 $OPTARG"'
        _4 'OPTIND=$((($#+1)/2)) OPTARG=$1; shift'
        _4 'while [ $# -gt "$OPTIND" ]; do OPTARG="$OPTARG, $1"; shift; done'
        _4 'set "Ambiguous option: $1 (could be $OPTARG)" ambiguous "$@"'
        [ "$_error" ] && _4 "$_error" '"$@" >&2 || exit $?'
        _4 'echo "$1" >&2'
        _4 'exit 1 ;;'
        _3 '?*)'
        _4 '[ "$2" = "$3" ] || OPTARG="$OPTARG=$2"'
        _4 "shift 3; eval 'set -- \"\${OPTARG# }\"' \${1+'\"\$@\"'}; OPTARG= ;;"
        _3 '*) shift 2'
        _2 'esac'
}

# Author of this function: ko1nksm (Koichi Nakashima)
# License of this function: CC0-1.0 License https://github.com/ko1nksm/getoptions/blob/v3.3.0/LICENSE
getoptions_help() {
        _width='30,12' _plus='' _leading='  '

        pad() { p=$2; while [ ${#p} -lt "$3" ]; do p="$p "; done; eval "$1=\$p"; }
        kv() { eval "${2-}${1%%:*}=\${1#*:}"; }
        sw() { pad sw "$sw${sw:+, }" "$1"; sw="$sw$2"; }

        args() {
                _type=$1 var=${2%% *} sw='' label='' hidden='' && shift 2
                while [ $# -gt 0 ] && i=$1 && shift && [ "$i" != -- ]; do
                        case $i in
                                --*) sw $((${_plus:+4}+4)) "$i" ;;
                                -?) sw 0 "$i" ;;
                                +?) [ ! "$_plus" ] || sw 4 "$i" ;;
                                *) [ "$_type" = setup ] && kv "$i" _; kv "$i"
                        esac
                done
                [ "$hidden" ] && return 0 || len=${_width%,*}

                [ "$label" ] || case $_type in
                        setup | msg) label='' len=0 ;;
                        flag | disp) label="$sw " ;;
                        param) label="$sw $var " ;;
                        option) label="${sw}[=$var] "
                esac
                [ "$_type" = cmd ] && label=${label:-$var } len=${_width#*,}
                pad label "${label:+$_leading}$label" "$len"
                [ ${#label} -le "$len" ] && [ $# -gt 0 ] && label="$label$1" && shift
                echo "$label"
                pad label '' "$len"
                for i; do echo "$label$i"; done
        }

        for i in setup flag param option disp 'msg -' cmd; do
                eval "${i% *}() { args $i \"\$@\"; }"
        done

        echo "$2() {"
        echo "cat<<'GETOPTIONSHERE'"
        "$@"
        echo "GETOPTIONSHERE"
        echo "}"
}

main "$@"
