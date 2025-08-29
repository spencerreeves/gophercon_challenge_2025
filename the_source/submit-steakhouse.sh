#!/usr/bin/env bash
set -euo pipefail

###
# Steakhouse Submission Script
# Hey Folks - This script creates a (secret) GitHub Gist from a Go file
# This should work for submitting your solutions to the Traitor's Steakhouse sequence
# It DOES require gh (GitHub CLI) and jq (for JSON parsing)
###

BASE_URL="https://station.white-rabbit.dev/le-vrai/"
TIMEOUT=60
VERBOSE=0
URL_ARG=""
FILE_ARG=""

usage() {
  echo "Usage: $0 [-e BASE_URL] [-t SECONDS] [-v] (-u CODE_URL | -f PATH_TO_LOCAL_GO_FILE)"
  echo "  -e  Override base URL (default: $BASE_URL)"
  echo "  -t  Timeout in seconds (default: $TIMEOUT)"
  echo "  -v  Verbose: print raw JSON response"
  echo "  -u  URL to a .go file or a GitHub repo (server fetches it)"
  echo "  -f  Path to a local .go file (will be uploaded as a GitHub Gist using your credentials)"
  exit 1
}

have_cmd() { command -v "$1" >/dev/null 2>&1; }

# Parse flags
while getopts ":e:t:vu:f:" opt; do
  case "$opt" in
    e) BASE_URL="$OPTARG" ;;
    t) TIMEOUT="$OPTARG" ;;
    v) VERBOSE=1 ;;
    u) URL_ARG="$OPTARG" ;;
    f) FILE_ARG="$OPTARG" ;;
    *) usage ;;
  esac
done
shift $((OPTIND-1))

# Allow bare positional URL as convenience if -u not provided and -f not used
if [ -z "$URL_ARG" ] && [ -z "$FILE_ARG" ] && [ $# -ge 1 ]; then
  URL_ARG="$1"
fi

# Validate input
if [ -z "$URL_ARG" ] && [ -z "$FILE_ARG" ]; then
  usage
fi
if [ -n "$URL_ARG" ] && [ -n "$FILE_ARG" ]; then
  echo "Provide either -u or -f, not both." >&2
  exit 1
fi

final_url=""

have_cmd() { command -v "$1" >/dev/null 2>&1; }

create_gist_raw_url() {
  # args: <file_path>
  local path="$1"
  local fname
  fname="$(basename "$path")"
  case "$fname" in
    *.go) : ;;
    *) fname="code.go" ;;
  esac
  local desc
  desc="steakhouse submission $(date -u +%Y-%m-%dT%H:%M:%SZ)"

  if have_cmd gh; then
    # Attempt 1: gh gist create as secret (no external jq needed)
    local gist_url gist_id raw_url
    if [ "$VERBOSE" -eq 1 ]; then
      echo "Using gh gist create (secret)..." >&2
    fi
    # Prefer --secret; if unsupported, this command may fail and we'll fallback
    if ! gist_url="$(gh gist create -s -d "$desc" --filename "$fname" --no-open "$path" 2>&1)"; then
      [ "$VERBOSE" -eq 1 ] && echo "gh gist create failed: $gist_url" >&2
    else
      gist_id="$(printf '%s' "$gist_url" | awk -F/ '{print $NF}')"
      raw_url="$(gh api "/gists/$gist_id" --jq ".files['$fname'].raw_url // (.files | to_entries[0].value.raw_url)")"
      if [ -n "$raw_url" ]; then
        printf '%s' "$raw_url"
        return 0
      fi
      [ "$VERBOSE" -eq 1 ] && echo "Failed to resolve raw_url via gh api" >&2
    fi

    # Attempt 2: gh api with JSON payload (requires jq), set public:false
    if have_cmd jq; then
      [ "$VERBOSE" -eq 1 ] && echo "Falling back to gh api with JSON payload..." >&2
      local resp
      if ! resp="$(jq -Rs --arg d "$desc" --arg f "$fname" '{description:$d,public:false,files:{($f):{content:.}}}' < "$path" | gh api -X POST /gists --input -)"; then
        [ "$VERBOSE" -eq 1 ] && echo "gh api POST /gists failed" >&2
      else
        raw_url="$(printf '%s' "$resp" | gh api --paginate -X GET /rate_limit >/dev/null 2>&1; printf '%s' "$resp" | jq -r --arg f "$fname" '.files[$f].raw_url // (.files | to_entries[0].value.raw_url)')"
        if [ -n "$raw_url" ] && [ "$raw_url" != "null" ]; then
          printf '%s' "$raw_url"
          return 0
        fi
      fi
    fi
  fi

  # Fallback: use curl + token + jq to create via API and read raw_url
  local token
  token="${GH_TOKEN:-${GITHUB_TOKEN:-}}"
  if [ -z "$token" ]; then
    echo "No gh CLI and no GH_TOKEN/GITHUB_TOKEN found for API auth" >&2
    return 1
  fi
  if ! have_cmd jq; then
    echo "jq is required for the token-based fallback" >&2
    return 1
  fi
  # Safely JSON-encode file content
  local content_json
  content_json="$(jq -Rs . < "$path")"
  local payload
  payload=$(cat <<JSON
{
  "description": $(jq -Rn --arg d "$desc" '$d'),
  "public": false,
  "files": { "$fname": { "content": $content_json } }
}
JSON
)
  local resp
  resp="$(curl -sS --fail-with-body -X POST \
    -H "Authorization: token $token" \
    -H "Accept: application/vnd.github+json" \
    https://api.github.com/gists \
    -d "$payload")" || return 1
  printf '%s' "$resp" | jq -r --arg f "$fname" '.files[$f].raw_url // (.files | to_entries[0].value.raw_url)'
}

