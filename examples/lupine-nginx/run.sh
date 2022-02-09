docker run -d -p -v .:/mnt --name lupine-build lupine-nginx
docker exec -ti lupine-build sh /mnt/build.sh
docker rm lupine-build
docker run  -v ~/wayfinder/examples/lupine-nginx/:/mnt/ --privileged -ti wayfinder/lupine 
docker build -f Dockerfile  -t wayfinder/lupine:latest . 
