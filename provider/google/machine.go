package google

type machineType struct {
	name string
	cpu  float32
	ram  float32
}

var (
	machineTypes = []machineType{
		{name: "f1-micro", cpu: 0.2, ram: 0.6},
		{name: "g1-small", cpu: 0.5, ram: 1.7},
		{name: "n1-standard-1", cpu: 1, ram: 3.75},
		{name: "n1-standard-2", cpu: 2, ram: 7.5},
		{name: "n1-standard-4", cpu: 4, ram: 15},
		{name: "n1-standard-8", cpu: 8, ram: 30},
		{name: "n1-standard-16", cpu: 16, ram: 60},
		{name: "n1-standard-32", cpu: 32, ram: 120},
		{name: "n1-standard-64", cpu: 64, ram: 240},
		{name: "n1-standard-96", cpu: 96, ram: 360},
	}
)

func getMachineType(cpu int, ram int) string {
	for _, machineType := range machineTypes {
		if float32(cpu) <= machineType.cpu && float32(ram) <= machineType.ram {
			return machineType.name
		}
	}

	return ""
}
