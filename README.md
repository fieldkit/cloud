# README for Cloud

## Overview

This repository contains:

1. All backend server-side code, written in Golang.
2. All public-facing custom websites, specifically the portal.
3. PostgreSQL schema, sample/seed data, and tools.

## Developers

### Dependencies

In order to build this project and be able to run it yourself, you'll need the following dependencies:

- [Docker (19.03.12)](https://docs.docker.com/get-docker/)
- [Golang (1.17)](https://golang.org/dl/)
- [nodejs (12.18.1)](https://nodejs.org/en/download/)
- [yarn (1.22.4)](https://classic.yarnpkg.com/en/docs/install#windows-stable)

**NOTE**: It should be possible to use `npm` if you'd like to avoid installing `yarn`.

### Overview

The server-side architecture is broken into two main parts:

- **Server**: A monolithic `golang` service that provides the backend REST API for FieldKit. It requires a PostgreSQL instance, which is covered later. It also reads objects from S3 for a more seamless developer experience. Any write that would go to S3 will go to the local file system when running in single-machine configuration.

- **Portal**: A Vuejs (JavaScript/TypeScript) single-page application that communicates with the backend server's API. This is the shiny and interesting stuff.

# Local Setup
## Prerequisites

Ensure your system is up to date:
```bash
sudo apt update
sudo apt upgrade
```

Clone the project:
```bash
git clone https://github.com/fieldkit/cloud
```

## Docker Setup

### Install Docker Desktop

1. Navigate to Downloads:
```bash
cd ~/Downloads
```

2. Install [Docker Desktop](https://www.docker.com/products/docker-desktop/):
```bash
sudo apt-get update
sudo apt-get install ./docker-desktop-<version>-<arch>.deb
```

### Install Docker CLI

Refer to the [StackOverflow](https://stackoverflow.com/questions/72299444/docker-desktop-doesnt-install-saying-docker-ce-cli-not-installable) thread for issues regarding Docker CLI installation.

To install, follow the steps below:

```bash
sudo apt install -y ca-certificates curl gnupg lsb-release
sudo mkdir -p /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt update -y
```

For additional information, check [Docker's official documentation](https://docs.docker.com/engine/install/ubuntu/#install-using-the-repository).

### Install Docker Compose

```bash
sudo apt install docker-compose
```

## NodeJS Setup

### Install NVM

```bash
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.4/install.sh | bash
```

### Install Node Versions

After installing NVM, restart the terminal and then install the desired versions of Node.js:

```bash
nvm install stable
nvm install lts/hydrogen
```

## Install Make

```bash
sudo apt install make
```

## Install Go

```bash
sudo snap install go --classic
```

## Install Yarn

```bash
npm install --global yarn
```

Navigate to the project's portal folder and run:

```bash
yarn install
```

## Additional Steps

1. Install tiptap:

```bash
npm install tiptap
```

## Running the Project

1. In one terminal, start Docker:

```bash
sudo docker-compose up
```

2. In a separate terminal, serve the project:

```bash
yarn serve
```

## Known Issues

- Issue with `@secrets` that's not getting resolved.
```

This should give you a nicely formatted README.md on GitHub with all the instructions and references intact. Adjustments might be needed depending on the specific layout you had in mind or the specifics of the project.

*EOF*
