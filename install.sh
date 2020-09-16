mkdir -p db

git clone https://github.com/nahamsec/bbht.git
cd bbht
chmod +x install.sh
./install.sh
rm -r bbht

go get github.com/003random/getJS
GO111MODULE=auto go get -u -v github.com/projectdiscovery/httpx/cmd/httpx
GO111MODULE=auto go get -u -v github.com/projectdiscovery/subfinder/cmd/subfinder
apt-get install amass

go get github.com/balook/teleman
