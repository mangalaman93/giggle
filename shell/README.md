## Software Installation
```
sudo apt-get install moreutils ssmtp git
```

## SMTP setup
This step is optional. Setting SMTP server allows you to
receive emails if there is any issue in syncing the repositories.

###  Modify file /etc/ssmtp/ssmtp.conf
```
root=<YOUR GMAIL>
mailhub=smptp.gmail.com:587
UseSTARTTLS=YES
TLS_CA_File=/etc/ssl/certs/ca-certificates.crt
AuthUser=<GMAIL_USERNAME>
AuthPass=<PASSWORD>
FromLineOverride=NO
```

### Modify file /etc/ssmtp/revaliases
ubuntu:<YOUR_GMAIL>:smtp.gmail.com:587

## Setup Cron Job
* Run command `crontab -e -u <UBUNTU_USERNAME>`
* Add following text -
```MAILTO=<YOUR_EMAIL>
* */6 * * * chronic <LINK_TO_SCRIPT>
```
* Do not add your email if you have no set up the STMP server
* Above settings are for syncing every 6 hours

## Logrotate settings
These are optional too.

### Modify /etc/logrotate.d/overleaf_sync
```
/home/ubuntu/logs/overleaf_sync.log {
  monthly
  rotate 12
  compress
  delaycompress
  missingok
  notifempty
  create 600 ubuntu ubuntu
}
```

## Create dirs
```
mkdir overleaf transient
```

## git settings
```
git config --global http.sslVerify "false"
```
