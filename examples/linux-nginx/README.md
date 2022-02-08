# Linux Wayfinder Nginx Experiment

```
The old pond;
the frog.
Plop!

 - Bashō
```

## Building the base container

It's easy:
```
$ make docker
```

Note here: we don't build everything in the dockerfile. Some part of it is
built by the bash script `build-container.sh` beforehand, and that's
because this part of the build requires docker and docker in docker is a
pain. To interested readers: feel free to change this if you're motivated.

### Changing the Nginx configuration (in the base filesystem)

It's located here, under `resources/nginx.conf`.

## Testing the container image

It's easy too:
```
$ make sanitycheck
```

If you can see wrk outputing a few 10s of thousands of requests per second,
you're good to go.

## Cleanup the environment

It's easy:
```
make clean
```

## IMPORTANT NOTES

If you change something in the filesystem image, such as the Nginx
configuration, but also the startup scripts, you HAVE to either run `make
clean` or delete `generated-data/nginx.ext2` otherwise it WILL NOT be
regenerated.
