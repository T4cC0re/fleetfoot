#!/bin/bash

readonly interface="$1"
shift
readonly tty="$1"
shift
readonly speed="$1"
shift
readonly local_ip="$1"
shift
readonly remote_ip="$1"
shift
readonly ipparam="$1"

exec 2>&1 > >(tee -a /tmp/hook.log)
date

echo "Hook executed as '$0' with parameters:"
echo "interface : $interface"
echo "tty ......: $tty"
echo "speed ....: $speed"
echo "local IP .: $local_ip"
echo "remote IP : $remote_ip"
echo "ipparam ..: $ipparam"

if [[ $0 =~ ip-up ]]; then
  echo "Executing ip-up hook..."
  ip route del default
  ip route add default dev $interface
  iptables -t nat -A POSTROUTING -o $interface -j MASQUERADE
  iptables -A INPUT -i $interface -m state --state ESTABLISHED,RELATED -j ACCEPT
  iptables -A INPUT -i $interface -j DROP
{% if enable_mss_clamping %}
  iptables -A FORWARD -p tcp --tcp-flags SYN,RST SYN -j TCPMSS --clamp-mss-to-pmtu
{% endif %}
  echo "Executed ip-up hook"
fi
