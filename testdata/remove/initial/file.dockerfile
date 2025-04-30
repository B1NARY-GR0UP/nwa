# Copyright 2025 BINARY Members
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

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