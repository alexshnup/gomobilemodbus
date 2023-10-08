
gomobile bind -o gomobilemodbus.aar -target=android .


Don't use "go mod vendor" because it will create a vendor folder with all the dependencies and the gomobile tool will not be able to find the dependencies.