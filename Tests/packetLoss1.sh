
sudo iptables -A INPUT -p udp --dport 42034 -m statistic --mode random --probability 1 -j DROP
sudo iptables -A INPUT -p udp --dport 42035 -m statistic --mode random --probability 1 -j DROP

# Add rules by:      chmod +x filename     followed by    ./filename   
# Delete rules by:   sudo iptables --flush / sudo iptables -F
# To check rules:    sudo iptables -L -n