README

The camera module is based on the Motion V4L2 api. For more information visit https://motion-project.github.io/.

Installation:

Update your system:
sudo apt-get update
sudo apt-get upgrade
sudo apt-get install motion

Load drivers:
sudo modprobe gc2035
sudo modprobe vfe_v4l2

To automatically load the drivers on boot, add gc2035 and vfe_v4l2 to /etc/modules:

Modify configuration file:
sudo vim /etc/motion/motion.conf
daemon on
target_dir: change if you want but will be simpler to leave it the same
stream_localhost on
webcontrol_localhost on

cd /the/path/to/target_dir (default /var/lib/motion)
chmod 777 .

Change startup configuration
sudo vim /etc/default/motion
change “start_motion_demon” to  = “yes”

Useful controls for motion:
sudo service motion start
sudo service motion status
sudo service motion stop

Start motion using the appropriate command above.

If you did not change the default address of for file storage you can directly run the service
go run camera.go

If you have changed the default address, pass it to the service upon starting:
go run camera.go -dir /your/target/directory/here

You can now access the webcam remotely. Supported commands are start and last. More can be added as they become necessary, such as timelapse or livestream.

To take a picture from another computer on the network, you will need the ip address (enter below for xxx.xxx.x.xxx) of the Pi. The service listens on the port 8082.

To take a picture:
xxx.xxx.x.xxx:8082/take
To view the last picture
xxx.xxx.x.xxx:8082/last
