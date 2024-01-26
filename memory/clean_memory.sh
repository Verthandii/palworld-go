sudo sync; echo 3 > /proc/sys/vm/drop_caches
sudo echo 1 > /proc/sys/vm/drop_caches
sudo sync; echo 2 > /proc/sys/vm/drop_caches
sudo swapoff -a && sudo swapon -a