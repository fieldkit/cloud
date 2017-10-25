#+TITLE:	README for fk-cloud
#+EMAIL:	jacob@conservify.org

* fieldkit

** Developer Instructions

  sudo apt-get install docker.io docker-compose

  ./build.sh
 
 2) Add to /etc/hosts
127.0.0.1       fieldkit.org
127.0.0.1       api.fieldkit.org
127.0.0.1       www.fieldkit.org

** Rebuilding API
  # Rebuilding API
  # Ensure you have the latest go. >1.8

  go get -u github.com/goadesign/goa/..
  
  ge generate
  # On error about missing API definition, I had to rename my vendor directory.

