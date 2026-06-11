#!/bin/bash
# ChromaDesk Signaling Server — Nginx reverse proxy setup
# Usage: ./setup-nginx.sh <DOMAIN> <PORT>
#
# This script is called by deploy scripts when --domain is specified.
# It can also be run standalone to configure/reconfigure Nginx.

set -e

DOMAIN="$1"
PORT="$2"

if [ -z "$DOMAIN" ] || [ -z "$PORT" ]; then
    echo "Usage: $0 <DOMAIN> <PORT>"
    echo "  e.g. $0 example.com 8000"
    exit 1
fi

echo "Configuring Nginx reverse proxy for $DOMAIN → 127.0.0.1:$PORT ..."

# ---- Install Nginx if missing ----
if ! command -v nginx &>/dev/null; then
    echo "Installing Nginx..."
    if command -v dnf &>/dev/null; then
        sudo dnf install -y nginx
    elif command -v yum &>/dev/null; then
        sudo yum install -y nginx
    elif command -v apt-get &>/dev/null; then
        sudo apt-get update && sudo apt-get install -y nginx
    fi
fi

if ! command -v nginx &>/dev/null; then
    echo "ERROR: Failed to install Nginx. Please install it manually and re-run."
    exit 1
fi

# ---- Locate nginx.conf (system, BT Panel, custom) ----
NGINX_CONF=""
for candidate in \
    /etc/nginx/nginx.conf \
    /www/server/nginx/conf/nginx.conf \
    /usr/local/nginx/conf/nginx.conf \
    /opt/nginx/conf/nginx.conf; do
    if [ -f "$candidate" ]; then
        NGINX_CONF="$candidate"
        break
    fi
done

if [ -z "$NGINX_CONF" ]; then
    NGINX_CONF=$(nginx -t 2>&1 | grep -oP 'file \K\S+(?= syntax)' || true)
fi

if [ -z "$NGINX_CONF" ] || [ ! -f "$NGINX_CONF" ]; then
    NGINX_CONF="/etc/nginx/nginx.conf"
    echo "Creating minimal $NGINX_CONF ..."
    sudo mkdir -p /etc/nginx/conf.d /var/log/nginx
    sudo tee "$NGINX_CONF" > /dev/null << 'MINCONF'
user nginx;
worker_processes auto;
error_log /var/log/nginx/error.log;
pid /run/nginx.pid;

events {
    worker_connections 1024;
}

http {
    include       /etc/nginx/mime.types;
    default_type  application/octet-stream;
    sendfile      on;
    keepalive_timeout 65;
    include /etc/nginx/conf.d/*.conf;
}
MINCONF
fi

echo "Using nginx config: $NGINX_CONF"

# ---- Detect vhost include directory ----
NGINX_CONF_DIR=""
while IFS= read -r inc_path; do
    inc_dir="$(dirname "$inc_path")"
    if [ -d "$inc_dir" ]; then
        NGINX_CONF_DIR="$inc_dir"
        break
    fi
done < <(grep -E '^\s*include\s+' "$NGINX_CONF" \
    | grep -v 'mime' | grep -v 'lua' | grep -v 'proxy' | grep -v 'php' | grep -v 'enable-' | grep -v 'fastcgi' | grep -v '/tcp/' | grep -v '/stream/' \
    | grep -oP 'include\s+\K[^;]+' \
    | grep '\*\.conf' \
    | tail -n 5)

if [ -z "$NGINX_CONF_DIR" ]; then
    NGINX_CONF_DIR="$(dirname "$NGINX_CONF")/conf.d"
fi

echo "Using vhost directory: $NGINX_CONF_DIR"
sudo mkdir -p "$NGINX_CONF_DIR"

# ---- Write ChromaDesk vhost config ----
sudo tee "$NGINX_CONF_DIR/chromadesk.conf" > /dev/null << NGINX_EOF
map \$http_upgrade \$connection_upgrade {
    default upgrade;
    ''      close;
}

upstream chromadesk_signaling {
    server 127.0.0.1:$PORT;
}

server {
    listen 80;
    server_name $DOMAIN;

    client_max_body_size 100M;

    # WebSocket (long-lived connections, server pings every 60s)
    location /signal/ {
        proxy_pass http://chromadesk_signaling;
        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection \$connection_upgrade;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;

        proxy_read_timeout 300s;
        proxy_send_timeout 300s;
        proxy_buffering off;
    }

    # HTTP API and static files
    location / {
        proxy_pass http://chromadesk_signaling;
        proxy_http_version 1.1;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
        proxy_read_timeout 30s;
    }
}
NGINX_EOF

# ---- Ensure nginx.conf includes the vhost directory ----
NGINX_CONF_DIR_ESCAPED=$(echo "$NGINX_CONF_DIR" | sed 's/[\/&]/\\&/g')
if ! grep -q "include.*${NGINX_CONF_DIR_ESCAPED}" "$NGINX_CONF" 2>/dev/null; then
    echo "Adding include for $NGINX_CONF_DIR to $NGINX_CONF ..."
    if grep -qP 'http\s*\{' "$NGINX_CONF"; then
        sudo sed -i '/http\s*{/a \    include '"$NGINX_CONF_DIR"'/*.conf;' "$NGINX_CONF"
    elif grep -q '^http$' "$NGINX_CONF"; then
        sudo sed -i '/^http$/,/\{/{/\{/a \    include '"$NGINX_CONF_DIR"'/*.conf;
        }' "$NGINX_CONF"
    else
        echo "WARNING: Could not auto-add include. Please add manually:"
        echo "  include $NGINX_CONF_DIR/*.conf;"
    fi
fi

# ---- Test and reload ----
echo "Testing nginx config..."
if ! sudo nginx -t; then
    echo "ERROR: Nginx config test failed!"
    exit 1
fi

if sudo systemctl is-active nginx &>/dev/null; then
    sudo systemctl reload nginx
elif pgrep -x nginx &>/dev/null; then
    sudo nginx -s reload
else
    sudo nginx
    sudo systemctl enable nginx 2>/dev/null || true
fi
echo "Nginx configured for $DOMAIN → 127.0.0.1:$PORT"

# ---- Optional HTTPS with certbot ----
if command -v certbot &>/dev/null; then
    echo ""
    read -p "Setup HTTPS with Let's Encrypt? (y/N) " -r
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        sudo certbot --nginx -d "$DOMAIN"
    fi
else
    echo "TIP: Install certbot for HTTPS: sudo yum install -y certbot python3-certbot-nginx"
fi
