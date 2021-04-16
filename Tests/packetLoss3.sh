sudo iptables -A INPUT -p udp --dport 16000 -j ACCEPT
sudo iptables -A INPUT -p udp --sport 16000 -j ACCEPT

sudo iptables -A INPUT -p udp --dport 16001 -j ACCEPT
sudo iptables -A INPUT -p udp --sport 16001 -j ACCEPT

sudo iptables -A INPUT -m statistic --mode random --probability 0.2 -j DROP