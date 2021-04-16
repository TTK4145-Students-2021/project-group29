
sudo iptables -A INPUT -p udp --dport 42034 -j ACCEPT
sudo iptables -A INPUT -p udp --sport 42034 -j ACCEPT
# +more for other simulator ports you use
sudo iptables -A INPUT -j DROP