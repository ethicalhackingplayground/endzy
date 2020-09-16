git clone https://github.com/nahamsec/bbht.git
cd bbht
chmod +x install.sh
./install.sh
rm -r bbht

go get github.com/003random/getJS
GO111MODULE=auto go get -u -v github.com/projectdiscovery/httpx/cmd/httpx
GO111MODULE=auto go get -u -v github.com/projectdiscovery/subfinder/cmd/subfinder
apt-get install amass

curl -Lo slackcat https://github.com/bcicen/slackcat/releases/download/v1.6/slackcat-1.6-$(uname -s)-amd64
sudo mv slackcat /usr/local/bin/
sudo chmod +x /usr/local/bin/slackcat

slackcat --configure
