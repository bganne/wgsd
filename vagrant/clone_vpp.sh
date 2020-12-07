#!/bin/bash
set -eux

VPP_COMMIT=2df2f75ec
VPP_DIR=$1

if [ ! -d $VPP_DIR ]; then
	git clone "https://gerrit.fd.io/r/vpp" $VPP_DIR
	pushd $VPP_DIR
	git reset --hard ${VPP_COMMIT}
else
	pushd $VPP_DIR
	git fetch "https://gerrit.fd.io/r/vpp"
    git reset --hard ${VPP_COMMIT}
    git clean -dfx
fi

git fetch "https://gerrit.fd.io/r/vpp" refs/changes/46/30146/3 && git cherry-pick FETCH_HEAD # 30146: wireguard: return public key in api | https://gerrit.fd.io/r/c/vpp/+/30146
git fetch "https://gerrit.fd.io/r/vpp" refs/changes/54/30154/1 && git cherry-pick FETCH_HEAD # 30154: wireguard: run feature after gso | https://gerrit.fd.io/r/c/vpp/+/30154
git fetch "https://gerrit.fd.io/r/vpp" refs/changes/49/30249/1 && git cherry-pick FETCH_HEAD # 30249: virtio: fix the len offset | https://gerrit.fd.io/r/c/vpp/+/30249
git fetch "https://gerrit.fd.io/r/vpp" refs/changes/56/30256/2 && git cherry-pick FETCH_HEAD # 30256: virtio: fix the offloads in tx path

make build-release
rm -f ./build-root/*.deb ./build-root/*.changes ./build-root/*.buildinfo
make pkg-deb
popd

rm -rf ./debs
mkdir -p ./debs
cp ${VPP_DIR}/build-root/*.deb ./debs/
