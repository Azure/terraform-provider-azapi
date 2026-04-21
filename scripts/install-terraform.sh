#!/usr/bin/env sh

set -eu

echo "==> installing Terraform CLI..."

set_terraform_path() {
	terraform_path="$1"
	export TF_ACC_TERRAFORM_PATH="$terraform_path"
	echo "TF_ACC_TERRAFORM_PATH=$TF_ACC_TERRAFORM_PATH"
	if [ -n "${TF_BUILD:-}" ]; then
		echo "##vso[task.setvariable variable=TF_ACC_TERRAFORM_PATH]$TF_ACC_TERRAFORM_PATH"
	fi
}

if command -v terraform >/dev/null 2>&1; then
	echo "terraform is already installed; leaving existing version in place"
	set_terraform_path "$(command -v terraform)"
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

if [ -z "${TERRAFORM_VERSION:-}" ]; then
	echo "TERRAFORM_VERSION is not set" >&2
	exit 1
fi
tf_version="$TERRAFORM_VERSION"

base_url="https://releases.hashicorp.com/terraform/${tf_version}"
zip_name="terraform_${tf_version}_${os}_${arch}.zip"
zip_path="$tmp_dir/$zip_name"
sums_path="$tmp_dir/terraform_${tf_version}_SHA256SUMS"

curl -fsSL "${base_url}/${zip_name}" -o "$zip_path"
curl -fsSL "${base_url}/terraform_${tf_version}_SHA256SUMS" -o "$sums_path"

expected_sum="$(grep " ${zip_name}$" "$sums_path" | awk '{ print $1 }')"
if [ -z "$expected_sum" ]; then
	echo "checksum entry for ${zip_name} not found in SHA256SUMS" >&2
	exit 1
fi
actual_sum="$(sha256sum "$zip_path" | awk '{ print $1 }')"
if [ "$actual_sum" != "$expected_sum" ]; then
	echo "checksum mismatch for ${zip_name}: expected ${expected_sum}, got ${actual_sum}" >&2
	exit 1
fi

unzip -oq "$zip_path" -d "$tmp_dir"
install "$tmp_dir/terraform" "$bin_dir/terraform"

set_terraform_path "$bin_dir/terraform"
"$bin_dir/terraform" version
