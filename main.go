package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"syscall"

	auth_pb "github.com/arsonistgopher/junos-gnmi-telem-testclient/authentication"
	gnmipb "github.com/arsonistgopher/junos-gnmi-telem-testclient/proto/gnmi"
	"golang.org/x/crypto/ssh/terminal"
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
// const SUBSCRIPTIONPATH = "/interfaces/"

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

func main() {

	log.Println("-----------------------------------")
	log.Println("Junos gNMI Telemetry Test Tool")
	log.Println("-----------------------------------")
	log.Print("Run the app with -h for options\n\n")

	// Parse flags
	var host = flag.String("host", "127.0.0.1", "Set host to IP address or FQDN DNS record")
	var subscription = flag.String("subscription", "/interfaces/", "Set subscription to path")
	var user = flag.String("user", "testuser", "Set to username")
	var port = flag.String("port", "50051", "Set to Server Port")
	var cid = flag.String("cid", "1", "Set to Client ID")
	var certDir = flag.String("certdir", "", "Directory with clientCert.crt, clientKey.crt, CA.crt")
	flag.Parse()

	// Set host
	hostandport := *host + ":" + *port

	// Grab password
	fmt.Print("Enter Password: \n")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Fatalf("Error reading password: %v", err)
	}
	password := string(bytePassword)

	// gRPC options
	var opts []grpc.DialOption

	// Are we going to run with TLS?
	runningWithTLS := false
	if *certDir != "" {
		runningWithTLS = true
	}

	opts = append(opts, grpc.WithInitialWindowSize(524288))

	// If we're running with TLS
	if runningWithTLS {

		// Grab x509 cert/key for client
		cert, err := tls.LoadX509KeyPair(fmt.Sprintf("%s/client.crt", *certDir), fmt.Sprintf("%s/client.key", *certDir))

		if err != nil {
			log.Fatalf("Could not load certFile: %v", err)
		}
		// Create certPool for CA
		certPool := x509.NewCertPool()

		// Get CA
		ca, err := ioutil.ReadFile(fmt.Sprintf("%s/CA.crt", *certDir))
		if err != nil {
			log.Fatalf("could not read ca certificate: %s", err)
		}

		// Append CA cert to pool
		if ok := certPool.AppendCertsFromPEM(ca); !ok {
			log.Fatal("Failed to append client certs")
		}

		// build creds
		creds := credentials.NewTLS(&tls.Config{
			RootCAs:      certPool,
			Certificates: []tls.Certificate{cert},
			ServerName:   *host,
		})

		if err != nil {
			log.Fatalf("Could not load clientCert: %v", err)
		}

		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else { // Else we're not running with TLS
		opts = append(opts, grpc.WithInsecure())
	}

	conn, err := grpc.Dial(hostandport, opts...)
	if err != nil {
		log.Fatalf("Error opening grpc.Dial(): %v", err)
	}
	// lazy close
	defer conn.Close()

	// Check for auth
	l := auth_pb.NewLoginClient(conn)
	dat, err := l.LoginCheck(context.Background(), &auth_pb.LoginRequest{UserName: *user, Password: password, ClientId: *cid})

	if err != nil {
		log.Fatalf("Could not login: %v", err)
	}

	if dat.Result == false {
		log.Fatalf("LoginCheck failed\n")
	}

	s := &gnmipb.SubscribeRequest_Subscribe{
		Subscribe: &gnmipb.SubscriptionList{
			Mode:     gnmipb.SubscriptionList_STREAM,
			Prefix:   &gnmipb.Path{Target: ""},
			Encoding: gnmipb.Encoding_PROTO,
		},
	}

	gpath, err := xpathToGNMIpath(*subscription)
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
		log.Fatalf("%v", err.Error())
	}

	fmt.Println("Exit")
}
