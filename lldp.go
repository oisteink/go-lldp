package lldp

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

type LLDPS struct {
	LLDPS []LLDP `json:"lldp"`
}

type LLDP struct {
	LLDPInterfaces []LLDPInterface `json:"interface"`
}

type LLDPInterface struct {
	Name    string        `json:"name"`
	Via     string        `json:"via"`
	Rid     string        `json:"rid"`
	Age     string        `json:"age"`
	Chassis []LLDPChassis `json:"chassis"`
	Port    []LLDPPort    `json:"port"`
	Vlan    []LLDPVlan    `json:"vlan"`
}

type LLDPChassis struct {
	IDs           []LLDPChassisID          `json:"id"`
	Names         []LLDPChassisName        `json:"name"`
	Descriptions  []LLDPChassisDescription `json:"descr"`
	ManagementIPs []LLDPChassisManagemetIP `json:"mgmt-ip"`
	Capabilities  []LLDPChassisCapability  `json:"capability"`
}

type LLDPChassisID struct {
	ChassisType string `json:"type"`
	Name        string `json:"value"`
}

type LLDPChassisName struct {
	Name string `json:"value"`
}

type LLDPChassisDescription struct {
	Description string `json:"value"`
}

type LLDPChassisManagemetIP struct {
	IP string `json:"value"`
}

type LLDPChassisCapability struct {
	Type    string `json:"type"`
	Enabled bool   `json:"enabled"`
}

type LLDPPort struct {
	IDs          []LLDPPortID          `json:"id"`
	Descriptions []LLDPPortDescription `json:"descr"`
	TTLs         []LLDPPortTTL         `json:"ttl"`
	// driter i denne for nå. Sikker kjekt å vite om vi er på half/full duplex
	//AutoNegotiations []LLDPPortAutoNeg     `json:"auto-negotiation"`
}

type LLDPPortID struct {
	Type string `json:"type"`
	ID   string `json:"value"`
}

type LLDPPortDescription struct {
	Description string `json:"value"`
}

type LLDPPortTTL struct {
	TTL string `json:"value"`
}

type LLDPVlan struct {
	VlanID     string `json:"vlan-id"`
	IsPortVLAN bool   `json:"pvid"`
	Name       string `json:"value"`
}

func PrintLLDPS(value LLDPS) string {
	var result string
	for i := 0; i < len(value.LLDPS); i++ { //det er bare en, men vi looper uansett
		result += "Interfaces:\n"
		for a := 0; a < len(value.LLDPS[i].LLDPInterfaces); a++ {
			result += PrintInterface(value.LLDPS[i].LLDPInterfaces[a])
		}
	}
	return result
}

func PrintVLAN(vlan LLDPVlan) string {
	var result string

	result += fmt.Sprintln("\tVID:", vlan.VlanID)

	if vlan.Name != "" {
		result += fmt.Sprintln("(", vlan.Name, ")")
	}

	if vlan.IsPortVLAN {
		result += fmt.Sprintln("\tThis is Port VLAN id")
	}

	return result
}

func PrintChassis(chassis LLDPChassis) string {
	var result string
	//Name
	for i := 0; i < len(chassis.Names); i++ {
		result += fmt.Sprintln("\tName:", chassis.Names[i].Name)
	}

	//Type
	for i := 0; i < len(chassis.Capabilities); i++ {
		result += fmt.Sprintln("\tType:", chassis.Capabilities[i].Type, "Enabled:", chassis.Capabilities[i].Enabled)
		if chassis.Capabilities[i].Enabled {
			result += fmt.Sprintln("\tType:", chassis.Capabilities[i].Type)
		}
	}

	//ID
	for i := 0; i < len(chassis.IDs); i++ {
		result += fmt.Sprintln("\tID:", chassis.IDs[i].Name)
		result += fmt.Sprintln("\tID type", chassis.IDs[i].ChassisType)
	}

	//Managmenet IP
	for i := 0; i < len(chassis.ManagementIPs); i++ {
		result += fmt.Sprintln("\tManagement IP:", chassis.ManagementIPs[i].IP)
	}

	//Description
	for i := 0; i < len(chassis.Descriptions); i++ {
		//strip out any number of newlines
		result += fmt.Sprintln("\tDescription:", strings.Replace(chassis.Descriptions[i].Description, "\n", " ", -1))
	}
	return result
}

func PrintPort(port LLDPPort) string {
	var result string
	//ID
	for i := 0; i < len(port.IDs); i++ {
		result += fmt.Sprintln("\tID:", port.IDs[i].ID, "ID type:", port.IDs[i].Type)
	}

	//Description
	for i := 0; i < len(port.Descriptions); i++ {
		result += fmt.Sprintln("\tDescription:", port.Descriptions[i].Description)
	}

	//Time to live
	for i := 0; i < len(port.TTLs); i++ {
		result += fmt.Sprintln("\tTTL:", port.TTLs[i].TTL)
	}
	return result
}

func PrintInterface(iface LLDPInterface) string {
	var result string
	result += fmt.Sprintln("interface: ", iface.Name)
	result += fmt.Sprintln("Protocol:", iface.Via, "Age:", iface.Age)

	//chassis
	if len(iface.Chassis) > 0 {
		result += fmt.Sprintln("Chassis:")

		for i := 0; i < len(iface.Chassis); i++ {
			result += PrintChassis(iface.Chassis[i])
		}
	}

	//Port
	if len(iface.Port) > 0 {
		result += fmt.Sprintln("Ports:")
		for i := 0; i < len(iface.Port); i++ {
			result += PrintPort(iface.Port[i])
		}
	}

	//VlAN
	if len(iface.Vlan) > 0 {
		result += fmt.Sprintln("VLANs:")
		for i := 0; i < len(iface.Vlan); i++ {
			result += PrintVLAN(iface.Vlan[i])
		}
	}

	return result
}

//heller ha det som var? kanskje kjekt å ikke lage nye hele tiden, eller er go god på det?
func ReadLLDP() LLDPS {
	out, err := exec.Command("/usr/sbin/lldpcli", "show", "neigh", "det", "-f", "json0").Output()
	if err != nil {
		log.Fatal(err)
	}
	var result LLDPS
	if err = json.Unmarshal([]byte(out), &result); err != nil {
		log.Fatal(err)
	}
	return result
}
