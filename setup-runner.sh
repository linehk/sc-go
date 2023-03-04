#!/bin/bash

# Check if curl is installed
if ! [ -x "$(command -v curl)" ]; then
    echo "curl is not installed. Installing curl..."
    sudo apt-get update
    sudo apt-get install -y curl
    echo "curl installed successfully."
fi

# Check if Docker is installed
if ! [ -x "$(command -v docker)" ]; then
    echo "Docker is not installed. Installing Docker..."
    curl -fsSL https://get.docker.com -o get-docker.sh
    sudo sh get-docker.sh
    sudo usermod -aG docker $USER
    echo "Docker installed successfully."
fi

# Check if Dockerx is installed
if ! [ -x "$(command -v dockerx)" ]; then
    echo "Dockerx is not installed. Installing Dockerx..."
    sudo curl -L https://github.com/mayflower/docker-x/releases/latest/download/docker-x_linux_amd64 -o /usr/local/bin/dockerx
    sudo chmod +x /usr/local/bin/dockerx
    echo "Dockerx installed successfully."
fi

# Install GitHub runner
if [ ! -d ~/actions-runner ]; then
    echo "GitHub runner is not installed. Installing GitHub runner..."
    mkdir ~/actions-runner && cd ~/actions-runner
    LATEST_RELEASE=$(curl -s https://api.github.com/repos/actions/runner/releases/latest | grep -oP '"tag_name": "\K(.*)(?=")')
    curl -o actions-runner-linux-x64.tar.gz -L https://github.com/actions/runner/releases/download/${LATEST_RELEASE}/actions-runner-linux-x64-${LATEST_RELEASE}.tar.gz
    tar xzf ./actions-runner-linux-x64.tar.gz
    echo "GitHub runner installed successfully."
fi

# Configure GitHub runner
cd ~/actions-runner
./config.sh --url $REPO --token $TOKEN --unattended
echo "GitHub runner configured successfully."

# Launch GitHub runner as a service
./svc.sh install
./svc.sh start
echo "GitHub runner launched successfully."