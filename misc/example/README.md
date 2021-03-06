# KVM

Linux KVM/Qemu

## CentOS7

The following that start a new CentOS7 instance.

### Initialize Local Datastore

    mkdir -p /lightstar/datastore && cd /lightstar/datastore
    mkdir -p 09b191af-b82a-4736-b492-d43224bb5379
    ln -s 09b191af-b82a-4736-b492-d43224bb5379 0 && cd 0
   
### Create One Image

    mkdir -p centos7 && cd centos7
    qemu-img create -f qcow2 disk0.qcow2 10G
    
### Add New Network

    yum install bridge-utils
    brctl addbr virbr0
    brctl addbr virbr1
    brctl addbr virbr2

### Start Instance 

    virsh create centos7.xml
    virsh start centos7
    
    virsh list --all
    virsh vncdisplay panabit-01

## Virtual Network

    virsh net-create virbr1.xml
    virsh attach-interface --domain centos7 --type bridge --source virbr1 --model virtio --config --persistent
    virsh domiflist centos7
    virsh detach-interface --domain centos7 --type bridge --mac 52:54:00:8e:20:b5 --config --persistent
    
## KVM UEFI

    virsh create kvm-uefi.xml

## KVM BIOS

    virsh create kvm-biso.xml
   
## Panabit
 
    qemu-img create -f qcow2 /var/lib/libvirt/images/panabit.disk-0.img 10G
    virsh create panabit.xml
    virsh list --all
    virsh vncdisplay panabit-01
    
    
## Sound 

    [root@249openstack ~]# cat /etc/libvirt/qemu.conf  | grep host_audio
    vnc_allow_host_audio = 0
    #nographics_allow_host_audio = 1
    [root@249openstack ~]# 
    
        

