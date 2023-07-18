#!/bin/bash

# Check for version argument
if [[ -n $1 ]]; then
  VERSION=$1
else
  VERSION="latest"
fi

# Detect platform
if [[ -n $2 ]]; then
  PLATFORM=$2
else
  PLATFORM=$(uname -s | tr '[:upper:]' '[:lower:]')
fi

if [[ $PLATFORM == "linux" || $PLATFORM == "darwin" ]]; then
  echo "$PLATFORM platform is supported"
else
  echo "$PLATFORM platform is not supported - try linux or darwin (for mac)"
  exit 1
fi

# port
if [[ -n $3 ]]; then
  PORT=$3
else
  PORT="7777"
fi

# Set the systemd service name and description
SERVICE_NAME="makr"
SERVICE_DESC="A CLI tool for converting source files to container images."

# Set the GitHub repository and file URL
GITHUB_REPO="https://github.com/origine-run/makr"
FILE_PATH="assets/$VERSION/makr-$VERSION-$PLATFORM"

# Set the username and home directory for the service user
SERVICE_USER="makr"
SERVICE_USER_HOME="/home/$SERVICE_USER"

# Create the service user if it doesn't exist
# if ! id -u $SERVICE_USER > /dev/null 2>&1; then
#   echo "Creating service user: $SERVICE_USER"
#   sudo useradd -m -s /bin/bash -d $SERVICE_USER_HOME $SERVICE_USER
# fi

# Set the Go application binary and working directory
APP_BINARY="/usr/local/bin/makr"
WORKING_DIR="$SERVICE_USER_HOME/makr"

# Download the Go application binary from GitHub
echo "Downloading the Go application binary from GitHub..."
wget -O $APP_BINARY $GITHUB_REPO/raw/main/$FILE_PATH

# Set ownership of the application directory to the service user
sudo chown -R $SERVICE_USER:$SERVICE_USER $WORKING_DIR

# Create a systemd service unit file
sudo tee /etc/systemd/system/$SERVICE_NAME.service > /dev/null << EOF
[Unit]
Description=$SERVICE_DESC
After=network.target

[Service]
ExecStart=$APP_BINARY serve --port $PORT
WorkingDirectory=$WORKING_DIR
Restart=always
User=$SERVICE_USER

[Install]
WantedBy=multi-user.target
EOF

# Reload systemd daemon to read the new service unit file
sudo systemctl daemon-reload

# Enable the systemd service to run on system startup
sudo systemctl enable $SERVICE_NAME

# Start the systemd service
sudo systemctl start $SERVICE_NAME

# Make the service auto run on system startup
systemctl enable --now myapp

# Check the status of the systemd service
# sudo systemctl status $SERVICE_NAME