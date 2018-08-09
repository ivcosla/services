#!bin/sh
# Since docker swarm pulls images using digest it won't tag the downloaded images
# hence, we need to tag the id with the image name plus latest tag, so autoupdater
# can perform docker inspect over it with name+tag

image_list=$(docker images skycoin/skywire -q)

# docker will exit with error code 0 even if the image doesn't exists, thus we must check it out
# and exit with error if this happens
if [ "${image_list}" == "" ]; then
    exit 1
fi

first_id=$(echo ${image_list} | awk '{print $1:}')

echo "found id ${first_id}"

docker tag ${first_id} skycoin/skywire:latest