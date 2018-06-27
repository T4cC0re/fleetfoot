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

{% for interface in interfaces %}
{% if interfaces[interface].enable_intel_driver_bug_workaround is defined and interfaces[interface].enable_intel_driver_bug_workaround %}
ethtool -K {{interface}} gso off gro off tso off
{% endif %}
{% endfor %}

if [[ $0 =~ ip-up ]]; then
  echo "Executing ip-up hook..."
  ip route del default
  ip route add default dev $interface
  iptables -t nat -F PREROUTING
  iptables -F FORWARD
  iptables -F INPUT
  iptables -t nat -A POSTROUTING -o $interface -j MASQUERADE
  iptables -A INPUT -i $interface -m state --state ESTABLISHED,RELATED -j ACCEPT
  iptables -A INPUT -i $interface -j LOG --log-prefix "fleetfoot-firewall::INPUT:DROP:" --log-level 6
  iptables -A INPUT -i $interface -j DROP
{% if enable_mss_clamping %}
  iptables -A FORWARD -p tcp --tcp-flags SYN,RST SYN -j TCPMSS --clamp-mss-to-pmtu
{% endif %}

{# iptables boilerplate to make optional params possible #}
{% if portForwarding is defined %}
{% for portConf in portForwarding %}
{% if portConf.externalPort is defined %}
{% set externalPort = portConf.externalPort %}
{% else %}
{% set externalPort = portConf.targetPort %}
{% endif %}
{% if portConf.proto is defined %}
{% set proto = portConf.proto %}
{% else %}
{% set proto = 'tcp' %}
{% endif %}
  echo "Forward external {{ externalPort }}/{{ proto }} to {{ portConf.target }} on port {{ portConf.targetPort }}..."
  iptables -t nat -A PREROUTING -p {{ proto }} -d ${local_ip} --dport {{ externalPort }} -j DNAT --to-destination {{ portConf.target }}:{{ portConf.targetPort }}
  iptables -A FORWARD -p {{ proto }} -d {{ portConf.target }} --dport {{ portConf.targetPort }} -m state --state NEW,ESTABLISHED,RELATED -j ACCEPT
{% endfor %}
{% endif %}

  echo "Executed ip-up hook"
fi
