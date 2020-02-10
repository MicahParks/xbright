# Assume that binary is already built in pwd.

$GOPATH/bin/fyne package -os linux -icon icon.png
tar xzf xbright.tar.gz
mkdir xbright
mv usr xbright/
mkdir xbright/DEBIAN
cat <<EOF | sudo tee xbright/DEBIAN/control
# TODO
EOF
# Use dpkg to build and publish.

