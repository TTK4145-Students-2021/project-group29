sudo iptables -A INPUT -p udp --dport 16000 -m statistic --mode random --probability 0.2 -j DROP
sudo iptables -A INPUT -p udp --dport 16001 -m statistic --mode random --probability 0.2 -j DROP

# Add rules by:      chmod +x filename     followed by    ./filename   
# Delete rules by:   sudo iptables --flush / sudo iptables -F
# To check rules:    sudo ipterables -L -n