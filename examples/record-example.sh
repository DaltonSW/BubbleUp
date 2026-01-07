#!/usr/bin/env bash
set -euo pipefail

# record-demo.sh
#
# Purpose:
#   Runs VHS to generate a GIF/MP4 demo recording from a .tape file.
#   This script validates prerequisites so contributors can run it without
#   knowing anything about VHS/ttyd/ffmpeg/fonts.
#
# Why the Nerd Font check matters:
#   Nerd Font PUA glyph sizing can differ depending on which font face is used.
#   For correct sizing in ttyd/VHS captures, use:
#     JetBrainsMono Nerd Font Mono
#
# Typical usage:
#   ./record-demo.sh
#   ./record-demo.sh -t example_main.tape
#   ./record-demo.sh --force
#
# Notes:
#   - The .tape file usually controls output filename via an `Output ...` line.
#   - On macOS, VHS may need VHS_NO_SANDBOX=1 depending on your environment.
#     This script enables it by default on macOS.

SCRIPT_NAME="$(basename "$0")"
OS="$(uname -s)"

DEFAULT_TAPE="example_main.tape"
TAPE_FILE="$DEFAULT_TAPE"
VHS_NO_SANDBOX_DEFAULT=""
FONT_FAMILY_REQUIRED="JetBrainsMono Nerd Font Mono"
FORCE=0

usage() {
  cat <<EOF
Usage:
  $SCRIPT_NAME [-t <tape-file>] [-f] [-h]

Options:
  -t <tape-file>   Path to VHS tape file (default: $DEFAULT_TAPE)
  -f, --force      Bypass prerequisite checks (tools/fonts) and run anyway
  -h               Show this help

What this does:
  - Verifies required tools: go, vhs, ttyd, ffmpeg
  - Verifies required font family is installed: "$FONT_FAMILY_REQUIRED"
  - Runs VHS to produce the recording defined by the tape file

Examples:
  $SCRIPT_NAME
  $SCRIPT_NAME -t example_main.tape
  $SCRIPT_NAME --force
  $SCRIPT_NAME -f -t example_main.tape

EOF
}

die() {
  echo "Error: $*" >&2
  exit 1
}

have_cmd() {
  command -v "$1" >/dev/null 2>&1
}

require_cmd() {
  local cmd="$1"
  local hint="$2"
  have_cmd "$cmd" || die "Missing required command: $cmd
$hint"
}

# system_profiler is slow and its output is not stable for
# grepping specific family names; prefer checking installed
# font files in ~/Library/Fonts and /Library/Fonts.
font_installed_macos() {
  # Fast path: check common install locations for Nerd Font files.
  # Homebrew font casks typically install into ~/Library/Fonts or /Library/Fonts.
  local patterns=(
    "JetBrainsMono*Nerd*Mono*"
    "JetBrains Mono*Nerd*Mono*"
    "JetBrainsMonoNerdFontMono*"
  )

  local dirs=(
    "$HOME/Library/Fonts"
    "/Library/Fonts"
  )

  local d p
  for d in "${dirs[@]}"; do
    [[ -d "$d" ]] || continue
    for p in "${patterns[@]}"; do
      # Use compgen to avoid "no matches found" issues with strict shells.
      if compgen -G "$d/$p*.ttf" >/dev/null || compgen -G "$d/$p*.otf" >/dev/null; then
        return 0
      fi
    done
  done

  return 1
}

font_installed_linux() {
  # Prefer fc-list if available (fontconfig). Most distros have it installed.
  if have_cmd fc-list; then
    fc-list : family 2>/dev/null | grep -Fq "$FONT_FAMILY_REQUIRED"
    return $?
  fi

  # Fallback: try listing common font directories; less reliable.
  # This is best-effort and may miss user-installed fonts.
  find /usr/share/fonts /usr/local/share/fonts "$HOME/.local/share/fonts" \
    -type f 2>/dev/null | head -n 1 >/dev/null || return 1
  return 1
}

