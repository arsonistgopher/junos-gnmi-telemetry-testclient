## Basic Junos GNMI Telemetry Test Client

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
2018/05/22 17:32:45 -----------------------------------
2018/05/22 17:32:45 Junos gNMI Telemetry Test Tool
2018/05/22 17:32:45 -----------------------------------
2018/05/22 17:32:45 Run the app with -h for options

Usage of ./junos-gnmi-telem-testclient:
  -certdir string
    	Directory with clientCert.crt, clientKey.crt, CA.crt
  -cid string
    	Set to Client ID (default "1")
  -host string
    	Set host to IP address or FQDN DNS record (default "127.0.0.1")
  -port string
    	Set to Server Port (default "32767")
  -subscription string
    	Set subscription to path (default "/interfaces/")
  -user string
    	Set to username (default "testuser")
```

Here is how to run it in case this still doesn't make sense.

```bash
./junos-gnmi-telem-testclient -certdir CLIENTCERT -cid 1 -host vmx -port 50051 -subscription /interfaces/ -user jet
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

Finally finally, here is some test output! Use `Ctrl+c` to quit.

```bash
./junos-gnmi-telem-testclient -certdir CLIENTCERT -cid 1 -host vmx -port 50051 -subscription /interfaces/ -user jet

2018/05/22 17:40:02 -----------------------------------
2018/05/22 17:40:02 Junos gNMI Telemetry Test Tool
2018/05/22 17:40:02 -----------------------------------
2018/05/22 17:40:02 Run the app with -h for options

Enter Password:
gnmiPath: ([]string) ["interfaces"]
pathToString: "interfaces"
gNMIPath: (*gnmi.Path) "element:\"interfaces\" elem:<name:\"interfaces\" > "


"subscribe:<prefix:<> subscription:<path:<element:\"interfaces\" elem:<name:\"interfaces\" > > mode:SAMPLE > encoding:PROTO > "

Path:  [name:"__juniper_telemetry_header__" ]
Value:  any_val:<type_url:"type.googleapis.com/GnmiJuniperTelemetryHeader" value:"\n\005vmx02\020\377\377\003\"/sensor_1000_4_1:/interfaces/:/interfaces/:mib2d(\200\200\200\001" >
Path:  [name:"__timestamp__" ]
Value:  uint_val:1527014406155
Path:  [name:"state"  name:"type" ]
Value:  string_val:"other"
Path:  [name:"state"  name:"mtu" ]
Value:  uint_val:65535
Path:  [name:"state"  name:"name" ]
Value:  string_val:"lsi"
Path:  [name:"state"  name:"description" ]
<snip>
```
