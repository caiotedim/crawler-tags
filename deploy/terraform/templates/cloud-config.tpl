#cloud-config
manage_etc_hosts: "localhost"
ssh_authorized_keys:
%{ for key in ssh_keys.sre ~}
- ${key}
%{ endfor }