require_font() {
  case "$OS" in
    Darwin)
      font_installed_macos || die "Required font family not found: \"$FONT_FAMILY_REQUIRED\"

Install suggestion (macOS, Homebrew):
  brew install --cask font-jetbrains-mono-nerd-font

Then set your tape to:
  Set FontFamily \"$FONT_FAMILY_REQUIRED\"
"
      ;;
    Linux)
      font_installed_linux || die "Required font family not found: \"$FONT_FAMILY_REQUIRED\"

Install suggestion (Linux):
  - Install a Nerd Font (JetBrains Mono Nerd Font), then ensure the *Mono* family is available.
  - If you have fontconfig installed, you can verify with:
      fc-list | grep -i \"jetbrains\" | grep -i \"nerd\"

Then set your tape to:
  Set FontFamily \"$FONT_FAMILY_REQUIRED\"
"
      ;;
    *)
      die "Unsupported OS: $OS (expected Darwin or Linux)"
      ;;
  esac
}

# Parse args (support -f and --force)
if [[ "${1:-}" == "-h" || "${1:-}" == "--help" ]]; then
  usage
  exit 0
fi

# Handle long opts before getopts
while [[ $# -gt 0 ]]; do
  case "$1" in
    --force)
      FORCE=1
      shift
      ;;
    --)
      shift
      break
      ;;
    -*)
      break
      ;;
    *)
      break
      ;;
  esac
done

while getopts ":t:hf" opt; do
  case "$opt" in
    t) TAPE_FILE="$OPTARG" ;;
    f) FORCE=1 ;;
    h) usage; exit 0 ;;
    :)
      usage
      die "Missing value for -$OPTARG"
      ;;
    \?)
      usage
      die "Unknown option: -$OPTARG"
      ;;
  esac
done
shift $((OPTIND-1))

# No unexpected positional args
if [[ $# -ne 0 ]]; then
  usage
  die "Unexpected argument(s): $*"
fi

[[ -f "$TAPE_FILE" ]] || die "Tape file not found: $TAPE_FILE"

# OS-specific defaults
case "$OS" in
  Darwin)
    # VHS sandboxing can cause issues on macOS; enable by default.
    VHS_NO_SANDBOX_DEFAULT="1"
    ;;
  Linux)
    VHS_NO_SANDBOX_DEFAULT=""
    ;;
  *)
    die "Unsupported OS: $OS (expected Darwin or Linux)"
    ;;
esac

if [[ "$FORCE" -eq 1 ]]; then
  echo "Warning: --force enabled; skipping prerequisite checks (tools/fonts)." >&2
  echo >&2
else
  # Requirements
  require_cmd go     "Install Go from https://go.dev/ or your package manager."
  require_cmd vhs    "Install VHS from https://github.com/charmbracelet/vhs (or your package manager)."
  require_cmd ttyd   "Install ttyd (required by VHS). macOS: brew install ttyd"
  require_cmd ffmpeg "Install ffmpeg. macOS: brew install ffmpeg"

  require_font
fi

# Informational echo
echo "Using tape: $TAPE_FILE"
echo "OS: $OS"
if [[ "$FORCE" -eq 0 ]]; then
  echo "Verified: go, vhs, ttyd, ffmpeg"
  echo "Verified font family: $FONT_FAMILY_REQUIRED"
else
  echo "Checks: SKIPPED (--force)"
fi
echo

# Run VHS
# - On macOS: set VHS_NO_SANDBOX=1 by default unless caller already set it.
# - On Linux: don't force it.
if [[ "$OS" == "Darwin" ]]; then
  : "${VHS_NO_SANDBOX:=$VHS_NO_SANDBOX_DEFAULT}"
  export VHS_NO_SANDBOX
  echo "Running: VHS_NO_SANDBOX=$VHS_NO_SANDBOX vhs \"$TAPE_FILE\""
else
  echo "Running: vhs \"$TAPE_FILE\""
fi

vhs "$TAPE_FILE"
