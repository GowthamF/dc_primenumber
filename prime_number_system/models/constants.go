package models

var PrimeNumberNode = "PRIMENUMBERNODE"
var MasterNode = "MASTERNODE"
var LearnerNode = "LEARNERNODE"
var ProposerNode = "PROPOSERNODE"
var AcceptorNode = "ACCEPTORNODE"

var UpStatus = "UP"
var DownStatus = "DOWN"

var SidecarPortNumber *string
var NodeId *string
var NumberOfProposers *int32
var HasLeaderElected bool = false
var MasterNodeId *string

var ElectionLock = "Election"
var MasterLock = "Master"
var ValidPrimeNumberMessage = "PRIMENUMBER"
var InvalidPrimeNumberMessage = "NOTAPRIMENUMBER"
