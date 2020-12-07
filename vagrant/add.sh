#!/bin/sh
set -eux
VM=$1
ADDR=$2

SERVER_KEY=$(vagrant ssh registry -- cat /etc/wireguard/publickey)

vagrant ssh $VM -- sudo bash -s << EOF
wg genkey | tee /etc/wireguard/privatekey | wg pubkey | tee /etc/wireguard/publickey
# linux config
cat > /etc/wireguard/wg0.conf << CLIENTEOF
[Interface]
PrivateKey = \$(cat /etc/wireguard/privatekey)
Address = $ADDR/24
ListenPort = 51820
[Peer]
PublicKey = $SERVER_KEY
Endpoint = 192.168.33.10:51820
AllowedIPs = 192.168.100.10/32
CLIENTEOF
chmod 600 /etc/wireguard/{privatekey,wg0.conf}
chmod 644 /etc/wireguard/publickey
chmod 711 /etc/wireguard
# vpp config
cp /vagrant/startup.conf /etc/vpp/startup.conf
cat > /etc/vpp/init.vpp << CLIENTEOF
create int af_xdp host-if enp0s8
set int mac addr enp0s8/0 \$(cat /sys/class/net/enp0s8/address)
set int ip addr enp0s8/0 \$(ip addr show enp0s8|awk '\$1=="inet"{print \$2}')
wireguard create listen-port 51820 private-key \$(cat /etc/wireguard/privatekey) src \$(ip addr show enp0s8|awk -F '[ /]' '\$5=="inet"{print \$6}')
wireguard peer add wg0 public-key $SERVER_KEY endpoint 192.168.33.10 allowed-ip 192.168.100.10/32 port 51820
set int unnum wg0 use enp0s8/0
create tap id 0 host-ip4-addr $ADDR/24 host-if-name wg0 host-mtu-size 1420 tun gso
set int unnum tun0 use enp0s8/0
set int feature gso wg0 enable
ip route add $ADDR/32 via tun0
ip route add 192.168.100.0/24 via wg0
set int st enp0s8/0 up
set int st wg0 up
set int st tun0 up
wait 1
set int rx-mode enp0s8/0 adaptive
set int rx-mode tun0 adaptive
CLIENTEOF
EOF

CLIENT_KEY=$(vagrant ssh $VM -- cat /etc/wireguard/publickey)

vagrant ssh registry -- sudo wg set wg0 peer $CLIENT_KEY allowed-ips $ADDR/32
