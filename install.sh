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

 CYGWIN*|MINGW*|MSYS*)
   PLATFORM='windows'
   ;;

 *)
   echo "Platform may or may not be supported. Will attempt to install."
   PLATFORM='linux'
   ;;
esac

if [[ "$(uname -m)" == 'x86_64' ]]; then
  ARCH="amd64"
elif [[ "$(uname -m)" == armv5* ]]; then
  ARCH="armv5"
elif [[ "$(uname -m)" == armv6* ]]; then
  ARCH="armv6"
elif [[ "$(uname -m)" == armv7* ]]; then
  ARCH="armv7"
elif [[ "$(uname -m)" == 'arm64' ]]; then
  ARCH="arm64"
else
  ARCH="386"
fi

if [[ "$(cat /proc/1/cgroup 2> /dev/null | grep docker | wc -l)" > 0 ]] || [ -f /.dockerenv ]; then
  IS_DOCKER=true
else
  IS_DOCKER=false
fi

curl -s -o .ek_tmp_version https://raw.githubusercontent.com/envkey/envkey-source/master/version.txt
VERSION=$(cat .ek_tmp_version)
rm .ek_tmp_version

welcome_envkey () {
  echo "envkey-source $VERSION Quick Install"
  echo "Copyright (c) 2019 Envkey Inc. - MIT License"
  echo "https://github.com/envkey/envkey-source"
  echo ""
}

cleanup () {
  rm envkey-source.tar.gz
  rm -f envkey-source
}

download_envkey () {
  echo "Downloading envkey-source binary for ${PLATFORM}-${ARCH}"
  url="https://github.com/envkey/envkey-source/releases/download/v${VERSION}/envkey-source_${VERSION}_${PLATFORM}_${ARCH}.tar.gz"
  echo "Downloading tarball from ${url}"
  curl -s -L -o envkey-source.tar.gz "${url}"

  tar zxf envkey-source.tar.gz envkey-source.exe 2> /dev/null
  tar zxf envkey-source.tar.gz envkey-source 2> /dev/null

  if [ "$PLATFORM" == "darwin" ] || $IS_DOCKER ; then
    if [[ -d /usr/local/bin ]]; then
      mv envkey-source /usr/local/bin/
      echo "envkey-source is installed in /usr/local/bin"
    else
      echo >&2 'Error: /usr/local/bin does not exist. Create this directory with appropriate permissions, then re-install.'
      cleanup
      exit 1
    fi
  elif [ "$PLATFORM" == "windows" ]; then
    # ensure $HOME/bin exists (it's in PATH but not present in default git-bash install)
    mkdir $HOME/bin 2> /dev/null
    mv envkey-source.exe $HOME/bin/
    echo "envkey-source is installed in $HOME/bin"
  else
    sudo mv envkey-source /usr/local/bin/
    echo "envkey-source is installed in /usr/local/bin"
  fi
}

welcome_envkey
download_envkey
cleanup

echo "Installation complete. Info:"
echo ""
envkey-source -h
