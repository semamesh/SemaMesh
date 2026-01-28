// Simplified Spoon-feed version using cilium/ebpf library
func LoadBPF() {
    // 1. Load the compiled object file
    spec, err := ebpf.LoadCollectionSpec("bpf/sema_redirect.o")

    // 2. Load it into the Kernel
    coll, err := ebpf.NewCollection(spec)

    // 3. Attach it to the "Root Cgroup"
    // This ensures EVERY pod on the node is subject to this rule
    cgroupPath := "/sys/fs/cgroup" // Standard K8s cgroup root
    link, err := link.AttachCgroup(link.CgroupOptions{
        Path:    cgroupPath,
        Attach:  ebpf.AttachCGroupInet4Connect,
        Program: coll.Programs["sema_connect4"],
    })

    fmt.Println("SemaMesh Interceptor is ACTIVE.")
}