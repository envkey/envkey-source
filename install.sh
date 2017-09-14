#!/usr/bin/env bash

case "$(uname -s)" in
 Darwin)
   platform='darwin'
   ;;

 Linux)
   platform='linux'
   ;;

 FreeBSD)
   platform='freebsd'
   ;;

 CYGWIN*|MINGW32*|MSYS*)
   echo "This install script is not compatible with windows. Download the latest .exe release and put it in your PATH: https://github.com/envkey/envkey-source/releases"
   exit 1
   ;;

   echo "Platform may or may not be supported. Will attempt to install."
   platform= 'linux'
   ;;
esac

if [[ "$(uname -m)" == 'x86_64' ]]; then
  arch="amd64"
else
  arch="386"
end

curl -s -o version.txt https://raw.githubusercontent.com/envkey/envkey-source/master/version.txt
VERSION=$(cat version.txt)
rm version.txt

function welcome_envkey {
  echo "envkey-source ($VERSION) Quick Install"
  echo "Copyright (c) 2017 Envkey Inc. - MIT License"
  echo "https://github.com/envkey/envkey-source"
  echo ""
}

function download_envkey {
  echo "Downloading envkey-source binary for ${platform}-${arch}"
  curl -s -o envkey-source.tar.gz "https://raw.githubusercontent.com/envkey/envkey-source/master/dist/envkey-source_${platform}_${arch}.tar.gz"
  tar zxf envkey-source.tar.gz

  sudo mv envkey-source/envkey-source /usr/local/bin/
  echo "envkey-source is installed in /usr/local/bin"

  rm envkey-source.tar.gz
  rm -rf envkey-source
}


welcome_envkey
download_envkey

echo "Installation complete."
/usr/local/bin/envkey-source -h
