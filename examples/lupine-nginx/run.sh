docker build -f Dockerfile  -t wayfinder/lupine:latest . 
docker run  -v ~/wayfinder/examples/lupine-nginx/:/mnt/ --privileged -ti wayfinder/lupine sh /mnt/build.sh  
