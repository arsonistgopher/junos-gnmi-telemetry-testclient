package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	auth_pb "github.com/arsonistgopher/gojtemtestgnmi/authentication"
	gnmipb "github.com/arsonistgopher/gojtemtestgnmi/proto/gnmi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// SUBSCRIPTIONPATH is for easy scripting
// "/interfaces"
// "/junos/system/linecard/packet/usage"
// "/bgp"
// "/components"
// "/interfaces/interface/subinterfaces"
// "/junos/npu-memory"
// "/junos/system/linecard/npu/memory"
// "/junos/task-memory-information"
// "/junos/system/linecard/firewall/"
// "/interfaces/interface[name='fxp0']/"
const SUBSCRIPTIONPATH = "/interfaces/"

func processGNMIResponse(resp *gnmipb.SubscribeResponse) {
	if notif := resp.GetUpdate(); notif != nil {
		// Have the full update
		// fmt.Printf("Update: %q\n", notif)

		// Print a section of each stuff
		// fmt.Println("Alias: ", notif.Alias)
		// fmt.Println("Delete: ", notif.Delete)
		// fmt.Println("Descriptor: ", notif.Descriptor)
		// fmt.Println("GetAlias: ", notif.GetAlias)
		// fmt.Println("Proto Message: ", notif.ProtoMessage)
		updates := notif.GetUpdate()

		for _, m := range updates {
			fmt.Println("Path: ", m.Path.GetElem())
			fmt.Println("Value: ", m.GetVal())

		}

		// fmt.Println("ProtoMessage :", notif.ProtoMessage)
	}
	if syncResp := resp.GetSyncResponse(); syncResp {
		fmt.Printf("Received sync-response\n")
		if false {
			os.Exit(0)
		}
	}
	if err := resp.GetError(); err != nil {
		fmt.Printf("Received error: %q\n", err)
	}
}

func subSendAndReceiveGNMI(conn *grpc.ClientConn, req *gnmipb.SubscribeRequest) {
	c := gnmipb.NewGNMIClient(conn)

	ctx := context.Background()

	client, err := c.Subscribe(ctx)
	if err != nil {
		log.Fatalf("Error invoking gnmi.subscribe(): %q", err)
	}
	if err := client.Send(req); err != nil {
		log.Fatalf("Error sending(): %q", err)
	}

	for {
		var resp *gnmipb.SubscribeResponse
		resp, err := client.Recv()
		if err != nil {
			log.Fatalf("Recv error: %s\n", err)
		}
		processGNMIResponse(resp)
	}
}

func xpathToGNMIpath(input string) ([]string, error) {
	path := strings.Trim(input, "/")
	var buf []rune
	inKey := false
	null := rune(0)
	for _, r := range path {
		switch r {
		case '[':
			if inKey {
				return nil, fmt.Errorf("malformed path, nested '[': %q ", path)
			}
			inKey = true
		case ']':
			if !inKey {
				return nil, fmt.Errorf("malformed path, unmatched ']': %q", path)
			}
			inKey = false
		case '/':
			if !inKey {
				buf = append(buf, null)
				continue
			}
		}
		buf = append(buf, r)
	}
	if inKey {
		return nil, fmt.Errorf("malformed path, missing trailing ']': %q", path)
	}
	return strings.Split(string(buf), string(null)), nil
}

func main() {

	// Set host
	hostandport := "vmx02.corepipe.co.uk:50051"

	// gRPC options
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInitialWindowSize(524288))

	creds, err := credentials.NewClientTLSFromFile("certs/vmx02.corepipe.co.uk.crt", "")
	if err != nil {
		logrus.Fatalf("Could not load certFile: %v", err)
	}
	opts = append(opts, grpc.WithTransportCredentials(creds))

	conn, err := grpc.Dial(hostandport, opts...)
	if err != nil {
		logrus.Fatalf("Error opening grpc.Dial(): %v", err)
	}
	// lazy close
	defer conn.Close()

	// Check for auth
	l := auth_pb.NewLoginClient(conn)
	dat, err := l.LoginCheck(context.Background(), &auth_pb.LoginRequest{UserName: "jet", Password: "Passw0rd", ClientId: "42"})

	if err != nil {
		logrus.Fatalf("Could not login: %v", err)
	}

	if dat.Result == false {
		logrus.Fatalf("LoginCheck failed\n")
	}

	if err != nil {
		logrus.Fatalf("%v", err)
	}

	s := &gnmipb.SubscribeRequest_Subscribe{
		Subscribe: &gnmipb.SubscriptionList{
			Mode:     gnmipb.SubscriptionList_STREAM,
			Prefix:   &gnmipb.Path{Target: ""},
			Encoding: gnmipb.Encoding_PROTO,
		},
	}

	gpath, err := xpathToGNMIpath(SUBSCRIPTIONPATH)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	fmt.Printf("gnmiPath: (%T) %q\n", gpath, gpath)
	fmt.Printf("pathToString: %q\n", pathToString(gpath))
	pp, err := StringToPath(pathToString(gpath), StructuredPath, StringSlicePath)
	if err != nil {
		log.Fatalf("Invalid path: %v", err)
	}
	s.Subscribe.Subscription = append(s.Subscribe.Subscription,
		&gnmipb.Subscription{
			Path:           pp,
			Mode:           getSMode("sample"),
			SampleInterval: uint64(0),
		})
	fmt.Printf("gNMIPath: (%T) %q\n", pp, pp)

	req := &gnmipb.SubscribeRequest{Request: s}
	fmt.Printf("\n\n%q\n\n", req)
	subSendAndReceiveGNMI(conn, req)

	if err != nil {
		logrus.Fatalf("%v", err.Error())
	}

	// /*
	fmt.Println("Exit")
}
