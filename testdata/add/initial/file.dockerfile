FROM ubuntu:22.04

# Metadata
LABEL maintainer="example@example.com"
LABEL version="1.0"

# Install requirements
RUN apt-get update && apt-get install -y \
    curl \
    && rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /hello

# Create hello world script
RUN echo '#!/bin/sh\necho "Hello, World!"' > /hello/hello.sh && \
    chmod +x /hello/hello.sh

# Set entry point
ENTRYPOINT ["/hello/hello.sh"]