#!/usr/bin/env sh

set -eu

echo "==> installing Terraform CLI..."

if command -v terraform >/dev/null 2>&1; then
	echo "terraform is already installed; leaving existing version in place"
	terraform version
	exit 0
fi

os="$(uname -s | tr '[:upper:]' '[:lower:]')"
arch="$(uname -m)"

case "$arch" in
	x86_64|amd64)
		arch="amd64"
		;;
	aarch64|arm64)
		arch="arm64"
		;;
	i386|i686)
		arch="386"
		;;
	armv7l|armv6l)
		arch="arm"
		;;
	*)
		echo "unsupported architecture: $arch" >&2
		exit 1
		;;
esac

case "$os" in
	linux|darwin)
		;;
	*)
		echo "unsupported operating system: $os" >&2
		exit 1
		;;
esac

bin_dir="$(go env GOBIN)"
if [ -z "$bin_dir" ]; then
	bin_dir="$(go env GOPATH)/bin"
fi
mkdir -p "$bin_dir"

tmp_dir="$(mktemp -d)"
trap 'rm -rf "$tmp_dir"' EXIT INT TERM

tf_version="$(
	curl -fsSL https://releases.hashicorp.com/terraform/ \
	| sed -n 's|.*terraform/\(1\.[0-9][0-9.]*\)/.*|\1|p' \
	| sort -V \
	| tail -n 1
)"

if [ -z "$tf_version" ]; then
	echo "failed to determine latest stable Terraform 1.x version" >&2
	exit 1
fi

zip_path="$tmp_dir/terraform_${tf_version}_${os}_${arch}.zip"
download_url="https://releases.hashicorp.com/terraform/${tf_version}/terraform_${tf_version}_${os}_${arch}.zip"

curl -fsSL "$download_url" -o "$zip_path"
unzip -oq "$zip_path" -d "$tmp_dir"
install "$tmp_dir/terraform" "$bin_dir/terraform"

"$bin_dir/terraform" version
