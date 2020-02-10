# Assume that binary is already built in pwd.

$GOPATH/bin/fyne package -os linux -icon icon.png
tar xzf xbright.tar.gz
mkdir xbright_1.0
mv usr xbright_1.0/
mkdir xbright_1.0/DEBIAN
echo \
'Package: xbright
Version: 1.0-1
Priority: optional
Architecture: amd64
Maintainer: Micah Parks <micahleviparks@gmail.com>
Description: xbright
 A simple brightness GUI for xrandr.' \
> xbright_1.0/DEBIAN/control
# Use dpkg to build and publish.
dpkg-deb --build xbright_1.0/
