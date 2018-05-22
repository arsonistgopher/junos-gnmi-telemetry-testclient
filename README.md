## Basic Junos GNMI Telemetry Test Client

This is a Go script that demonstrates how to retrieve gNMI data from Junos via gRPC.

This test script takes the concept of [Nilesh Simaria's JTIMon](https://github.com/nileshsimaria/jtimon) and boils it down to the raw basics. 

Do not use this for anything other than curiosity!

## Usage

This package has been created with Godep support for dependencies.

```bash
go get github.com/arsonistgopher/junos-gnmi-telemetry-testclient.git
cd $GOHOME/src/github.com/arsonistgopher/junos-gnmi-telemetry-testclient
godep restore
go build
./junos-gnmi-telemetry-testclient
```

If you do not want to build, that's fine. I've got you covered. There are three binaries pre-compiled and ready to go.

```bash
junos-gnmi-telemetry-testclient-junos-32-0.1
junos-gnmi-telemetry-testclient-linux-64-0.1
junos-gnmi-telemetry-testclient-osx-0.1
```

To run the application, some command line inputs are required as below.

```bash
./junos-gnmi-telemetry-testclient-osx-0.1 -h
2018/05/22 18:10:19 -----------------------------------
2018/05/22 18:10:19 Junos gNMI Telemetry Test Tool
2018/05/22 18:10:19 -----------------------------------
2018/05/22 18:10:19 Run the app with -h for options

Usage of ./junos-gnmi-telemetry-testclient-osx-0.1:
  -certdir string
    	Directory with client.crt, client.key, CA.crt
  -cid string
    	Set to Client ID (default "1")
  -host string
    	Set host to IP address or FQDN DNS record (default "127.0.0.1")
  -loops int
    	Set number of times we should go through receive and print loop (default 2)
  -port string
    	Set to Server Port (default "32767")
  -subscription string
    	Set subscription to path (default "/interfaces/")
  -user string
    	Set to username (default "testuser")
```

Here is how to run it in case this still doesn't make sense.

```bash
./junos-gnmi-telemetry-testclient -certdir CLIENTCERT -cid 1 -host vmx -port 50051 -subscription /interfaces/ -user jet
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
./junos-gnmi-telemetry-testclient-osx-0.1 -loops 1 -host vmx -port 50051 -subscription /interfaces/ -user jet -certdir CLIENTCERT
2018/05/22 18:10:47 -----------------------------------
2018/05/22 18:10:47 Junos gNMI Telemetry Test Tool
2018/05/22 18:10:47 -----------------------------------
2018/05/22 18:10:47 Run the app with -h for options

Enter Password:
gnmiPath: ([]string) ["interfaces"]
pathToString: "interfaces"
gNMIPath: (*gnmi.Path) "element:\"interfaces\" elem:<name:\"interfaces\" > "


"subscribe:<prefix:<> subscription:<path:<element:\"interfaces\" elem:<name:\"interfaces\" > > mode:SAMPLE > encoding:PROTO > "

Path:  [name:"__juniper_telemetry_header__" ]
Value:  any_val:<type_url:"type.googleapis.com/GnmiJuniperTelemetryHeader" value:"\n\005vmx02\020\377\377\003\"/sensor_1000_4_1:/interfaces/:/interfaces/:mib2d(\200\200\200\001" >
Path:  [name:"__timestamp__" ]
Value:  uint_val:1527016251752
Path:  [name:"state"  name:"type" ]
Value:  string_val:"other"
Path:  [name:"state"  name:"mtu" ]
Value:  uint_val:65535
Path:  [name:"state"  name:"name" ]
Value:  string_val:"lsi"
Path:  [name:"state"  name:"description" ]
Value:  string_val:""
Path:  [name:"state"  name:"enabled" ]
Value:  bool_val:true
Path:  [name:"state"  name:"ifindex" ]
Value:  uint_val:4
Path:  [name:"state"  name:"admin-status" ]
Value:  string_val:"UP"
Path:  [name:"state"  name:"oper-status" ]
Value:  string_val:"UP"
Path:  [name:"state"  name:"last-change" ]
Value:  uint_val:203
Exit
```
