all: docker clean
	echo "All good!"

docker:
	./build-container.sh

clean:
	rm -rf generated-data

properclean:
	docker image rm hlefeuvre/linux-nginx

sanitycheck:
	docker run -it --privileged --security-opt seccomp:unconfined \
		hlefeuvre/linux-nginx /root/sanity_check.sh
