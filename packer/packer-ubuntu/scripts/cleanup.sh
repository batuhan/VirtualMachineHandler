#!/bin/bash -eux

if [[ $ANSIBLE_CLEANUP  =~ true || $ANSIBLE_CLEANUP =~ 1 || $ANSIBLE_CLEANUP =~ yes ]]; then
  # Uninstall Ansible and remove PPA.
  apt -y remove --purge ansible
  apt-add-repository --remove ppa:ansible/ansible

  # Delete Ansible leftovers in home directory
  rm -rf ~/.ansible*
  rm -rf /home/vagrant/.ansible*
fi

# Clean temporary files
rm -rf /tmp/*

# Cleanup apt cache
apt-get -y autoremove --purge
apt-get -y clean
apt-get -y autoclean

# Remove Bash history
unset HISTFILE
rm -f /root/.bash_history
rm -f /home/vagrant/.bash_history

# Clean up log files
find /var/log -type f | while read f; do echo -ne '' > $f; done;

# Clear last login information
>/var/log/lastlog
>/var/log/wtmp
>/var/log/btmp

if [[ $WHITEOUT  =~ true || $WHITEOUT =~ 1 || $WHITEOUT =~ yes ]]; then

	# Whiteout root
	count=$(df --sync -kP / | tail -n1  | awk -F ' ' '{print $4}')
	let count--
	dd if=/dev/zero of=/tmp/whitespace bs=1024 count=$count
	rm /tmp/whitespace

	# Whiteout /boot
	count=$(df --sync -kP /boot | tail -n1 | awk -F ' ' '{print $4}')
	let count--
	dd if=/dev/zero of=/boot/whitespace bs=1024 count=$count
	rm /boot/whitespace

	# Clear out swap and disable until reboot
	set +e
	swapuuid=$(/sbin/blkid -o value -l -s UUID -t TYPE=swap)
	case "$?" in
	    2|0) ;;
	    *) exit 1 ;;
	esac
	set -e
	if [ "x${swapuuid}" != "x" ]; then
	    # Whiteout the swap partition to reduce box size
	    # Swap is disabled till reboot
	    swappart=$(readlink -f /dev/disk/by-uuid/$swapuuid)
	    /sbin/swapoff "${swappart}"
	    dd if=/dev/zero of="${swappart}" bs=1M || echo "dd exit code $? is suppressed"
	    /sbin/mkswap -U "${swapuuid}" "${swappart}"
	fi

	# Zero out the free space to save space in the final image
	dd if=/dev/zero of=/EMPTY bs=1M || echo "dd exit code $? is suppressed"
	rm -f /EMPTY

fi

# Add `sync` so Packer doesn't quit too early, before the large file is deleted.
sync
