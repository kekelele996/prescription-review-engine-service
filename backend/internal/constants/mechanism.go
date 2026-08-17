package constants

// 重复用药审核：命中以下作用机制时，重复开具判定为高风险（如两种 ACEI/ARB 降压药）。
var DuplicationHighRiskMechanisms = []string{
	"arb",           // ARB 类降压药
	"antiplatelet",  // 抗血小板药
	"vitamin_k_antagonist", // 维生素 K 拮抗剂
	"benzodiazepine", // 苯二氮䓬类
}

// IsDuplicationHighRiskMechanism 判断机制是否属于重复用药高风险机制。
func IsDuplicationHighRiskMechanism(m string) bool {
	for _, v := range DuplicationHighRiskMechanisms {
		if v == m {
			return true
		}
	}
	return false
}
