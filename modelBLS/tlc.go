package modelBLS

import (
	"crypto/sha256"
	"fmt"
	"go.dedis.ch/kyber/v3/sign"
	"go.dedis.ch/kyber/v3/sign/bdn"
	"log"
	"strconv"
	"time"
)

const ChanLen = 500

var Logger1 *log.Logger

// Advance will change the step of the node to a new one and then broadcast a message to the network.
func (node *Node) Advance(step int) {
	node.TimeStep = step
	node.Acks = 0
	node.Wits = 0

	fmt.Printf("node %d , Broadcast in timeStep %d,%#v\n", node.Id, node.TimeStep, node.History)
	Logger1.SetPrefix(strconv.FormatInt(time.Now().Unix(), 10) + " ")
	Logger1.Printf("%d,%d\n", node.Id, node.TimeStep)

	msg := MessageWithSig{
		Source:  node.Id,
		MsgType: Raw,
		Step:    node.TimeStep,
		History: make([]MessageWithSig, 0),
	}
	node.CurrentMsg = msg
	msgBytes := node.ConvertMsg.MessageToBytes(msg)
	node.Comm.Broadcast(*msgBytes)
}

// waitForMsg waits for upcoming messages and then decides the next action with respect to msg's contents.
func (node *Node) WaitForMsg(stop int) {
	end := false
	msgChan := make(chan *[]byte, ChanLen)
	for node.TimeStep <= stop && !end {
		// For now we assume that the underlying receive function is blocking
		rcvdMsg := node.Comm.Receive()
		if rcvdMsg == nil {
			continue
		}
		msgChan <- rcvdMsg

		go func() {
			msgBytes := <-msgChan

			msg := node.ConvertMsg.BytesToModelMessage(*msgBytes)

			fmt.Printf("node %d in step %d ;Received MSG with step %d type %d source: %d\n", node.Id, node.TimeStep, msg.Step, msg.MsgType, msg.Source)

			// Used for stopping the execution after some timesteps
			if node.TimeStep == stop {
				fmt.Println("Break reached by node ", node.Id)
				end = true
				return
			}

			// If the received message is from a lower step, send history to the node to catch up
			if msg.Step < node.TimeStep {
				if msg.MsgType == Raw {
					msg.MsgType = Catchup
					msg.Step = node.TimeStep
					msg.History = node.History
					msgBytes := node.ConvertMsg.MessageToBytes(*msg)
					node.Comm.Broadcast(*msgBytes)
				}
				return
			}

			switch msg.MsgType {
			case Wit:

				if msg.Step > node.TimeStep+1 {
					return
				}

				// Verify that it's really witnessed by majority of nodes by checking the signature and number of them
				sig := msg.Signature
				mask := msg.Mask

				msg.Signature = nil
				msg.Mask = nil
				msg.MsgType = Raw

				h := sha256.New()
				h.Write(*node.ConvertMsg.MessageToBytes(*msg))
				msgHash := h.Sum(nil)

				keyMask, _ := sign.NewMask(node.Suite, node.PublicKeys, nil)
				err := keyMask.SetMask(mask)
				if err != nil {
					return
				}
				// Only new place!
				if keyMask.CountEnabled() < node.ThresholdAck {
					return
				}

				aggPubKey, err := bdn.AggregatePublicKeys(node.Suite, keyMask)
				if err != nil {
					panic(err)
				}

				// Verify message signature
				fmt.Println("RCVD AggSig: ", sig, "RCVD AggPub :", aggPubKey, "RCVD Hash :", msgHash)
				err = bdn.Verify(node.Suite, aggPubKey, msgHash, sig)
				if err != nil {
					panic(err)
					return
				}

				if msg.Step == node.TimeStep+1 { // Node needs to catch up with the message
					// Update nodes local history. Append history from message to local history
					node.History = append(node.History, *msg)

					// Advance
					node.Advance(msg.Step)
					node.Wits += 1

				} else if msg.Step == node.TimeStep {
					// Count message toward the threshold
					node.Wits += 1
					fmt.Printf("WITS: node %d , %d\n", node.Id, node.Wits)

					if node.Wits >= node.ThresholdWit {
						// Log the message in history
						node.History = append(node.History, *msg)
						// Advance to next time step
						node.Advance(node.TimeStep + 1)
					}
				}

			case Ack:
				// Checking that the ack is for message of this step
				if (msg.Source != node.CurrentMsg.Source) || (msg.Step != node.CurrentMsg.Step) || (node.Acks >= node.ThresholdAck) {
					return
				}
				fmt.Printf("received ACK. node %d %d\n", node.Id, msg.Source)

				// TODO First you have to verify signature! you have to change sig and type field for verification.
				sig := msg.Signature
				mask := msg.Mask

				msg.Signature = nil
				msg.Mask = nil
				msg.MsgType = Raw

				h := sha256.New()
				h.Write(*node.ConvertMsg.MessageToBytes(*msg))
				msgHash := h.Sum(nil)

				//fmt.Println("RCVD hash :",msgHash," MASK",mask)

				keyMask, _ := sign.NewMask(node.Suite, node.PublicKeys, nil)
				err := keyMask.SetMask(mask)
				if err != nil {
					panic(err)
					return
				}

				//fmt.Println(node.PublicKeys[keyMask.IndexOfNthEnabled(0)],"		",sig)
				//fmt.Println(node.PublicKeys)

				PubKey := node.PublicKeys[keyMask.IndexOfNthEnabled(0)]
				// Verify message signature
				err = bdn.Verify(node.Suite, PubKey, msgHash, sig)
				if err != nil {
					panic(err)
					return
				}

				// add message's mask to existing mask
				err = node.SigMask.Merge(mask)

				//

				// Count acks toward the threshold
				node.Acks += 1

				// Add signature to the list of signatures
				node.Signatures = append(node.Signatures, sig)

				// TODO a flaw here! Only send the TW message once not after every ack after reaching majority!
				if node.Acks >= node.ThresholdAck {
					// Send witnessed message if the acks are more than threshold
					msg.MsgType = Wit

					// Add aggregate signatures to message
					msg.Mask = node.SigMask.Mask()
					aggSignature, err := bdn.AggregateSignatures(node.Suite, node.Signatures, node.SigMask)
					if err != nil {
						panic(err)
					}
					msg.Signature, err = aggSignature.MarshalBinary()
					if err != nil {
						panic(err)
					}

					aggPubKey, err := bdn.AggregatePublicKeys(node.Suite, node.SigMask)

					// TODO signature is invalid here!
					fmt.Println("AggSig: ", msg.Signature, "AggPub :", aggPubKey, "Hash :", msgHash, "mask :", msg.Mask)
					err = bdn.Verify(node.Suite, aggPubKey, msgHash, msg.Signature)
					if err != nil {
						fmt.Println("PANIC AggSig: ", msg.Signature, "AggPub :", aggPubKey, "Hash :", msgHash, "mask :", msg.Mask)
						panic(err)
						return
					}
					fmt.Println("SIG OKAY")

					msgBytes := node.ConvertMsg.MessageToBytes(*msg)
					node.Comm.Broadcast(*msgBytes)
				}

			case Raw:
				if msg.Step > node.TimeStep+1 {
					return
				} else if msg.Step == node.TimeStep+1 { // Node needs to catch up with the message
					// Update nodes local history. Append history from message to local history
					node.History = append(node.History, *msg)

					// Advance
					node.Advance(msg.Step)
				}
				//fmt.Printf("ACKing by node %d, for msg %d\n", node.Id, msg.Source)

				// Node has to sign message hash
				h := sha256.New()
				h.Write(*msgBytes)
				msgHash := h.Sum(nil)
				//fmt.Println("SENT Hash :",msgHash)

				signature, err := bdn.Sign(node.Suite, node.PrivateKey, msgHash)
				if err != nil {
					panic(err)
				}

				// Adding signature and ack to message. These fields were empty when message got signed
				msg.Signature = signature

				// Add mask for the signature
				keyMask, _ := sign.NewMask(node.Suite, node.PublicKeys, nil)
				err = keyMask.SetBit(node.Id, true)
				if err != nil {
					panic(err)
				}
				msg.Mask = keyMask.Mask()

				// Send ack for the received message
				msg.MsgType = Ack
				msgBytes := node.ConvertMsg.MessageToBytes(*msg)
				node.Comm.Send(*msgBytes, msg.Source)

			case Catchup:
				if msg.Source == node.CurrentMsg.Source && msg.Step > node.TimeStep {
					fmt.Printf("Catchup: node (%d,step %d), msg(source %d ,step %d)\n", node.Id, node.TimeStep, msg.Source, msg.Step)
					node.History = append(node.History, msg.History[node.TimeStep:]...)
					node.Advance(msg.Step)
				}
			}
		}()

	}
}