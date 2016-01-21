# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure(2) do |config|
  config.vm.box = "coreos-stable"
  config.vm.box_url = "http://stable.release.core-os.net/amd64-usr/current/coreos_production_vagrant.json"
  config.vm.network "forwarded_port", guest: 80, host: 3000
  config.vm.synced_folder ".", "/home/core/uuid-service"

  config.vm.provider "virtualbox" do |vb|
    vb.memory = "2048"
  end

  config.vm.provision "shell", inline: <<-SHELL
    docker run --name uuid-service-mongo \
      --restart=always \
      -d \
      mongo
    (cd /home/core/uuid-service && docker build -t uuid-service-golang .)
    docker run --name uuid-service-golang \
      -p 80:80 \
      --link uuid-service-mongo:mongo \
      --restart=always \
      -d \
      uuid-service-golang
  SHELL
end
