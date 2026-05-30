package platform_test

import (
	"strings"
	"testing"

	"github.com/mitslab-de/omniinstall/internal/platform"
)

func TestDetectFromReader_Ubuntu(t *testing.T) {
	content := `NAME="Ubuntu"
ID=ubuntu
ID_LIKE=debian
PRETTY_NAME="Ubuntu 22.04.3 LTS"
VERSION_ID="22.04"
`
	family := platform.DetectFromReader(strings.NewReader(content))
	if family != platform.FamilyDebian {
		t.Errorf("expected FamilyDebian for Ubuntu, got %q", family)
	}
}

func TestDetectFromReader_Debian(t *testing.T) {
	content := `NAME="Debian GNU/Linux"
ID=debian
PRETTY_NAME="Debian GNU/Linux 12 (bookworm)"
`
	family := platform.DetectFromReader(strings.NewReader(content))
	if family != platform.FamilyDebian {
		t.Errorf("expected FamilyDebian for Debian, got %q", family)
	}
}

func TestDetectFromReader_Fedora(t *testing.T) {
	content := `NAME="Fedora Linux"
ID=fedora
PRETTY_NAME="Fedora Linux 38 (Workstation Edition)"
`
	family := platform.DetectFromReader(strings.NewReader(content))
	if family != platform.FamilyFedora {
		t.Errorf("expected FamilyFedora for Fedora, got %q", family)
	}
}

func TestDetectFromReader_RHEL(t *testing.T) {
	content := `NAME="Red Hat Enterprise Linux"
ID="rhel"
ID_LIKE="fedora"
`
	family := platform.DetectFromReader(strings.NewReader(content))
	if family != platform.FamilyFedora {
		t.Errorf("expected FamilyFedora for RHEL, got %q", family)
	}
}

func TestDetectFromReader_AlmaLinux(t *testing.T) {
	content := `NAME="AlmaLinux"
ID="almalinux"
ID_LIKE="rhel centos fedora"
`
	family := platform.DetectFromReader(strings.NewReader(content))
	if family != platform.FamilyFedora {
		t.Errorf("expected FamilyFedora for AlmaLinux, got %q", family)
	}
}

func TestDetectFromReader_Arch(t *testing.T) {
	content := `NAME="Arch Linux"
ID=arch
PRETTY_NAME="Arch Linux"
`
	family := platform.DetectFromReader(strings.NewReader(content))
	if family != platform.FamilyArch {
		t.Errorf("expected FamilyArch for Arch Linux, got %q", family)
	}
}

func TestDetectFromReader_Manjaro(t *testing.T) {
	content := `NAME="Manjaro Linux"
ID=manjaro
ID_LIKE=arch
`
	family := platform.DetectFromReader(strings.NewReader(content))
	if family != platform.FamilyArch {
		t.Errorf("expected FamilyArch for Manjaro (ID_LIKE=arch), got %q", family)
	}
}

func TestDetectFromReader_OpenSUSE(t *testing.T) {
	content := `NAME="openSUSE Leap"
ID="opensuse-leap"
ID_LIKE="suse opensuse"
`
	family := platform.DetectFromReader(strings.NewReader(content))
	if family != platform.FamilyOpenSUSE {
		t.Errorf("expected FamilyOpenSUSE for openSUSE, got %q", family)
	}
}

func TestDetectFromReader_Unknown(t *testing.T) {
	content := `NAME="SomeDistro"
ID=somedistro
`
	family := platform.DetectFromReader(strings.NewReader(content))
	if family != platform.FamilyUnknown {
		t.Errorf("expected FamilyUnknown for unknown distro, got %q", family)
	}
}

func TestDetectFromReader_Empty(t *testing.T) {
	family := platform.DetectFromReader(strings.NewReader(""))
	if family != platform.FamilyUnknown {
		t.Errorf("expected FamilyUnknown for empty input, got %q", family)
	}
}

func TestDetectFromReader_CommentLines(t *testing.T) {
	content := `# This is a comment
NAME="Ubuntu"
# Another comment
ID=ubuntu
ID_LIKE=debian
`
	family := platform.DetectFromReader(strings.NewReader(content))
	if family != platform.FamilyDebian {
		t.Errorf("expected FamilyDebian (comments should be ignored), got %q", family)
	}
}

func TestDetectFromReader_QuotedValues(t *testing.T) {
	content := `ID="ubuntu"
ID_LIKE="debian"
`
	family := platform.DetectFromReader(strings.NewReader(content))
	if family != platform.FamilyDebian {
		t.Errorf("expected FamilyDebian for quoted ID values, got %q", family)
	}
}

func TestDistroFamilyConstants(t *testing.T) {
	// Ensure constants have the expected string values (important for logging/output).
	if string(platform.FamilyDebian) != "debian" {
		t.Errorf("FamilyDebian should be 'debian', got %q", platform.FamilyDebian)
	}
	if string(platform.FamilyFedora) != "fedora" {
		t.Errorf("FamilyFedora should be 'fedora', got %q", platform.FamilyFedora)
	}
	if string(platform.FamilyArch) != "arch" {
		t.Errorf("FamilyArch should be 'arch', got %q", platform.FamilyArch)
	}
	if string(platform.FamilyOpenSUSE) != "opensuse" {
		t.Errorf("FamilyOpenSUSE should be 'opensuse', got %q", platform.FamilyOpenSUSE)
	}
	if string(platform.FamilyUnknown) != "unknown" {
		t.Errorf("FamilyUnknown should be 'unknown', got %q", platform.FamilyUnknown)
	}
}
