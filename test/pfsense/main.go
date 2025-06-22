package main

import (
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi-terraform-provider/sdks/go/libvirt/libvirt"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func stack(ctx *pulumi.Context) error {
	suffix, err := random.NewRandomString(ctx, "pfsense", &random.RandomStringArgs{
		Length: pulumi.Int(8),
	})
	if err != nil {
		return err
	}

	name := pulumi.Sprintf("pfsense-%s", suffix.Result)

	pool, err := libvirt.NewPool(ctx, "pfsense", &libvirt.PoolArgs{
		Name: name,
		Type: pulumi.String("dir"),
		Target: &libvirt.PoolTargetArgs{
			Path: pulumi.StringPtr("/tmp/pfsense"),
		},
	})
	if err != nil {
		return err
	}

	iso, err := libvirt.NewVolume(ctx, "pfsense-iso", &libvirt.VolumeArgs{
		Name:   pulumi.Sprintf("pfsense-iso-%s", suffix.Result),
		Pool:   pool.Name,
		Source: pulumi.StringPtr("netgate-installer-amd64.iso"),
	})
	if err != nil {
		return err
	}

	disk, err := libvirt.NewVolume(ctx, "pfsense", &libvirt.VolumeArgs{
		Name:   name,
		Pool:   pool.Name,
		Format: pulumi.String("qcow2"),
		Size:   pulumi.Float64Ptr(25_000),
	})
	if err != nil {
		return err
	}

	net, err := libvirt.NewNetwork(ctx, "pfsense", &libvirt.NetworkArgs{
		Name:      name,
		Mode:      pulumi.StringPtr("nat"),
		Addresses: pulumi.ToStringArray([]string{"192.168.100.0/24"}),
	})
	if err != nil {
		return err
	}

	_, err = libvirt.NewDomain(ctx, "pfsense", &libvirt.DomainArgs{
		Name:   name,
		Memory: pulumi.Float64Ptr(512),
		Vcpu:   pulumi.Float64Ptr(1),
		NetworkInterfaces: libvirt.DomainNetworkInterfaceArray{
			libvirt.DomainNetworkInterfaceArgs{
				NetworkId: net.NetworkId,
			},
		},
		Disks: libvirt.DomainDiskArray{
			libvirt.DomainDiskArgs{
				VolumeId: iso.VolumeId,
			},
			libvirt.DomainDiskArgs{
				VolumeId: disk.VolumeId,
			},
		},
		Consoles: libvirt.DomainConsoleArray{
			libvirt.DomainConsoleArgs{
				Type:       pulumi.String("pty"),
				TargetType: pulumi.StringPtr("virtio"),
				TargetPort: pulumi.String("1"),
			},
		},
		Graphics: libvirt.DomainGraphicsArgs{
			Type:       pulumi.StringPtr("spice"),
			ListenType: pulumi.StringPtr("address"),
			Autoport:   pulumi.BoolPtr(true),
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func main() {
	pulumi.Run(stack)
}
