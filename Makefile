branch=master
build: 
	(cd $(HOME)/cco && git reset --hard HEAD && git fetch && git checkout ${branch} && git pull origin ${branch})
	(cd $(HOME)/cco/webapp/go && make build)
	sudo systemctl restart nginx
	sudo systemctl restart cco.golang.service 
analyze:
	sudo cp /dev/null /var/log/mysql/mysql-slow.log
	sudo cp /dev/null /var/log/nginx/access.log
	(cd $(HOME)/cco/bench && ./bench -remotes localhost:80)
	sudo alp --file=/var/log/nginx/access.log ltsv -r --sort sum | head -n 30
.PHONY: build
