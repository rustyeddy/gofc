package gofc

import (
	"fmt"
	"net"

	"github.com/rustyeddy/gofc/ofprotocol/ofp13"
	log "github.com/rustyeddy/logrus"
)

/**
 * basic controller
 */
type OFController struct {
	echoInterval int32 // echo interval
}

func NewOFController() *OFController {
	ofc := new(OFController)
	ofc.echoInterval = 60
	return ofc
}

func (c *OFController) HandleSwitchFeatures(msg *ofp13.OfpSwitchFeatures, dp *Datapath) {
	// create match
	fmt.Println("Handling switch stuff ! ")

	ethdst, _ := ofp13.NewOxmEthDst("00:00:00:00:00:00")
	if ethdst == nil {
		log.Info("  ~ ethdst prob ", ethdst)
		return
	}
	match := ofp13.NewOfpMatch()
	match.Append(ethdst)

	// create Instruction
	instruction := ofp13.NewOfpInstructionActions(ofp13.OFPIT_APPLY_ACTIONS)

	// create actions
	seteth, _ := ofp13.NewOxmEthDst("11:22:33:44:55:66")
	instruction.Append(ofp13.NewOfpActionSetField(seteth))

	// append Instruction
	instructions := make([]ofp13.OfpInstruction, 0)
	instructions = append(instructions, instruction)

	// create flow mod
	fm := ofp13.NewOfpFlowModModify(
		0, // cookie
		0, // cookie mask
		0, // tableid
		0, // priority
		ofp13.OFPFF_SEND_FLOW_REM,
		match,
		instructions,
	)

	fmt.Printf("  ~ send flow mod .. %#v ", fm)

	// send FlowMod
	dp.Send(fm)
	fmt.Println("  ~ create and send aggregate ... ")

	// Create and send AggregateStatsRequest
	mf := ofp13.NewOfpMatch()
	mf.Append(ethdst)
	mp := ofp13.NewOfpAggregateStatsRequest(0, 0, ofp13.OFPP_ANY, ofp13.OFPG_ANY, 0, 0, mf)
	dp.Send(mp)
}

func (c *OFController) HandleAggregateStatsReply(msg *ofp13.OfpMultipartReply, dp *Datapath) {
	log.Info("Handle AggregateStats")
	for _, mp := range msg.Body {
		if obj, ok := mp.(*ofp13.OfpAggregateStats); ok {
			fmt.Printf(" packets {%d} bytes {%d} flows {%d}\n",
				obj.PacketCount, obj.ByteCount, obj.FlowCount)
		}
	}
}

func (c *OFController) HandleEchoRequest(msg *ofp13.OfpHeader, dp *Datapath) {
	// send EchoReply
	echo := ofp13.NewOfpEchoReply()
	(*dp).Send(echo)
}

func (c *OFController) ConnectionUp() {
	// handle connection up
	log.Info("  TODO Connection Up ")
}

func (c *OFController) ConnectionDown() {
	// handle connection down
	log.Info("  TODO Connection Down ")
}

func (c *OFController) sendEchoLoop() {
	// send echo request forever
	log.Info("  TODO sendEchoLoop")
}

func ServerLoop(serverStr string) {
	if serverStr == "" {
		serverStr = ":6633"
	}
	tcpAddr, err := net.ResolveTCPAddr("tcp", serverStr)
	listener, err := net.ListenTCP("tcp", tcpAddr)

	ofc := NewOFController()
	GetAppManager().RegistApplication(ofc)

	if err != nil {
		return
	}

	// wait for connect from switch
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			return
		}
		log.Info("  connection to handle ")
		go handleConnection(conn)
	}
}

/**
 * hanleConnection hello style :)
 */
func handleConnection(conn *net.TCPConn) {

	// send hello
	hello := ofp13.NewOfpHello()
	_, err := conn.Write(hello.Serialize())
	if err != nil {
		fmt.Println("hello write ", err)
	}

	// create datapath
	dp := NewDatapath(conn)

	// launch goroutine
	go dp.recvLoop()
	go dp.sendLoop()
}
