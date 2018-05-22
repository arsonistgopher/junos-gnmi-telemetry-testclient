## Basic Junos OpenConfig telemetry tester

This is a Go script that demonstrates how to retrieve gNMI data from Junos via gRPC.

This test script takes the concept of [Nilesh Simaria's JTIMon](https://github.com/nileshsimaria/jtimon) and boils it down to the raw basics. 

Do not use this for anything other than curiosity!

## Usage

This package has been created with Godep support for dependencies.

```bash
go get github.com/arsonistgopher/junos-gnmi-telem-testclient.git
cd $GOHOME/src/github.com/arsonistgopher/junos-gnmi-telem-testclient
godep restore
go build
./junos-gnmi-telem-testclient
```

The script requires some command line inputs as below.

```bash
./junos-gnmi-telem-testclient -h
Usage of ./junos-gnmi-telem-testclient:
  -certdir string
    	Directory with clientCert.crt, clientKey.crt, CA.crt
  -cid string
    	Set to Client ID (default "1")
  -host string
    	Set host to IP address or FQDN DNS record (default "127.0.0.1")
  -loops int
    	Set loops to desired iterations (default 1)
  -port string
    	Set to Server Port (default "50051")
  -resource string
    	Set resource to resource path (default "/interfaces")
  -smpfreq int
    	Set to sample frequency in milliseconds (default 1000)
  -user string
    	Set to username (default "testuser")
```

Here is how to run it in case this still doesn't make sense.

```bash
./junos-gnmi-telem-testclient -cid 42 -host HOST -port 50051 -loops 1 -resource /interfaces -smpfreq 1000 -user jet
```
Replace `HOST` with the hostname or IP address of your code. Replace `50051` with the port your gRPC server on Junos is listening on. For the resource you want telemetry on, replace `/interfaces` with your chosen OpenConfig sensor.

For the readers amongst you, note that the password field is missing. This is requested from you and the output is masked to prevent shoulder surfer dangers!

Finally, here is a list of paths you can use for subscriptions to gNMI!

```bash
"/interfaces"
"/junos/system/linecard/packet/usage"
"/bgp"
"/components"
"/interfaces/interface/subinterfaces"
"/junos/npu-memory"
"/junos/system/linecard/npu/memory"
"/junos/task-memory-information"
"/junos/system/linecard/firewall/"
"/interfaces/interface[name='fxp0']/"
```