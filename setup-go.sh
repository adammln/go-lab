# create download dir if not exists
mkdir $HOME/Downloads
cd $HOME/Downloads

# download installation
wget https://go.dev/dl/go1.18.3.linux-amd64.tar.gz

echo "Extracting zip"
# remove any previous go installation
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.18.3.linux-amd64.tar.gz

echo "Setting up environment variable"
# add /usr/local/go/bin to PATH env var
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source $HOME/.bashrc

echo "Installation Done!"
go version
