package libvirtc

import (
	"github.com/danieldin95/lightstar/libstar"
	"github.com/libvirt/libvirt-go"
	"strings"
)

type HyperVisor struct {
	Name    string
	Schema  string
	Address string
	Path    string
	Conn    *libvirt.Connect
}

func parseQemuTCP(name string) (address, path string) {
	if strings.Contains(name, "://") {
		addrs := strings.SplitN(name, "://", 2)[1]
		address = strings.SplitN(addrs, "/", 2)[0]
		if strings.Contains(addrs, "/") {
			path = strings.SplitN(addrs, "/", 2)[1]
		}
	}
	return address, path
}

func parseQemuSSH(name string) (address, path string) {
	if strings.Contains(name, "://") {
		addrs := strings.SplitN(name, "://", 2)[1]
		address = strings.SplitN(addrs, "/", 2)[0]
		if strings.Contains(addrs, "/") {
			path = strings.SplitN(addrs, "/", 2)[1]
		}
		if strings.Contains(address, "@") {
			address = strings.SplitN(address, "@", 2)[1]
		}
	}
	return address, path
}

func (h *HyperVisor) Open() error {
	if hyper.Conn == nil {
		conn, err := libvirt.NewConnect(hyper.Name)
		if err != nil {
			return err
		}
		hyper.Conn = conn
	}
	if hyper.Conn == nil {
		return libstar.NewErr("Not connect.")
	}
	return nil
}
func (h *HyperVisor) Init() {
	if h.Name != "" {
		h.Schema = strings.SplitN(h.Name, ":", 2)[0]
		switch h.Schema {
		case "qemu+ssh":
			h.Address, h.Path = parseQemuSSH(h.Name)
		case "qemu+tcp", "qemu+tls":
			h.Address, h.Path = parseQemuTCP(h.Name)
		default:
			h.Address = "localhost"
			h.Path = "system"
		}
		if strings.Contains(h.Address, ":") {
			h.Address = strings.SplitN(h.Address, ":", 2)[0]
		}
	}
}

func (h *HyperVisor) GetCPU() (uint, string) {
	if err := h.Open(); err != nil {
		return 0, ""
	}
	if info, err := h.Conn.GetNodeInfo(); err == nil {
		return info.Cpus, info.Model
	}
	return 0, ""
}

func (h *HyperVisor) GetMem() (t uint64, f uint64, p float64) {
	if err := h.Open(); err != nil {
		return 0, 0, 0
	}
	if info, err := h.Conn.GetNodeInfo(); err == nil {
		t = info.Memory * 1024
	}
	if free, err := h.Conn.GetFreeMemory(); err == nil {
		f = free
	}
	return t, f, p
}

func (h *HyperVisor) GetRootfs() string {
	if err := h.Open(); err != nil {
		return ""
	}
	return ""
}

func (h *HyperVisor) ListAllDomains() ([]Domain, error) {
	if err := h.Open(); err != nil {
		return nil, err
	}

	domains, err := h.Conn.ListAllDomains(DOMAIN_ALL)
	if err != nil {
		return nil, err
	}
	newDomains := make([]Domain, 0, 32)
	for _, m := range domains {
		newDomains = append(newDomains, *NewDomainFromVir(&m))
	}
	return newDomains, nil
}

func (h *HyperVisor) LookupDomainByUUIDString(id string) (*Domain, error) {
	if err := h.Open(); err != nil {
		return nil, err
	}
	domain, err := hyper.Conn.LookupDomainByUUIDString(id)
	if err != nil {
		return nil, err
	}
	return NewDomainFromVir(domain), nil
}

func (h *HyperVisor) LookupDomainByUUIDName(id string) (*Domain, error) {
	if err := h.Open(); err != nil {
		return nil, err
	}
	domain, err := hyper.Conn.LookupDomainByUUIDString(id)
	if err != nil {
		domain, err := hyper.Conn.LookupDomainByName(id)
		if err != nil {
			return nil, err
		}
		return NewDomainFromVir(domain), nil
	}
	return NewDomainFromVir(domain), nil
}

func (h *HyperVisor) LookupDomainByName(id string) (*Domain, error) {
	if err := h.Open(); err != nil {
		return nil, err
	}
	domain, err := hyper.Conn.LookupDomainByName(id)
	if err != nil {
		return nil, err
	}
	return &Domain{*domain}, nil
}

func (h *HyperVisor) DomainDefineXML(xmlConfig string) (*Domain, error) {
	if err := h.Open(); err != nil {
		return nil, err
	}
	domain, err := h.Conn.DomainDefineXML(xmlConfig)
	if err != nil {
		return nil, err
	}
	return &Domain{*domain}, nil
}

var hyper = HyperVisor{
	Name:    "qemu:///system",
	Address: "localhost",
}

func GetHyper() (*HyperVisor, error) {
	if hyper.Conn == nil {
		conn, err := libvirt.NewConnect(hyper.Name)
		if err != nil {
			return &hyper, err
		}
		hyper.Conn = conn
	}
	return &hyper, nil
}

func SetHyper(name string) (*HyperVisor, error) {
	if name == hyper.Name {
		return &hyper, nil
	}
	hyper.Name = name
	hyper.Init()

	conn, err := libvirt.NewConnect(hyper.Name)
	if err != nil {
		return &hyper, err
	}
	hyper.Conn = conn
	return &hyper, nil
}

func CloseHyper(name string) {
	hyper.Conn.Close()
}

func init() {
	hyper.Init()
}

func LookupDomainByUUIDString(uuid string) (*Domain, error) {
	hyper, err := GetHyper()
	if err != nil {
		return nil, err
	}
	dom, err := hyper.LookupDomainByUUIDString(uuid)
	if err != nil {
		return nil, err
	}
	return dom, nil
}

func LookupDomainByUUIDName(uuid string) (*Domain, error) {
	hyper, err := GetHyper()
	if err != nil {
		return nil, err
	}
	dom, err := hyper.LookupDomainByUUIDName(uuid)
	if err != nil {
		return nil, err
	}
	return dom, nil
}