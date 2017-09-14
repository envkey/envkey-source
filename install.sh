# !/usr/bin/env bash

case "$(uname -s)" in
 Darwin)
   PLATFORM='darwin'
   ;;

 Linux)
   PLATFORM='linux'
   ;;

 FreeBSD)
   PLATFORM='freebsd'
   ;;

 CYGWIN*|MINGW32*|MSYS*)
   echo "This install script is not compatible with windows. Download the latest .exe release and put it in your PATH: https://github.com/envkey/envkey-source/releases"
   exit 1
   ;;

 *)
   echo "Platform may or may not be supported. Will attempt to install."
   PLATFORM='linux'
   ;;
esac

if [[ "$(uname -m)" == 'x86_64' ]]; then
  ARCH="amd64"
else
  ARCH="386"
fi

curl -s -o .ek_tmp_version https://raw.githubusercontent.com/envkey/envkey-source/master/version.txt
VERSION=$(cat .ek_tmp_version)
rm .ek_tmp_version

function welcome_envkey {
  echo "envkey-source $VERSION Quick Install"
  echo "Copyright (c) 2017 Envkey Inc. - MIT License"
  echo "https://github.com/envkey/envkey-source"
  echo ""
}

function download_envkey {
  echo "Downloading envkey-source binary for ${PLATFORM}-${ARCH}"
  url="https://raw.githubusercontent.com/envkey/envkey-source/master/dist/envkey-source_${VERSION}_${PLATFORM}_${ARCH}.tar.gz"
  echo "Downloading tarball from ${url}"
  curl -s -o envkey-source.tar.gz "${url}"
  tar zxf envkey-source.tar.gz

  if [ "$PLATFORM" == "darwin" ]; then
    mv envkey-source /usr/local/bin/
  else
    sudo mv envkey-source /usr/local/bin/
  fi

  echo "envkey-source is installed in /usr/local/bin"

  rm envkey-source.tar.gz
  rm -f envkey-source
}

welcome_envkey
download_envkey

echo "Installation complete. Info:"
echo ""
/usr/local/bin/envkey-source -h
