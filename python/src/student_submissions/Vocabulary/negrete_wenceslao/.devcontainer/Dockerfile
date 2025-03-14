FROM python:3.10-slim

# Avoid warnings by switching to noninteractive
ENV DEBIAN_FRONTEND=noninteractive

# Configure apt and install packages
RUN apt-get update \
    && apt-get -y install --no-install-recommends apt-utils dialog 2>&1 \
    && apt-get -y install git procps lsb-release \
    && apt-get -y install gcc python3-dev \
    # Clean up
    && apt-get autoremove -y \
    && apt-get clean -y \
    && rm -rf /var/lib/apt/lists/*

# Switch back to dialog for any ad-hoc use of apt-get
ENV DEBIAN_FRONTEND=dialog

# Set up the Python development environment
RUN python -m pip install --upgrade pip \
    && pip install pylint black autopep8 flake8 mypy pycodestyle pydocstyle yapf \
    && mkdir -p /usr/local/py-utils/bin \
    && ln -s /usr/local/bin/pylint /usr/local/py-utils/bin/pylint \
    && ln -s /usr/local/bin/black /usr/local/py-utils/bin/black \
    && ln -s /usr/local/bin/autopep8 /usr/local/py-utils/bin/autopep8 \
    && ln -s /usr/local/bin/flake8 /usr/local/py-utils/bin/flake8 \
    && ln -s /usr/local/bin/mypy /usr/local/py-utils/bin/mypy \
    && ln -s /usr/local/bin/pycodestyle /usr/local/py-utils/bin/pycodestyle \
    && ln -s /usr/local/bin/pydocstyle /usr/local/py-utils/bin/pydocstyle \
    && ln -s /usr/local/bin/yapf /usr/local/py-utils/bin/yapf

# Create a non-root user to use if preferred
ARG USERNAME=vscode
ARG USER_UID=1000
ARG USER_GID=$USER_UID

# Create the user
RUN groupadd --gid $USER_GID $USERNAME \
    && useradd --uid $USER_UID --gid $USER_GID -m $USERNAME \
    && apt-get update \
    && apt-get install -y sudo \
    && echo $USERNAME ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/$USERNAME \
    && chmod 0440 /etc/sudoers.d/$USERNAME

# Set the default user
USER $USERNAME

# Add user to conda group for permissions
ENV PATH="/home/$USERNAME/.local/bin:${PATH}"

# Set the working directory
WORKDIR /workspace

# Warning: This is only for development, not for production use
ENV PYTHONPATH=${PYTHONPATH}:/workspace