if [ -n "$FILE_ARG" ]; then
  [ -f "$FILE_ARG" ] || { echo "File not found: $FILE_ARG" >&2; exit 1; }
  echo "Creating GitHub Gist from local file..." >&2
  final_url="$(create_gist_raw_url "$FILE_ARG")" || { echo "Gist creation failed" >&2; exit 1; }
else
  final_url="$URL_ARG"
fi

# Minimal JSON string escape for quotes and backslashes
esc_url="${final_url//\\/\\\\}"
esc_url="${esc_url//\"/\\\"}"
payload='{"url":"'"$esc_url"'"}'

response_and_code="$(curl -sS --fail-with-body \
  -X POST "$BASE_URL/validate" \
  -H "Content-Type: application/json" \
  --max-time "$TIMEOUT" \
  -d "$payload" \
  -w $'\n%{http_code}')"

http_code="${response_and_code##*$'\n'}"
body="${response_and_code%$'\n'*}"

if [ "$VERBOSE" -eq 1 ]; then
  echo "HTTP $http_code"
  echo "$body"
fi

if have_cmd jq; then
  if [ "$http_code" = "200" ]; then
    msg="$(printf '%s' "$body" | jq -r '.msg // empty')"
    [ -n "$msg" ] && { echo "$msg"; exit 0; }
  else
    err="$(printf '%s' "$body" | jq -r '.error // empty')"
    [ -n "$err" ] && { echo "Error: $err" >&2; exit 1; }
  fi
fi

if [ "$http_code" = "200" ]; then
  msg="$(printf '%s' "$body" | sed -n 's/.*\"msg\"[[:space:]]*:[[:space:]]*\"\(.*\)\".*/\1/p')"
  [ -n "$msg" ] && { echo "$msg"; exit 0; }
  echo "$body"
  exit 0
else
  err="$(printf '%s' "$body" | sed -n 's/.*\"error\"[[:space:]]*:[[:space:]]*\"\(.*\)\".*/\1/p')"
  [ -n "$err" ] && { echo "Error: $err" >&2; exit 1; }
  echo "HTTP $http_code" >&2
  echo "$body" >&2
  exit 1
fi
