# Assume that binary is already built in pwd.

$GOPATH/bin/fyne package -os linux -icon icon.png
tar xzf xbright.tar.gz
mkdir xbright
mv usr xbright/
mkdir xbright/DEBIAN
echo \
'Package: xbright
Version: 1.0-1
Priority: optional
Architecture: amd64
Maintainer: Micah Parks <micahleviparks@gmail.com>
Description: xbright
 A simple brightness GUI for xrandr.' \
> xbright/DEBIAN/control
# Use dpkg to build and publish.
dpkg-deb --build xbright/
