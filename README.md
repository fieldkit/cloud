# fieldkit


# Developer Instructions

  sudo apt-get install docker.io docker-compose

  ./build.sh

  # Rebuilding API
  # Ensure you have the latest go. >1.8

  go get -u github.com/goadesign/goa/..
  
  ge generate
  # On error about missing API definition, I had to rename my vendor directory.

