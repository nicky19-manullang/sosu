package service

type breakpoint struct {
	concLo, concHi float64
	ispuLo, ispuHi float64
}

var pm25Breakpoints = []breakpoint{
	{0, 15.5, 0, 50},
	{15.6, 55.4, 51, 100},
	{55.5, 150.4, 101, 200},
	{150.5, 250.4, 201, 300},
	{250.5, 500, 301, 500},
}

func PM25ToISPU(concentration float64) int {
	for _, bp := range pm25Breakpoints {
		if concentration >= bp.concLo && concentration <= bp.concHi {
			ispu := (bp.ispuHi-bp.ispuLo)/(bp.concHi-bp.concLo)*(concentration-bp.concLo) + bp.ispuLo
			return int(ispu + 0.5)
		}
	}
	if concentration > 500 {
		return 500
	}
	return 0
}

func ISPUCategory(ispu int) string {
	switch {
	case ispu <= 50:
		return "Baik"
	case ispu <= 100:
		return "Sedang"
	case ispu <= 200:
		return "Tidak Sehat"
	case ispu <= 300:
		return "Sangat Tidak Sehat"
	default:
		return "Berbahaya"
	}
}

func IsAQISafe(ispu int) bool {
	return ispu <= 100
}